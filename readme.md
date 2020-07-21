# LogAnalysis 日志分析系统

实现应用日志收集和分析功能

- [LogAnalysis 日志分析系统](#loganalysis-日志分析系统)
  - [技术架构](#技术架构)
    - [日志收集](#日志收集)
  - [Elastic 安装教程](#elastic-安装教程)
  - [系统配置](#系统配置)
    - [禁用交换空间](#禁用交换空间)
    - [增加虚拟内存区域数量](#增加虚拟内存区域数量)
  - [开发记录](#开发记录)
    - [配置 ELK 集群](#配置-elk-集群)
    - [收集服务日志](#收集服务日志)
    - [filebeat + Logstash](#filebeat--logstash)

## 技术架构

通过 filebeat 收集日志，然后 Logstash 对日志数据进行处理，存储到 Elasticsearch，然后通过 Kibana 进行分析和展示

### 日志收集

filebeat 插件进行日志的收集，

## Elastic 安装教程

- [filebeat 安装配置](https://www.elastic.co/guide/en/beats/filebeat/7.8/running-on-docker.html)
- [Logstash 配置](https://www.elastic.co/guide/en/logstash/current/docker-config.html)
- [通过 Docker 安装 Elasticsearch](https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html)
- [通过 Docker 安装 Kibana](https://www.elastic.co/guide/en/kibana/current/docker.html)

## 系统配置

### 禁用交换空间

```shell
sudo swapoff -a
```

### 增加虚拟内存区域数量

`vm.max_map_count`: [elasticSearch Docker 启动报错：max virtual memory areas vm.max_map_count [65530] is too low, increase to at least [262144]](http://www.fecmall.com/topic/1172)

## 开发记录

### 配置 ELK 集群

按照官方教程通过 Docker 安装集群，集群配置写在 `docker-compose.yml` 中。

对于 Elasticsearch 相关的

### 收集服务日志

通过官方教程将 ELK 集群搭建完成后，需要配置日志收集的渠道。

日志收集通过 filebeat 进行，通过 Docker 部署方式中，已经自动配置了 Docker 服务发现，会自动收集 Docker 服务的日志，我们对容器 Label 进行设置之后，filebeat 会自动调用对应的 module 对日志进行处理。

我们还需要通过文件来收集其他日志。

首先查找 filebeat 读取和跟踪文件日志相关的文档：[Log input](https://www.elastic.co/guide/en/beats/filebeat/7.8/filebeat-input-log.html#filebeat-input-log) 用于收集单行的 `.log` 文件，多行日志可以使用 [Configure inputs](https://www.elastic.co/guide/en/beats/filebeat/7.8/configuration-filebeat-options.html) 中的 multiline messages 输入。

最主要需要收集的是 Nginx 访问日志和错误日志

如何配置 filebeat 区分不同站点的 nginx 日志？

Nginx module 中没有提到相关内容，只能配置使用 nginx 模块的目录和文件

想法如下：

- 使用多个 filebeat 实例（占用资源多）
- 通过配置 log input 添加 field 的方式（实际效果未知，没有找到相关文档）

网络上大多是通过 json 收集 nginx 日志，使用 nginx 模块的也没有找到多站点的记录，因此先尝试第二种，第二种不行再使用第一种。

文件日志，映射到 filebeat 容器的 `logs` 目录，然后二级目录区分不同应用，通过配置 Input 的方式来读取不同日志文件。

### filebeat + Logstash

20/07/21 filebeat module 和 input 是独立的，在 module 中配置的路径会被自动读取而无需设置 input，因此无法添加字段。

通过可以采用另一种方式，即通过 Logstash 对 filebeat 日志作二次分析，
