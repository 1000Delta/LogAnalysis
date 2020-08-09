# LogAnalysis 日志分析系统

![Travis Status](https://travis-ci.org/1000Delta/LogAnalysis.svg?branch=master)

实现应用日志收集和分析功能

- [LogAnalysis 日志分析系统](#loganalysis-日志分析系统)
  - [技术架构](#技术架构)
    - [日志收集](#日志收集)
  - [Elastic 安装教程](#elastic-安装教程)
  - [依赖环境](#依赖环境)
  - [系统配置](#系统配置)
    - [禁用交换空间](#禁用交换空间)
    - [增加虚拟内存区域数量](#增加虚拟内存区域数量)
  - [部署配置](#部署配置)
    - [指定日志目录](#指定日志目录)
    - [X-PACK 配置](#x-pack-配置)
      - [ES 证书生成](#es-证书生成)
      - [创建/启动集群](#创建启动集群)
      - [配置 ES 集群内建用户密码](#配置-es-集群内建用户密码)
      - [配置 Kibana 连接 ES 密码](#配置-kibana-连接-es-密码)
      - [配置 filebeat 密码](#配置-filebeat-密码)
  - [开发记录](#开发记录)
    - [配置 ELK 集群](#配置-elk-集群)
    - [收集服务日志](#收集服务日志)
    - [filebeat + Logstash](#filebeat--logstash)
    - [filebeat + Elasticsearch Ingest](#filebeat--elasticsearch-ingest)
      - [方案 1](#方案-1)
      - [方案 2](#方案-2)
    - [ES pipeline 加载器](#es-pipeline-加载器)
    - [添加 pipeline processor 自动添加域名字段](#添加-pipeline-processor-自动添加域名字段)
      - [filebeat 自动加载自定义 pipeline](#filebeat-自动加载自定义-pipeline)
      - [pipeline 解析日志报错](#pipeline-解析日志报错)
    - [多行日志处理问题](#多行日志处理问题)
    - [部署问题](#部署问题)
    - [配置生成器](#配置生成器)
      - [模板空白符处理](#模板空白符处理)
      - [模板分文件处理](#模板分文件处理)
      - [docker-compose.yml 内容解读](#docker-composeyml-内容解读)
      - [配置生成存在的一些问题](#配置生成存在的一些问题)
  - [Kibana 配置](#kibana-配置)
    - [Elastic 安全问题](#elastic-安全问题)
      - [项目实际配置](#项目实际配置)
        - [kibana 安全性](#kibana-安全性)
        - [filebeat 安全性](#filebeat-安全性)
    - [Kibana 面板开发](#kibana-面板开发)

## 技术架构

通过 filebeat 收集日志，然后 Logstash 对日志数据进行处理，存储到 Elasticsearch，然后通过 Kibana 进行分析和展示

### 日志收集

filebeat 插件进行日志的收集，

## Elastic 安装教程

- [filebeat 安装配置](https://www.elastic.co/guide/en/beats/filebeat/7.8/running-on-docker.html)
- [Logstash 配置](https://www.elastic.co/guide/en/logstash/current/docker-config.html)
- [通过 Docker 安装 Elasticsearch](https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html)
- [通过 Docker 安装 Kibana](https://www.elastic.co/guide/en/kibana/current/docker.html)

## 依赖环境

go >= 1.11

Docker

docker-compose

## 系统配置

### 禁用交换空间

```shell
sudo swapoff -a
```

### 增加虚拟内存区域数量

`vm.max_map_count`: [elasticSearch Docker 启动报错：max virtual memory areas vm.max_map_count [65530] is too low, increase to at least [262144]](http://www.fecmall.com/topic/1172)

## 部署配置

### 指定日志目录

在 [`docker-compose.yml`](./docker-compose.yml) 中，对 filebeat 的 volume 配置项，把 `/1000delta/logs` 替换成部署机器的日志集合目录。二级目录将被识别为域名。

### X-PACK 配置

#### ES 证书生成

**部署集群之前**需要生成 ES 证书，执行：

```bash
make configure
```

将会创建临时 ES 容器生成证书，用于部署时映射到集群中。

#### 创建/启动集群

```bash
make run
```

启动 ELKB 集群。

#### 配置 ES 集群内建用户密码

只能执行一次，在**集群启动后**执行

```bash
make pass-es
```

#### 配置 Kibana 连接 ES 密码

输入配置 ES 时 `kibana_system` 的密码：

```bash
make pass-kibana
```

#### 配置 filebeat 密码

有两个密码，包括 ES 监控和推送日志到 ES。

ES 监控使用内建用户 `beats_system`，推送日志需要在 Kibana > Stack Management > 角色/用户 中配置角色和用户 `filebeat_publisher` 然后使用。

[具体操作](#filebeat-安全性)

记录密码命令：

```bash
make pass-filebeat
```

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

需要配置 Nginx 的 `log_format`，Nginx 官方文档：[Nginx log_format](http://nginx.org/en/docs/http/ngx_http_log_module.html#log_format)

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

20/07/25

完成 pipelineloader 的基本逻辑和测试用例，直接编译运行即可加载 pipeline 到 `127.0.0.1:9200` 上运行的 ES

20/07/26

> 题外话：阅读了一下 filebeat 的源码，其大量使用了工厂模式，之后可以去看一遍源码学习一下 filebeat 的架构

filebeat 对配置文件要求所有者为 filebeat 的用户并且仅用户可写，而 filebeat 是使用容器的 root 用户运行的，因此在外部修改需要使用 sudo 或者 root 用户修改，并且权限问题导致不能自动加载，需要手动重启非常不方便。

目前的解决方案是通过 Makefile 定制命令 `edit` 和 `reconfigure` 来实现修改文件权限和重新加载。 // TODO 优化 filebeat 配置加载

针对 pipelineloader 的编译和执行制定了 make 命令，之后可以尝试优化 // TODO 优化 make 命令

> Makefile 中，对于需要变更工作目录的命令，比如 `go build`，需要使用 `cd target && go build` 来执行命令，因为每条命令在 make 中都是开启一个 “sub shell”，上一条命令的 `cd` 对于下一条命令是无效的

20/07/27

目前从 filebeat 项目中下载了 Nginx module 的 pipeline 配置，然后需要根据生产环境的日志来做定制，添加 processor 来区分不同站点。

生产环境中，日志文件名相同而日志内容格式相同，没有区分不同的站点，对每个站点做配置又过于麻烦，通过 nginx module 获取到的日志中，我们可以看到一条字段：

![image-20200727164843328](readme.assets/image-20200727164843328.png)

`log.file.path` 这个字段包含了我们日志的目录规则：日志文件的上一级目录是站点域名，此时我们就可以添加一个 processor 来分割路径获取到站点域名，也就是我们 pipelineloader 发挥作用的时候了。

### 添加 pipeline processor 自动添加域名字段

Elastic search 官方文档关于 processor 编写

我们可以模仿 nginx module 中添加域名的字段 `destination.domain` 来添加字段。

通过使用 grok 和正则可以快速的提取出站点域名信息：

```yaml
- grok:
    field: log.file.path
    patterns:
      - "/(%{DATA}/)*%{DATA:destination.domain}/access.log"
    if: "ctx.destination.domain == null"
```

然后模拟 filebeat Nginx module 的标识，让数据条目在 Kibana 中可以看作 Nginx module 输出的:

```yaml
- set:
    field: event.module
    value: nginx
    if: ctx.event.module == null
```

运行 ES 集群之后，通过 pipelineloader 将 pipeline 加载到 ingest node 中，输出的数据已经符合了我们的要求。

上述 processor 是通用于 `access.log` 和 `error.log`，对于 `error.log`，Nginx module 中定义的 grok processor 只划分了基本的数据字段，对于错误的详细信息只存储在 `message` 中，对于错误定位和错误类型没有做划分，因此对于统计错误不太方便，我在 Logstash 定制 Grok filter 时，对于错误信息做了一个简单的分割器，可以对常见的错误定位和错误类型进行分割，便于统计错误类型。

在 Grok 文件中的格式：

```grok
NGINXERROR_DATE %{YEAR}/%{MONTHNUM}/%{MONTHDAY} %{TIME}
NGINXERROR_MESSAGE (?:%{GREEDYDATA:error.detail_before})?\(%{NUMBER:error.code}: %{GREEDYDATA:error_info}\)(?:%{GREEDYDATA:error.detail_end})?

# Error logs
NGINX_ERRORLOG %{NGINXERROR_DATE:timestamp} \[%{WORD:level}\] %{POSINT:pid}#%{NUMBER}: (?<error_message>%{NGINXERROR_MESSAGE}|%{GREEDYDATA})(?:, client: (?<remote_addr>%{IP}|%{HOSTNAME}))(?:, server: %{IPORHOST:server}?)(?:, request: %{QS:request})?(?:, upstream: (?<upstream>\"%{URI}\"|%{QS}))?(?:, host: %{QS:request_host})?(?:, referrer: \"%{URI:referrer}\")?
```

此处匹配逻辑借用了 Grok 默认模式中的 `httpd` 格式，对于错误信息和后续的 IP 和域名都做了划分，重点在于自定义的字段 `NGINXERROR_MESSAGE` 中，使用正则分别捕获错误定位和错误类型，Nginx 错误日志的错误类型有固定的格式为 `({code}: {info})`，而对于不符合内部错误格式的错误使用 `GREEDYDATA` 匹配即可。

> Elasticsearch 内置的 grok 模式可以参考 [`grok patterns`](https://github.com/elastic/elasticsearch/blob/5a5e11cf7d151636932a793ddbcc033675bd05ee/libs/grok/src/main/resources/patterns/grok-patterns)，对常用 patterns 做了定义。

修改为 pipeline 格式如下：

```yaml
  - grok:
      field: message
      patterns:
        - "%{NGINXERROR_DATE:timestamp} \[%{WORD:level}\] %{POSINT:pid}#%{NUMBER}: (?<error_message>%{NGINXERROR_MESSAGE}|%{GREEDYDATA})(?:, client: (?<remote_addr>%{IP}|%{HOSTNAME}))(?:, server: %{IPORHOST:server}?)(?:, request: %{QS:request})?(?:, upstream: (?<upstream>\"%{URI}\"|%{QS}))?(?:, host: %{QS:request_host})?(?:, referrer: \"%{URI:referrer}\")?"
      pattern_definitions:
        NGINXERROR_DATE: "%{YEAR}/%{MONTHNUM}/%{MONTHDAY} %{TIME}"
        NGINXERROR_MESSAGE: "(?:%{GREEDYDATA:error.detail_before})?\(%{NUMBER:error.code}: %{GREEDYDATA:error_info}\)(?:%{GREEDYDATA:error.detail_end})?"
```

由于在生产环境中，Nginx 可能对接到 `php-fpm` 等 FastCGI 协议的引擎，此时错误输出可能是引擎内部错误或者代码运行时错误，输出内容就不会遵循 Nginx 内部错误的格式，因此我们需要对可能的多行内容进行匹配，此时可以使用已经定义在配置中的 `GREEDYMULTILINE` 来指定，我们参考原本实现，最大限度复用原有配置，原 grok 如下：

```yaml
- grok:
    field: message
    patterns:
      - '%{DATA:nginx.error.time} \[%{DATA:log.level}\] %{NUMBER:process.pid:long}#%{NUMBER:process.thread.id:long}:
        (\*%{NUMBER:nginx.error.connection_id:long} )?%{GREEDYMULTILINE:message}'
    pattern_definitions:
      GREEDYMULTILINE: (.|\n|\t)*
    ignore_missing: true
```

最终修改配置如下：

```yaml
- grok:
    field: message
    patterns:
      - '%{DATA:nginx.error.time} \[%{DATA:log.level}\] %{NUMBER:process.pid:long}#%{NUMBER:process.thread.id:long}:(\*%{NUMBER:nginx.error.connection_id:long} )?(?<nginx.error.message>%{NGINXERROR_MESSAGE}|%{GREEDYDATA})(?:, client: (?<remote_addr>%{IP}|%{HOSTNAME}))(?:, server: %{IPORHOST:server}?)(?:, request: %{QS:request})?(?:, upstream: (?<upstream>\"%{URI}\"|%{QS}))?(?:, host: %{QS:request_host})?(?:, referrer: \"%{URI:referrer}\")?'
    pattern_definitions:
      NGINXERROR_DATE: "%{YEAR}/%{MONTHNUM}/%{MONTHDAY} %{TIME}"
      NGINXERROR_MESSAGE: '(?:%{GREEDYDATA})?\(%{NUMBER:nginx.error.code}: %{GREEDYDATA:nginx.error.info}\)(?:%{GREEDYDATA})?'
      GREEDYMULTILINE: |-
        (.|
        |	)*
    ignore_missing: true
```

实际上就是对最后一部分 `%{GREEDYMULTILINE:message}` 做了扩展，将报错信息进行了详细划分。

#### filebeat 自动加载自定义 pipeline

20/08/08

考虑到配置通过工具加载到 ES 中会导致先启动的 filebeat 传输数据无法被指定的 pipeline 解析，我查找了一下资料，发现都是通过手动加载，没有很好的自动加载的方式。

我通过 `docker cp` 导出 filebeat Docker 容器配置文件目录 `/usr/share/fielbeat` 后发现，其 `module/nginx/access/ingest` 和 `module/nginx/error/ingest` 分别存储了其加载到 ES 中的 pipeline 的配置，那么我们直接把这两个文件映射成我们需要的配置然后使用 `nginx` module 可不可以呢？

> 经过分析，在 `module/nginx` 文件夹中，配置对应的就是在 `filebeat.yml` 中写的配置，比如我自己对于 Nginx error 多行日志的处理，使用了 `multiline` 参数等，在 `module/nginx/error/config/nginx-error.yml` 中几乎一模一样，唯独对于多文件路径使用了 go template 进行遍历替换到配置文件中。

测试之后确认，可以通过替换 `pipeline.yml` 来使用自动加载，并且使用 nginx module 收集日志。

20/07/28

进行测试发现报错：`Provided Grok expressions do not match field value: [/logs/nginx/wsl.dev/error.log]`，说明编写的模式有问题

#### pipeline 解析日志报错

在检查日志解析的时候发现，对于 `access.log` 的解析经常出现报错信息：`[script] Too many dynamic script compilations within, max: [75/5m]; please use indexed, or scripts with parameters instead; this limit can be changed by the [script.max_compilations_rate] setting`，通过检查代码，发现是 processor 中 `if` 条件的表达式有问题：`"ctx?.destination?.domain? == null"`，修改成 `"ctx?.destination?.domain == null"` 即可，我对这里 `?` 的理解是先检查前置值是否存在，而最终字段不应该加上 `?`。

查询文档表明，ES 的 Painless 语法提供了 `?` 来保证 `null safe`（空值安全），否则会抛出 Java 经典的 NullPointerException 更多详细信息可以查阅文档：[Handling Nested Fields in Conditionals](https://www.elastic.co/guide/en/elasticsearch/reference/current/ingest-conditional-nullcheck.html#ingest-conditional-nullcheck)。

### 多行日志处理问题

多行日志，比如 `error.log` 中 php-fpm 输出的 PHP track，此时如果使用默认的 Log input，那么对于多行信息就会被当作多条信息，导致 pipeline 无法解析。通过设置 `multiline` 参数可以设置用于匹配多行数据起始的模式等参数，文档：[Manage multiline messages](https://www.elastic.co/guide/en/beats/filebeat/7.8/multiline-examples.html)

可以使用正则来编写模式，然后指定匹配与否和匹配位置，对于错误日志，起始为时间戳，格式为 `yyyy/MM/dd hh:mm:ss`，那么我们用以下设置：

```yaml
multiline.pattern: "^\d{4}\/\d{2}\/\d{2}"
multiline.negate: true
multiline.match: after
```

以上表示将“（以模式时间戳作为一行起始的数据）和（之后的行）”作为一条多行数据来处理。

20/07/29

### 部署问题

部署 ES 集群之后，不能直接启动 filebeat 收集日志，必须先将 pipeline 设置好，不然无法解析日志字段，那么我们需要在启动 filebeat 之前先运行 pipelineloader，要找到一种方法能够自动初始化。

### 配置生成器

对于日志在服务器的目录和 ES 节点数等配置，需要一种可以更快进行配置的方法，可以考虑使用 go template 对配置文件进行生成，通过交互式填写配置的方式填充模板内容，然后生成或覆盖配置文件。

go template 文档：[Package template](https://golang.org/pkg/text/template/)

通过对复合模板格式的文本创建模板对象然后导入模板中引用的结构体就可以导出解析后的文本，可以快速生成配置。

初步构想通过定义模型结构体和模板文件，最后编写交互逻辑，就可以方便地使用。

> go template 中有两个很相近的 Action，`with` 和 `if`，它们都是条件结构，区别在于，`with` 会将内部的 `.` 设置为 pipeline 的值，而 `if` 不会影响 `.` 的值

20/07/31

在模板中，占有一行的**无输出**插值表达式会导致出现空行，通过对分隔符加上 `-` 裁剪符可以除去空白符，比如 `{{-` 除去左侧空白符 `-}}` 除去右侧空白符，文档：[Text and spaces](https://golang.org/pkg/text/template/#hdr-Text_and_spaces).

#### 模板空白符处理

yaml 格式配置中，缩进是一个很重要的点，因此对于空白符的删除也要注意：

- 对于表达式和下一行文本在同一缩进的格式，使用右侧的裁剪符，相当于将表达式放在**下一个非空白符**之前。

- 对于表达式下一行文本不在同一缩进的格式，使用左侧裁剪符，相当于把表达式放在**上一个非空白符**末尾。

```template
services:
  {{ range .Nodes -}}
  {{ .Name }}:
    container_name: {{ .Name }}
    image: docker.elastic.co/elasticsearch/elasticsearch:7.8.0
      - LogAnalysis
  {{- end }}
networks:
```

此处 `{{ range .Nodes -}}` 使用右裁剪模式，相当于把 `{{ .Name }}` 提到了它这一行的开头，而因为 `range` 表达式缩进正确，因此输出也会正确。

对于 `{{- end }}` 表达式使用的是左裁剪，这样相当于它在上一行 `LogAnalysis` 的末尾，如果使用右裁剪的话，相当于把下一行的 `networks` 提到了和 `end` 一样的缩进，此时格式就错误了。

#### 模板分文件处理

对于不同部分的模板可以分开做处理，可以优化模板的格式，对不同内容的模板可以单独编辑，更加人性化。

使用嵌套模板：`{{ template "tmpl" }}` 表示使用名称为 tmpl 的模板替换当前内容，一般使用 `{{- template "tmpl" -}}` 的形式忽略语句带来的换行效果

嵌套模板的文件需要提前引入到模板对象中才能使用，否则会报错模板未定义

20/08/01

#### docker-compose.yml 内容解读

一开始 docker-compose.yml 中的内容是根据 Elastic 官方教程来的，虽然绝大多是配置项都能够理解作用是什么，但是还是有少数配置不理解，担心在部署后会因此导致生产环境的问题，因此需要查询文档将配置文件彻底弄清楚。

`ulimits`

对应的是 `docker run` 时的参数 `--ulimit`，用于设置 Linux ulimit.

Linux ulimit 对 shell 生成的进程的资源做限制，相关设置可以参照[这里](https://man.linuxde.net/ulimit)

官方文档对这个参数的描述只有寥寥几句：[Set ulimits in container (--ulimit)](https://docs.docker.com/engine/reference/commandline/run/#set-ulimits-in-container---ulimit)

对于配置中的 `type`，即配置项名称和含义没有在官方文档中列出来，参照 linux ulimit 的描述可以大概理解。

> 实际上 `--ulimit` 参数的 type 使用的是 Linux /etc/security/limits.conf 配置文件中对于系统参数的缩略词，可以参考这篇博客： [linux limits.conf 文件重要参数描述](https://blog.csdn.net/u012085379/article/details/53390203)
>
> 其正好也提到了 ES 相关的配置。

配置文件中 Elasticsearch 容器参数这几行：

```yaml
ulimits:
  memlock:
    soft: -1
    hard: -1
```

- `memlock` 配置项表示最大锁定内存地址空间(kb)
  - `soft` 表示 ulimit 的弹性限制，超过会发出警告
  - `hard` 表示 ulimit 的硬限制，超过会报错

#### 配置生成存在的一些问题

20/08/06

1. 配置生成对于配置容器的节点做了抽象，但是节点之间的依赖关系（`depends_on`），还有 Kibana 对 ES 的依赖 `ELASTICSEARCH_HOSTS` 这些应该如何定义？
2. 对于已经存在的配置，我如何对其中的部分属性做修改？比如我要调整 Kibana 的 basepath，但是我不希望修改其他的定制配置。

对于第二点可以考虑添加一个本地存储来实现，比如 一个 json/yaml/toml 等，每次启动 CLI 时加载，修改之后写入。 // TODO 定制配置缓存

## Kibana 配置

### Elastic 安全问题

Kibana 面板部署到服务器上对于公网是处于开放状态，如果限制 ip 访问控制效果太强，而且不灵活，优先考虑使用 Elastic 官方提供的安全措施。

配置安全性需要启用 xpack，然后要启用 TLS 认证，因此需要通过 es 生成证书，具体命令可以参考 [Encrypting communications in Elasticsearchedit](https://www.elastic.co/guide/en/elasticsearch/reference/current/configuring-tls.html)

使用证书启动 es 集群之后，需要配置内建用户的密码，用于管理集群和 ES 与 Kibana 之间的通信还有 Beats 的监控。

用户和角色管理：

1. 对于 filebeat 推送数据的角色，需要到 Kibana 的角色管理中设置，相关参数在 [Create a publishing user](https://www.elastic.co/guide/en/beats/filebeat/7.8/privileges-to-publish-events.html) 中。
2. 将角色的权限赋予 filebeat 使用的用户，你需要创建一个用户，然后赋予权限。

设置密码可以使用 keystore 来交互式设置，而不用硬编码到配置文件中。

- [Kibana](https://www.elastic.co/guide/en/kibana/current/secure-settings.html)
- [filebeat](https://www.elastic.co/guide/en/beats/filebeat/7.8/keystore.html)

#### 项目实际配置

##### kibana 安全性

使用内建用户 `kibana_system`，已经硬编码到配置中。

密码通过 make 命令 `pass-kibana` 来设置，调用 keystore 设置 `elasticsearch.password` 配置项。

##### filebeat 安全性

用户名称已经硬编码到配置文件中，为 `filebeat_publisher`。

需要在 Kibana 面板创建一个具有文档 [Create a publishing user](https://www.elastic.co/guide/en/beats/filebeat/7.8/privileges-to-publish-events.html) 中权限的角色，然后创建用户 `filebeat_publisher` 并赋予角色，密码通过 make 命令 `pass-filebeat` 使用 keystore 来添加。

### Kibana 面板开发
