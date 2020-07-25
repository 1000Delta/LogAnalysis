# pipelineloader

用于向 Elasticsearch 集群添加 ingest pipeline。

## Usage

```shell
go build
./pipelineloader
```

## Test

```shell
go generate # 创建es单节点容器，如果 127.0.0.1:9200 已经存在 es master node 请忽略此命令
go test
```
