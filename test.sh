#!/bin/bash

root=$(pwd)

# pipelineloader
cd "./elasticsearch/pipelineloader" && go test
cd $root

# generator
cd "./generator" && go test 
cd $root
