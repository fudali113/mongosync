package main

import (
	"flag"
	"fmt"
	"github.com/fudali113/mongosync/sync"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

const version = 0.1

func main() {
	var src, dst, name, opStr, includeNS, excludeNS string
	var limit, updateTsLen, interval int
	flag.StringVar(&src, "src", "localhost:27017", "数据源数据库地址,支持任何 mongo 官方支持的连接字符串")
	flag.StringVar(&dst, "dst", "localhost:27017", "目标数据库地址,支持任何 mongo 官方支持的连接字符串")
	flag.StringVar(&name, "name", "", "转换上下文的名字, 推荐为每个转换设置一个特殊的名字; 默认值为 dst 参数")
	flag.IntVar(&limit, "limit", 1000, "每次从oplog.rs读取多少条数据进行转化")
	flag.StringVar(&opStr, "op-str", sync.DefaultOpStr, "加载哪些 op type 的数据进行转换， 默认以 `,` 分割")
	flag.IntVar(&updateTsLen, "update-ts-len", 10, "转换多少条数据同步一次 mongo.sync.log 里面的 ts 参数， 该 ts 参数用于下次获取数据的起点")
	flag.IntVar(&interval, "interval", 60, "同步间隔时间; unit: second")
	flag.StringVar(&includeNS, "includes", "", "只有在此集合中的 NS 才会被同步, 多个可以使用 `，` 分割; 可以使用 `dbName.*`（目前只支持这一种格式）只匹配某条数据库下面的 ns")
	flag.StringVar(&excludeNS, "excludes", "", "只有不在此集合中的 NS 才会被同步， 多个可以使用 `，` 分割； 可以使用 `dbName.*`（目前只支持这一种格式）只匹配某条数据库下面的 ns")

	help := flag.Bool("h", false, "帮助信息")
	showVersion := flag.Bool("v", false, "版本信息")
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}
	if *showVersion {
		fmt.Printf("mongosync version: %g \n", version)
		return
	}
	if name == "" {
		name = dst
	}
	ctx := sync.SyncCtx{
		Src:         src,
		Dst:         dst,
		Name:        name,
		Limit:       limit,
		OpStr:       opStr,
		UpdateTsLen: updateTsLen,
		Interval:    interval,
		IncludeNS:   splitString(includeNS, ","),
		ExcludeNS:   splitString(excludeNS, ","),
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
	log.Println("exiting, sleep 2 s, wating save log")
	time.Sleep(time.Second)
	os.Exit(0)
}

func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, sep)
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
		panic(err)
	}
}
