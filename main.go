package main

import (
	"os"

	"github.com/op/go-logging"

	"github.com/ZanLabs/kai/db"
	"github.com/ZanLabs/kai/server"
	"github.com/ZanLabs/kai/types"
	"github.com/jinzhu/configor"
)

func initLogger(config types.ServerConfig) *logging.Logger {
	logger := logging.MustGetLogger(config.ServerName)
	logfile := config.System.Logfile
	f, err := os.Create(logfile)
	if err != nil {
		formatter := logging.MustStringFormatter(
			`%{color}%{time:15:04:05.000} package->%{shortpkg} func->%{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
		)
		backend := logging.NewLogBackend(os.Stderr, "", -1)
		backendFormatter := logging.NewBackendFormatter(backend, formatter)
		logging.SetBackend(backendFormatter)
	} else {
		formatter := logging.MustStringFormatter(
			`%{time:15:04:05.000} package->%{shortpkg} func->%{shortfunc} ▶ %{level:.4s} %{id:03x} %{message}`,
		)
		backend := logging.NewLogBackend(f, "", 0)
		backendFormatter := logging.NewBackendFormatter(backend, formatter)
		logging.SetBackend(backendFormatter)
	}

	return logger
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var config types.ServerConfig
	if err := configor.Load(&config, dir+"/config.yaml"); err != nil {
		panic(err)
	}

	logger := initLogger(config)

	db, err := db.GetDatabase(config.System)
	if err != nil {
		panic(err)
	}

	port := config.System.Port
	kaiServer := server.New(logger, config, "tcp", ":"+port, db)
	kaiServer.Start(true)
}
