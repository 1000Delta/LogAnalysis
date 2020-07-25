package main

import (
	"context"
	"log"
	"strings"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func main() {
	esClient, err := es.NewDefaultClient()
	if err != nil {
		log.Fatalf("create Elasticsearch client error: %v", err.Error())
	}

	pipelines, err := GetPipelines()
	if err != nil {
		log.Fatalf("get pipelines error: %v", err.Error())
	}
	for _, pipeline := range pipelines {
		req := esapi.IngestPutPipelineRequest{
			PipelineID: pipeline.ID,
			Body:       strings.NewReader(pipeline.Content),
		}
		rsp, err := req.Do(context.Background(), esClient)
		if err != nil {
			log.Fatalf("request es client error: %v", err.Error())
		}
		defer rsp.Body.Close()
		if rsp.IsError() {
			log.Printf("request es pipeline error: %v", rsp.Body)
		}
	}
}

func init() {
	log.SetPrefix("[ES PipelineLoader]")
}
