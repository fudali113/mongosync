package sync

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

const (
	INSERT = "i"
	UPDATE = "u"
	DELETE = "d"
	// 声明当前数据库 (其中ns 被设置成为=>数据库名称+ '.')
	DB = "db"
	// no op,即空操作，其会定期执行以确保时效性 。
	NO_OP = "n"
	// db cmd
	DB_CMD  = "c"
	ALL_OPS = "*"
)

type SyncCtx struct {
	Src         string
	Dst         string
	Name        string
	Limit       int
	OpStr       string
	UpdateTsLen int
}

// Conn 封装一个数据库实体
type Conn struct {
	Url     string
	Ctx     *SyncCtx
	Session *mgo.Session
}

// SyncResult 同步结果
type SyncResult struct {
	Errs         []SyncError
	Total        int
	Begin        time.Time
	End          time.Time
	OplogsResult OplogsResult
}

type SyncError struct {
	Oplog Oplog
	Err   error
	Index int
}

// Oplogs mongo oplog.rs 实体数据结构
type Oplog struct {
	Ts bson.MongoTimestamp
	T  int32
	H  int32
	V  int8
	Op string
	Ns string
	O  bson.M
	O2 bson.M
}

type OplogsResult struct {
	Oplogs   []Oplog
	BeginTs  bson.MongoTimestamp
	Limit    int
	ConnUrl  string
	Criteria bson.M
}

type CollInfo struct {
	C  string
	DB string
}

// CollInfo get CollInfo by Oplog.Ns
func (oplog Oplog) CollInfo() (ci CollInfo, err error) {
	if oplog.Ns == "" {
		err = fmt.Errorf("Oplog.Ns is Empty")
		return
	}
	DBAndC := strings.Split(oplog.Ns, ".")
	if len(DBAndC) != 2 {
		err = fmt.Errorf("Oplog.Ns is %s , Invalid Format", oplog.Ns)
		return
	}
	ci.DB = DBAndC[0]
	ci.C = DBAndC[1]
	return
}

// MongoSyncLog save mongo sync log
type MongoSyncLog struct {
	Id       bson.ObjectId       `bson:"_id"`
	Dst      string              `bson:"dst"`
	SyncName string              `bson:"syncName"`
	Ts       bson.MongoTimestamp `bson:"ts"`
	Time     time.Time           `bson:"createAt"`
}
