package sync

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
	"time"
)

// Connection 创建一个数据库先关连接
func Connection(url string, ctx *SyncCtx) (conn *Conn, err error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败，url: %s, error: %s", url, err.Error())
	}
	return &Conn{Url: url, Session: session, Ctx: ctx}, nil
}

// Oplogs query oplog.rs
// param limit query number
// param opStr filter op type ; split with `,` , such as `i,d,u`;
// param ts    query ts timestamp
func (conn *Conn) Oplogs(limit int, opStr string, ts ...bson.MongoTimestamp) (oplogs []Oplog) {
	return conn.oplogs(limit, oplogsCriteria(opStr, ts...))
}

func oplogsCriteria(opStr string, ts ...bson.MongoTimestamp) bson.M {
	criteria := bson.M{}
	if len(ts) > 0 {
		criteria["ts"] = bson.M{"$gt": ts[0]}
	}
	if opStr != "" && opStr != ALL_OPS {
		criteria["op"] = bson.M{"$in": strings.Split(opStr, ",")}
	}
	return criteria
}

func (conn *Conn) oplogs(limit int, criteria bson.M) (oplogs []Oplog) {
	oplogColl := conn.Session.DB("local").C("oplog.rs")
	query := oplogColl.Find(criteria).Sort("ts")
	if limit != 0 {
		query.Limit(limit)
	}
	query.All(&oplogs)
	return
}

// MongoSyncLog get this context MongoSyncLog
func (conn *Conn) MongoSyncLog(syncName string) (log MongoSyncLog) {
	syncLogColl := conn.Session.DB("local").C("mongo.sync.log")
	syncLogColl.Find(bson.M{"syncName": syncName}).One(&log)
	return
}

func (conn *Conn) saveMongoSyncLog(oplog Oplog) (info *mgo.ChangeInfo, err error) {
	syncLogColl := conn.Session.DB("local").C("mongo.sync.log")
	return syncLogColl.Upsert(bson.M{"syncName": conn.Ctx.Name}, bson.M{
		"ts":       oplog.Ts,
		"time":     time.Now(),
		"dst":      conn.Ctx.Dst,
		"syncName": conn.Ctx.Name,
	})
}

// getNotDealOplogs
func (conn *Conn) GetNotDealOplogs() OplogsResult {
	mongoSyncLog := conn.MongoSyncLog(conn.Ctx.Name)
	var oplogs []Oplog
	var criteria bson.M
	if mongoSyncLog.Id != "" {
		criteria = oplogsCriteria(conn.Ctx.OpStr, mongoSyncLog.Ts)

	} else {
		criteria = oplogsCriteria(conn.Ctx.OpStr)
	}
	oplogs = conn.oplogs(conn.Ctx.Limit, criteria)
	return OplogsResult{
		Oplogs:   oplogs,
		BeginTs:  mongoSyncLog.Ts,
		Limit:    conn.Ctx.Limit,
		ConnUrl:  conn.Url,
		Criteria: criteria,
	}
}

// LoadOplog load oplog to database
func (conn *Conn) LoadOplog(oplog Oplog) (err error) {
	collInfo, err := oplog.CollInfo()
	if err != nil {
		return
	}
	coll := conn.Session.DB(collInfo.DB).C(collInfo.C)
	switch oplog.Op {
	case INSERT:
		err = coll.Insert(oplog.O)
	case UPDATE:
		_, err = coll.Upsert(oplog.O2, oplog.O)
	case DELETE:
		err = coll.Remove(oplog.O)
	case NO_OP, DB, DB_CMD:
		log.Printf("op type '%s' 默认不进行任何操作", oplog.Op)
	default:
		return fmt.Errorf("不被支持的 op type: '%s' ", oplog.Op)
	}
	return
}
