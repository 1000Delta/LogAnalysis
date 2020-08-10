# 环境变量，目录不要省略末尾斜线
PIPELINELOADER_DIR=./elasticsearch/pipelineloader/

# 编译 pipelineloader 工具
elasticsearch/pipelineloader/pipelineloader: elasticsearch/pipelineloader/main.go elasticsearch/pipelineloader/pipeline.go
	cd $(PIPELINELOADER_DIR) && go build

# 配置系统参数
.PHONY: sysconfigure
sysconfigure:
	sudo sysctl -w vm.max_map_count=262144

.PHONY: build
build: elasticsearch/pipelineloader/pipelineloader sysconfigure

# 启动容器
.PHONY: run
run: build configure
	docker-compose up -d --remove-orphans

# 编辑状态，修改部分文件权限
.PHONY: edit
edit:
	sudo chmod 775 -R filebeat/*

# 配置集群参数
.PHONY: configure
configure: cert/elastic-certificates.p12
	
cert/elastic-certificates.p12:
	@echo 生成证书
	docker run -itd --rm --name la-es-ca docker.elastic.co/elasticsearch/elasticsearch:7.8.0
	docker exec -it la-es-ca bin/elasticsearch-certutil ca
	docker exec -it la-es-ca bin/elasticsearch-certutil cert --ca elastic-stack-ca.p12
	docker cp la-es-ca:/usr/share/elasticsearch/elastic-certificates.p12 ./cert/
	docker stop la-es-ca

# 配置集群密码
.PHONY: pass-es
pass-es:
	docker-compose exec es01 elasticsearch-setup-passwords interactive

.PHONY: pass-filebeat
pass-filebeat:
	docker-compose exec filebeat filebeat keystore create
	docker-compose exec filebeat filebeat keystore add SYSTEM_PWD
	docker-compose exec filebeat filebeat keystore add PUBLISHER_PWD
	
.PHONY: pass-kibana
pass-kibana:
	docker-compose exec kibana kibana-keystore create
	docker-compose exec kibana kibana-keystore add elasticsearch.password

# 重载配置相关操作
.PHONY: kibana
kibana: kibana/config/kibana.yml
	docker-compose restart kibana

.PHONY: filebeat
filebeat: filebeat/*
	sudo chmod 755 -R filebeat/*
	docker-compose restart filebeat 

# 重新加载配置
.PHONY: reconfigure
reconfigure: kibana filebeat

# 查看运行状态
.PHONY: ps
ps:
	docker-compose ps

# 停止容器运行
.PHONY: stop
stop:
	docker-compose stop

# 删除构建产物
.PHONY: rm
rm: stop
	rm $(PIPELINELOADER_DIR)pipelineloader