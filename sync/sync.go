package sync

import (
	"context"
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

func isDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// Sync 同步两个连接之间的数据
func Sync(ctx context.Context, src *Conn, dst *Conn) SyncResult {
	begin := time.Now()
	log.Printf("开始一个同步周期, begin: %d", begin.Unix())
	oplogsResult := src.GetNotDealOplogs()
	oplogsLen := len(oplogsResult.Oplogs)
	defer func() {
		log.Printf("结束一个同步周期,len: %d, begin: %d, 耗时: %d s", oplogsLen, begin.Unix(), time.Now().Unix()-begin.Unix())
	}()
	errs := make([]SyncError, 0, 8)
	for i, oplog := range oplogsResult.Oplogs {
		if !needSync(src.Ctx, oplog) {
			log.Printf("不匹配相关信息， 不进行相关同步操作, oplog: %#v", oplog)
			continue
		}
		err := dst.LoadOplog(oplog)
		if err != nil {
			log.Printf("同步出错, err: %s", err.Error())
			errs = append(errs, SyncError{
				Err:   err,
				Index: i,
				Oplog: oplog,
			})
		}
		isSaveMongoSyncLog := false
		if num := i + 1; src.Ctx.UpdateTsLen < 2 || num%src.Ctx.UpdateTsLen == 0 || num == oplogsLen {
			_, err := src.saveMongoSyncLog(oplog)
			if err != nil {
				log.Printf("index: %d ; 更新 mongo.sync.log 中 ts 字段失败， 严重bug: %s", i, err.Error())
			} else {
				log.Printf("index: %d ; 同步了相关 oplog.ts 到 mongo.sync.log , ts: %+v", i, oplog.Ts)
			}
			isSaveMongoSyncLog = true
		}
		if isDone(ctx) {
			if !isSaveMongoSyncLog {
				src.saveMongoSyncLog(oplog)
			}
			log.Println("收到 ctx 结束信号， 退出循环并保存 操作时间日志")
			break
		}
	}
	return SyncResult{
		Errs:  errs,
		Total: oplogsLen,
		Begin: begin,
		End:   time.Now(),
		OplogsResult: OplogsResult{
			Name:     src.Ctx.Name,
			BeginTs:  oplogsResult.BeginTs,
			Limit:    oplogsResult.Limit,
			ConnUrl:  oplogsResult.ConnUrl,
			Criteria: oplogsResult.Criteria,
		},
	}
}

// Sync 执行同步，并将结束存入数据库
func (conn *Conn) Sync(ctx context.Context, dst *Conn) {
	syncResult := Sync(ctx, conn, dst)
	conn.Session.DB("local").C("mongo.sync.result.log").Insert(syncResult)
}

// Run 根据参数执行同步
func Run(sCtx SyncCtx) (cancelFunc context.CancelFunc, err error) {
	err = valid(sCtx)
	if err != nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	sCtxPtr := &sCtx
	src, err := Connection(sCtx.Src, sCtxPtr)
	if err != nil {
		return
	}
	src.Session.BuildInfo()
	localNames, _ := src.Session.DB("local").CollectionNames()
	hasOplogRs := false
	for _, name := range localNames {
		if name == "oplog.rs" {
			hasOplogRs = true
			break
		}
	}
	if !hasOplogRs {
		return nil, fmt.Errorf("请确保你的 src 数据库开启了 oplog 功能")
	}
	dst, err := Connection(sCtx.Dst, sCtxPtr)
	if err != nil {
		return
	}
	src.Sync(ctx, dst)
	go func() {
		for {
			select {
			case <-time.NewTicker(time.Duration(sCtx.Interval) * time.Second).C:
				src.Sync(ctx, dst)
			case <-ctx.Done():
				return
			}
		}
	}()
	return cancel, nil
}

// needSync 判断是否需要执行同步
func needSync(sc *SyncCtx, oplog Oplog) bool {
	result := true
	if len(sc.IncludeNS) > 0 {
		result = false
		for _, ns := range sc.IncludeNS {
			if matchNS(oplog.Ns, ns) {
				result = true
			}
		}
	}
	if len(sc.ExcludeNS) > 0 && result {
		for _, ns := range sc.ExcludeNS {
			if matchNS(oplog.Ns, ns) {
				result = false
			}
		}
	}
	return result
}

// matchNS 判断是否匹配
// 1. 相等
// 2. db 参数相等， 切 pattern 的集合字段为 *
func matchNS(src, pattern string) bool {
	srcInfo, _ := collInfo(src)
	patternInfo, _ := collInfo(pattern)
	if patternInfo.C == "*" && srcInfo.DB == patternInfo.DB {
		return true
	}
	return src == pattern
}

// valid 验证参数是否有错
func valid(ctx SyncCtx) error {
	if ctx.Limit < 1 {
		return fmt.Errorf("limit 参数不能够小于1， 您的limit参数是: %d", ctx.Limit)
	}
	if ctx.Interval < 1 {
		return fmt.Errorf("interval 参数不能够小于1， 您的 interval 参数是: %d", ctx.Interval)
	}

	// 检查是否 src 与 dst 中有相同的地址
	src, err := mgo.ParseURL(ctx.Src)
	if err != nil {
		return err
	}
	dst, err := mgo.ParseURL(ctx.Dst)
	if err != nil {
		return err
	}
	eqaulAddr := ""
addr:
	for _, addr := range src.Addrs {
		for _, addrByDst := range dst.Addrs {
			if addr == addrByDst {
				eqaulAddr = addr
				break addr
			}
		}
	}
	if eqaulAddr != "" {
		return fmt.Errorf("src 与 dst 包含了相同的服务地址: %s", eqaulAddr)
	}

	// 检查 includeNS 与 excludeNS 中是否有相同的元素
	equalNS := ""
ns:
	for _, ns := range ctx.IncludeNS {
		_, err := collInfo(ns)
		if err != nil {
			return fmt.Errorf("ns 字符串必须满足 *.* 的格式， 你的 ns 参数是： %s, 解析错误是： %s", ns, err.Error())
		}
		for _, exNS := range ctx.ExcludeNS {
			_, err := collInfo(exNS)
			if err != nil {
				return fmt.Errorf("ns 字符串必须满足 *.* 的格式， 你的 ns 参数是： %s, 解析错误是： %s", ns, err.Error())
			}
			if ns == exNS {
				equalNS = ns
				break ns
			}
		}
	}
	if equalNS != "" {
		return fmt.Errorf("includes 与 exludes 包含了相同的 NS : %s", equalNS)
	}
	return nil
}
