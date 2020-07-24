package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

const (
	pipelineDir = "./pipelines/"
)

// Pipeline 描述 pipeline 结构
type Pipeline struct {
	ID      string
	Content string
}

// GetPipelines 返回一组 pipeline 配置中的定义
func GetPipelines() ([]Pipeline, error) {
	var pipelines []Pipeline
	cfgDir, err := ioutil.ReadDir(pipelineDir)
	if err != nil {
		return nil, err
	}
	for _, cfg := range cfgDir {
		if cfg.IsDir() {
			continue
		}
		pipeline := Pipeline{}
		pipeline.ID = strings.Replace(cfg.Name(), filepath.Ext(cfg.Name()), "", -1) // 删除扩展名
		content, err := readConfig(pipelineDir + cfg.Name())
		if err != nil {
			return nil, err
		}
		contentBytes, err := json.Marshal(content)
		pipeline.Content = string(contentBytes)
		if err != nil {
			return nil, err
		}
		pipelines = append(pipelines, pipeline)
	}

	return pipelines, nil
}

// readConfig 返回解析 pipeline 配置文件的 map
func readConfig(p string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var content map[string]interface{}
	ext := filepath.Ext(p)
	switch ext {
	case ".json":
		err = json.Unmarshal(data, &content)
	case ".yml", ".yaml":
		err = yaml.Unmarshal(data, &content)
	default:
		return nil, fmt.Errorf("Unsupport format %s", ext)
	}
	if err != nil {
		return nil, err
	}

	return content, nil
}
