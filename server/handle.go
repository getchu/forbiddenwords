package server

import (
	"forbiddenwords/lib"
	"time"
)

type Handle struct {
	//词库
	words *Words
	//正则库
	regular *Regular
	//过滤库
	filter *Filter
	//配置
	config *lib.Config
	//更新
	updateChan chan int
}

func NewHandle(config *lib.Config) *Handle {
	h := &Handle{
		words:      NewWords(config.Engine),
		regular:    NewRegular(),
		filter:     NewFilter(),
		config:     config,
		updateChan: make(chan int),
	}
	h.update()
	//启动定时器，1小时更新一次词库
	go h.updateTimer()
	return h
}

func (h *Handle) updateTimer() {
	timer := time.NewTicker(time.Hour * 1)
	//timer := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-timer.C:
			h.config.Logger.Output(lib.LOG_INFO, "update by timer")
			h.update()
		case <-h.updateChan:
			h.config.Logger.Output(lib.LOG_INFO, "update by chan")
			h.update()
		}
	}
}

func (h *Handle) update() {
	//更新词库
	err := h.words.UpdateWords()
	if err != nil {
		h.config.LoggerError.Output(lib.LOG_ERROR, "update words error: "+err.Error())
	}
	//更新正则库
	err = h.regular.UpdateRegular()
	if err != nil {
		h.config.LoggerError.Output(lib.LOG_ERROR, "update regular error: "+err.Error())
	}
	//更新过滤库
	err = h.filter.UpdateFilter()
	if err != nil {
		h.config.LoggerError.Output(lib.LOG_ERROR, "update filter error: "+err.Error())
	}
}

const (
	//词库匹配
	MATCH_WAY_WORDS = 1 << iota
	//正则匹配
	MATCH_WAY_REGULAR
)
