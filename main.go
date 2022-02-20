/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"github.com/balchua/demo-spicedb/cmd"
	"go.uber.org/zap"
)

func init() {
	var err error
	config := zap.NewDevelopmentConfig()
	zapLog, err := config.Build(zap.AddCallerSkip(1))
	zap.ReplaceGlobals(zapLog)
	if err != nil {
		panic(err)
	}
}
func main() {
	cmd.Execute()
}
