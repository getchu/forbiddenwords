package server

import (
	"fmt"
	"forbiddenwords/lib"
	"net/http"
	"os"
	"time"
)

type Server struct {
	Config *lib.Config
}

func NewServer(config *lib.Config) *Server {
	s := &Server{
		Config: config,
	}
	return s
}

func (s *Server) Start() {
	httpHandle := NewHTTPHandle(s.Config)
	http.HandleFunc("/", httpHandle.index)
	http.HandleFunc("/find", httpHandle.find)
	http.HandleFunc("/forbidLevel", httpHandle.forbidLevel)
	http.HandleFunc("/update", httpHandle.update)
	addr := s.Config.Ip + ":" + s.Config.Port
	fmt.Println("addr: " + addr + ". start...")
	go func() {
		time.Sleep(time.Second * 1)
		fmt.Println("ForbieednWord start \033[40G[\033[49;32;5mOK\033[0m]\n")
	}()
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("ForbieednWord start \033[40G[\033[49;32;5mFailed\033[0m]\n")
		fmt.Println(err)
		s.Config.LoggerError.Output(lib.LOG_ERROR, "start listen and serve error: "+err.Error())
		os.Exit(1)
	}
}
