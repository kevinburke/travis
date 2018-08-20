#!/bin/bash

pigeon ini.peg | goimports | gofmt -s > pigeon.go
