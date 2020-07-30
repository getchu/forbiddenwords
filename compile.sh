#!/usr/bin/env bash

start(){
    echo $"Start..."
    nohup ./forbiddenwords &
}

run(){
    echo $"Run..."
    ./forbiddenwords
}

stop(){
    echo $"Stop..."
    ps aux | grep forbiddenwords | grep -v grep | awk '{print $2}' | xargs kill -9
}

build_linux(){
    echo $"Build for linux..."
    a="./forbiddenwords-"
    b=`date "+%Y%m%d%H%M%S"`
    c=$a$b
    GOOS=linux GOARCH=amd64 go build -o $c
}

gitpush(){
        echo $"Git add + git commit + git push..."
        git add .
        git commit -m "shell auto push"
        git pull
        git push
}

gitpull(){
        echo $"Git pull..."
        git pull
}

case "$1" in
   start)
        start
        exit 1
        ;;
   stop)
        stop
        exit 1
        ;;
   restart)
        echo $"Restart..."
        build_linux
        stop
        start
        exit 1
        ;;
   build-linux)
        build_linux
        exit 1
        ;;
   gitpush)
        gitpush
        exit 1
        ;;
   publish)
        build_linux
        gitpush
        exit 1
        ;;
   pull)
        gitpull
        run
        exit 1
        ;;
   *)
        echo $"Usage: $0 {start|stop|restart|build-linux|git push|publish|pull}"
        exit 1
        ;;
esac
