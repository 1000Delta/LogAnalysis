package main

import (
	"context"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func TestReadConfig(t *testing.T) {
	// test json
	if _, err := readConfig("./test/test.json"); err != nil {
		t.Errorf("读取 json 错误：%v", err.Error())
	}
	// test yaml
	if _, err := readConfig("./test/test.yml"); err != nil {
		t.Errorf("读取 yaml 错误：%v", err.Error())
	}
	// test other
	if _, err := readConfig("./test/test.other"); err == nil {
		t.Errorf("读取不支持类型未报错")
	}
}

func TestGetPipelines(t *testing.T) {
	if _, err := GetPipelines(); err != nil {
		t.Errorf("加载 pipeline 错误：%v", err.Error())
	}
}

// before run these test you must run a elasticsearch master node at 127.0.0.1:9200
//
//go:generate docker run -d -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.8.0

var esClient *es.Client

func init() {
	var err error
	esClient, err = es.NewDefaultClient()
	if err != nil {
		log.Fatalf("create Elasticsearch client error: %v", err.Error())
	}
}

func TestPutPipelines(t *testing.T) {
	pipelines, err := GetPipelines()
	if err != nil {
		t.Fatalf("get pipelines error: %v", err.Error())
	}
	for _, pipeline := range pipelines {
		req := esapi.IngestPutPipelineRequest{
			PipelineID: pipeline.ID,
			Body:       strings.NewReader(pipeline.Content),
		}
		rsp, err := req.Do(context.Background(), esClient)
		if err != nil {
			t.Fatalf("request es client error: %v", err.Error())
		}
		defer rsp.Body.Close()
		if rsp.IsError() {
			t.Logf("request es pipeline error: %v", rsp.Body)
		}
	}
}

func TestESGetPipelines(t *testing.T) {
	ppls, err := GetPipelines()
	if err != nil {
		t.Errorf("GetPipelines run error：%v", err.Error())
	}

	for _, ppl := range ppls {
		req := esapi.IngestGetPipelineRequest{
			PipelineID: ppl.ID,
		}
		rsp, err := req.Do(context.Background(), esClient)
		if err != nil {
			t.Errorf("request es client error: %v", err.Error())
		}
		defer rsp.Body.Close()
		msg, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			t.Errorf("read es response body error: %v", err.Error())
		}
		if rsp.IsError() {
			t.Errorf("request ingestgetpipeline error: %v", string(msg))
		}
		t.Log(string(msg))
	}
}
