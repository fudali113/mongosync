package main

import (
	"flag"
	"github/fudali113/mongosync/sync"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

func main() {
	var src, dst, name, opStr string
	var limit, updateTsLen, interval int
	flag.StringVar(&src, "src", "localhost:27017", "数据源数据库地址")
	flag.StringVar(&dst, "dst", "localhost:27018", "目标数据库地址")
	flag.StringVar(&name, "name", "", "转换上下文的名字")
	flag.IntVar(&limit, "limit", 1000, "每次从oplog.rs读取多少条数据进行转化")
	flag.StringVar(&opStr, "op-str", sync.DefaultOpStr, "加载哪些 op type 的数据进行转换， 默认以 `,` 分割")
	flag.IntVar(&updateTsLen, "update-ts-len", 10, "转换多少条数据同步一次 mongo.sync.log 里面的 ts 参数， 该 ts 参数用于下次获取数据的起点")
	flag.IntVar(&interval, "interval", 60, "同步间隔时间")
	ctx := sync.SyncCtx{
		Src:         src,
		Dst:         dst,
		Name:        name,
		Limit:       limit,
		OpStr:       opStr,
		UpdateTsLen: updateTsLen,
		Interval:    interval,
	}
	cancelFunc, err := sync.Run(ctx)
	checkErr(err)
	signals := make(chan os.Signal, 1)
	go func() {
		signal.Notify(signals, os.Interrupt, os.Kill)
	}()
	exitInfo := <-signals
	cancelFunc()
	log.Println("exit with: ", exitInfo.String())
	log.Println("exiting, sleep 1 s, wating save log")
	time.Sleep(time.Second)
	os.Exit(0)
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
		panic(err)
	}
}
