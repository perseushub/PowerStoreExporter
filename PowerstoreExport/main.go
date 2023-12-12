package main

import (
	"github.com/go-kit/log"
	"powerstore/route"
	"powerstore/utils"
)

var (
	loggers log.Logger
	config  *utils.Config
)

func init() {
	config = utils.GetConfig()
	loggers = utils.GetLogger(config.Log.Level, config.Log.Path, config.Log.Type)
}

func main() {
	route.Run(config, loggers)
}
