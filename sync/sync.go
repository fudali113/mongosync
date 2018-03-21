package sync

import (
	"context"
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
	oplogsResult := src.GetNotDealOplogs()
	oplogsLen := len(oplogsResult.Oplogs)
	errs := make([]SyncError, 0, 8)
	for i, oplog := range oplogsResult.Oplogs {
		err := dst.LoadOplog(oplog)
		if err != nil {
			errs = append(errs, SyncError{
				Err:   err,
				Index: i,
				Oplog: oplog,
			})
		}
		isSaveMongoSyncLog := false
		if num := i + 1; src.Ctx.UpdateTsLen < 1 || num%src.Ctx.UpdateTsLen == 0 || num == oplogsLen {
			src.saveMongoSyncLog(oplog)
			isSaveMongoSyncLog = true
		}
		if isDone(ctx) {
			if !isSaveMongoSyncLog {
				src.saveMongoSyncLog(oplog)
			}
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

func (conn *Conn) Sync(ctx context.Context, dst *Conn) {
	syncResult := Sync(ctx, conn, dst)
	conn.Session.DB("local").C("mongo.sync.result.log").Insert(syncResult)
}

func Run(sCtx SyncCtx) (cancelFunc context.CancelFunc, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	sCtxPtr := &sCtx
	src, err := Connection(sCtx.Src, sCtxPtr)
	if err != nil {
		return
	}
	dst, err := Connection(sCtx.Dst, sCtxPtr)
	if err != nil {
		return
	}
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
