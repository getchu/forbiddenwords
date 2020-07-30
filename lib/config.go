package lib

import (
	"fmt"
	"strconv"
	"time"
)

type Config struct {
	Ip          string
	Port        string
	Engine      string
	Logger      *Logger
	LoggerError *Logger
}

func NewConfig(ip, port, engine, logpath string) (*Config, error) {
	if !checkConfig(ip, port, engine) {
		return nil, fmt.Errorf("config invalid")
	}
	logger := NewLogger(logpath+"access.log", "[INFO]", "info")
	loggerError := NewLogger(logpath+"error.log", "[ERROR]", "warning")

	c := &Config{
		Ip:          ip,
		Port:        port,
		Engine:      engine,
		Logger:      logger,
		LoggerError: loggerError,
	}
	go c.UpdateConfig(logpath)
	return c, nil
}

func (c *Config)UpdateConfig(logpath string){
	timer := time.NewTicker(time.Hour * 1)
	for {
		select {
		case <-timer.C:
			c.Logger = NewLogger(logpath+"access.log", "[INFO]", "info")
			c.LoggerError = NewLogger(logpath+"error.log", "[ERROR]", "warning")
		}
	}
}

func checkConfig(ip, port, engineType string) bool {
	ipInt, err := strconv.Atoi(port)
	if err != nil || len(ip) < 7 || ipInt < 3000 || ipInt > 65535 || (engineType != "hash" && engineType != "trie") {
		return false
	}
	return true
}
