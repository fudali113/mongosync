package sync

import "time"

// Sync 同步两个连接之间的数据
func Sync(src *Conn, dst *Conn) SyncResult {
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
		if num := i + 1; src.Ctx.UpdateTsLen < 1 || num%src.Ctx.UpdateTsLen == 0 || num == oplogsLen {
			src.saveMongoSyncLog(oplog)
		}
	}
	return SyncResult{
		Errs:  errs,
		Total: oplogsLen,
		Begin: begin,
		End:   time.Now(),
		OplogsResult: OplogsResult{
			BeginTs:  oplogsResult.BeginTs,
			Limit:    oplogsResult.Limit,
			ConnUrl:  oplogsResult.ConnUrl,
			Criteria: oplogsResult.Criteria,
		},
	}
}
