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
    - [filebeat + Elasticsearch Ingest](#filebeat--elasticsearch-ingest)
      - [方案 1](#方案-1)
      - [方案 2](#方案-2)
    - [ES pipeline 加载器](#es-pipeline-加载器)

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

查询字段信息可以发现，包含日志文件的路径信息，那么可以采用另一种方式，即通过 Logstash 对 filebeat 日志作二次分析，从路径中解析出站点名称等。

实操之后发现问题在于，使用 filebeat 作为 output 接收 module 处理的日志时，日志信息没有被解析！还是在 `message` 字段中.

经过查询资料了解到，filebeat module 处理 Nginx 日志是通过 elasticsearch 的 ingest node 进行预处理的，可以支持比如 Logstash 中的 Grok 等处理方式，那我还要个 p 的 Logstash，直接上 filebeat 处理就够了。

关于 ingest 预处理的介绍可以看这篇博客：[filebeat 解析日志 并发送到 Elasticsearch](https://www.hezehua.net/csdn_article-104184444)

关于 filebeat 在 ingest 中配置的 pipeline 详细内容可以看 github 上的配置文件：

- `access.log`: [filebeat/module/nginx/access/ingest/pipeline.yml](https://github.com/elastic/beats/blob/48160563134a8014119ceb96d3d404ac98cd9947/filebeat/module/nginx/access/ingest/pipeline.yml)
- `error.log`: [filebeat/module/nginx/error/ingest/pipeline.yml](https://github.com/elastic/beats/blob/bca0adcd4353d2a547e73cb9523a456971d9dc27/filebeat/module/nginx/error/ingest/pipeline.yml)

### filebeat + Elasticsearch Ingest

有两种方案：

1. 通过 filebeat log input 采集数据，并且通过 `add_fields` 添加标识字段，再由 Nginx module 内置的 pipeline 对数据进行分析。
2. 查看 filebeat nginx module 的 pipeline 模式可以发现在 access 标准日志之前有一个用于匹配的字段 `(%{NGINX_HOST} )?` 可以在这里增加站点名称用于分析。

#### 方案 1

经过测试可行，但是在配置文件中需要添加大量区分文件的字段。

#### 方案 2

需要配置 Nginx 的 `log_format`，官方文档：[Nginx log_format](http://nginx.org/en/docs/http/ngx_http_log_module.html#log_format)

存在的问题就是无法对旧日志进行分析，因为旧格式不包含站点字段。

还有一个问题在于，错误日志是没有站点信息的，而且 pipeline 的格式中也没有，因此只能自己做解析。

20/07/22

我们可以通过基于 fb nginx module 添加字段的方式来解析文件目录中站点名称来区分日志。

因此我们需要自行导入 pipeline，可以使用 HTTP 或者 Client API 来添加 pipeline。

### ES pipeline 加载器

我们可以参考 filebeat 的 pipeloader 进行设计，相关代码：[`fileset.LoadPipelines`](https://github.com/elastic/beats/blob/bca0adcd4353d2a547e73cb9523a456971d9dc27/filebeat/fileset/pipelines.go#L61-L113)

1. 这个函数的大致逻辑就是遍历已经注册的 module，对于每个 module 遍历其定义的 fileset（文件类型集合）
2. 调用 `fileset.GetRequiredProcessors` 检查 ES 是否启用依赖的 Processer
3. 调用 fileset 对象的方法 [`fileset.GetPipelines`](https://github.com/elastic/beats/blob/bca0adcd4353d2a547e73cb9523a456971d9dc27/filebeat/fileset/fileset.go#L420-L472) 获取写在配置文件（pipeline.yml）中的 ingest pipeline
4. 通过 `esClient.GetVersion` 获取 ES 版本号，检查版本号和 pipelines 的数量，只有版本号 >= 6.5.0 的 Elasticsearch 才支持 multi pipeline
5. 通过 `loadPipeline` 函数检查 pipeline 是否已经加载到 es 中
6. 如果加载报错，则回滚 pipeline，通过调用 `deletePipeline` 删除

20/07/23

加载 pipeline 的逻辑就在 [`loadPipeline`](https://github.com/elastic/beats/blob/bca0adcd4353d2a547e73cb9523a456971d9dc27/filebeat/fileset/pipelines.go#L115-L136) 中，函数首先调用 `makeIngestPipelinePath` 获取到 pipelineID 对应的路径，然后判断如果覆盖参数 `overwrite == true` 则检查 ID 是否存在，存在则直接返回；检查逻辑之后则是调用 `setECSProcessors` 添加 pipeline 然后进行错误检查。

我们需要检查 [`setECSProcessors`](https://github.com/elastic/beats/blob/bca0adcd4353d2a547e73cb9523a456971d9dc27/filebeat/fileset/pipelines.go#L140-L170) 中的插入逻辑。

20/07/24

上一部分看错了， `setECSProcessors` 是用于更改 Processor 属性为 ECS 的，加载 pipeline 的逻辑是之后调用的 `esClient.LoadJson`.

ECS 的介绍在 [这里](https://www.elastic.co/guide/en/ecs/current/ecs-reference.html)

filebeat 中，加载 pipeline 的方法是使用了其内置的 `esClient.Request` 来向节点发送请求，而我用的 [`go-elasticsearch`](https://github.com/elastic/go-elasticsearch) 并非如此使用，而是先通过 `esapi.IngestPutPipelineRequest` 创建请求，然后通过 `es.Client` 进行发送，处理逻辑上有一定差异。

首先需要实现一个配置文件解析器，参考官方的解析方式，我们不需要读取其字段，只需要判断格式有效即可，实际内容交给 es 去读取。

编写了 pipeline 解析函数和测试用例。
