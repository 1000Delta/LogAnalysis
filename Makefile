# 环境变量，目录不要省略末尾斜线
PIPELINELOADER_DIR=./elasticsearch/pipelineloader/

# 编译 pipelineloader 工具
pipelineloader:
	cd $(PIPELINELOADER_DIR) && go build

# 配置系统参数
.PHONY: sysconfigure
sysconfigure:
	sudo sysctl -w vm.max_map_count=262144

.PHONY: build
build: pipelineloader sysconfigure

# 启动容器
.PHONY: run
run: build
	docker-compose up -d --remove-orphans

# 编辑状态，修改部分文件权限
.PHONY: edit
edit:
	sudo chmod 775 -R filebeat/*

# 加载配置
.PHONY: configure
configure:
	$(PIPELINELOADER_DIR)pipelineloader

# 重新加载配置
.PHONY: reconfigure
reconfigure:
	sudo chmod 755 -R filebeat/*
	docker-compose restart filebeat

# 运行时操作
.PHONY: ps
ps:
	docker-compose ps

# 删除构建产物
.PHONY: rm
rm:
	rm $(PIPELINELOADER_DIR)pipelineloader