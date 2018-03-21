package sync

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

// Connection 创建一个数据库先关连接
func Connection(url string) (conn *Conn, err error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Conn{Url: url, Session: session}, nil
}

// Oplogs query oplog.rs
// param limit query number
// param opStr filter op type ; split with `,` , such as `i,d,u`;
// param ts    query ts timestamp
func (conn *Conn) Oplogs(limit int, opStr string, ts ...bson.MongoTimestamp) (oplogs []Oplog) {
	criteria := bson.M{}
	if len(ts) > 0 {
		criteria["ts"] = bson.M{"$gt": ts[0]}
	}
	if opStr != "" && opStr != ALL_OPS {
		criteria["op"] = bson.M{"$in": strings.Split(opStr, ",")}
	}
	return conn.oplogs(limit, criteria)
}

func (conn *Conn) oplogs(limit int, criteria bson.M) (oplogs []Oplog) {
	oplogColl := conn.Session.DB("local").C("oplog.rs")
	query := oplogColl.Find(criteria).Sort("-ts")
	if limit != 0 {
		query.Limit(limit)
	}
	query.All(&oplogs)
	return
}

// MongoSyncLog get this context MongoSyncLog
func (conn *Conn) MongoSyncLog(syncName string) (log MongoSyncLog) {
	oplogColl := conn.Session.DB("local").C("mongo.sync.log")
	oplogColl.Find(bson.M{"syncName": syncName}).One(&log)
	return
}

// getNotDealOplogs
func (conn *Conn) GetNotDealOplogs() []Oplog {
	mongoSyncLog := conn.MongoSyncLog(conn.Ctx.Name)
	if mongoSyncLog.Id != "" {
		return conn.Oplogs(conn.Ctx.Limit, conn.Ctx.OpStr, mongoSyncLog.Ts)
	}
	return conn.Oplogs(conn.Ctx.Limit, conn.Ctx.OpStr)
}
