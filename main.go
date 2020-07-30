package main

import (
	"flag"
	"forbiddenwords/lib"
	"forbiddenwords/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

func main() {
	var ip = flag.String("ip", "0.0.0.0", "IP")
	var port = flag.String("port", "9121", "端口")
	var engine = flag.String("engine", "trie", "存储引擎：hash和trie")
	var logpath = flag.String("logpath", "/home/log/forbiddenwords/", "日志的存储路径")
	flag.Parse()
	config, err := lib.NewConfig(*ip, *port, *engine, *logpath)
	if config == nil || err != nil {
		flag.Usage()
		os.Exit(0)
	}
	ser := server.NewServer(config)
	go ser.Start()
	//监听信号
	listenSignal(ser)
}

// 监听信号
func listenSignal(ser *server.Server) {
	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTRAP)
	for {
		sig := <-sigChannel
		ser.Config.Logger.Output(lib.LOG_INFO, "main receive signal: %s", sig.String())
		if sig == syscall.SIGTRAP {
			continue
		}
		ser.Config.Logger.Output(lib.LOG_INFO, "main will be close master process。。.")
		break
	}
}

//------CPU监控开始----
//需要监控的时候，将下面2个方法放在需要监控的开始结束位置
var cpuF *os.File
var memF *os.File

func pprofStart() {
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()
	//cpu
	cpuF, err := os.OpenFile("./pprof/cpu.prof", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpuF)
	//内存
	memF, err = os.OpenFile("./pprof/mem.prof", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func pprofStop() {
	cpuF.Close()
	pprof.StopCPUProfile()

	pprof.WriteHeapProfile(memF)
	memF.Close()
}

//------CPU监控结束----
