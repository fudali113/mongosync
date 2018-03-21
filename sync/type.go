package sync

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	Src   string
	Dst   string
	Name  string
	Limit int
	OpStr string
}

// Conn 封装一个数据库实体
type Conn struct {
	Url     string
	Ctx     SyncCtx
	Session *mgo.Session
}

// SyncResult 同步结果
type SyncResult struct {
	Errs  []error
	Total int64
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

// MongoSyncLog save mongo sync log
type MongoSyncLog struct {
	Id       bson.ObjectId       `bson:"_id"`
	Dst      string              `bson:"dst"`
	SyncName string              `bson:"syncName"`
	Ts       bson.MongoTimestamp `bson:"ts"`
	CreateAt time.Time           `bson:"createAt"`
}
