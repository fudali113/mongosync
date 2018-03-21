package sync

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
// param ts    query ts timestamp
func (conn *Conn) Oplogs(limit int, ts ...bson.MongoTimestamp) (oplogs []Oplog) {
	oplogColl := conn.Session.DB("local").C("oplog.rs")
	Criteria := bson.M{}
	if len(ts) > 0 {
		Criteria = bson.M{"ts": bson.M{"$gt": ts[0]}}
	}
	query := oplogColl.Find(Criteria).Sort("-ts")
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
		return conn.Oplogs(conn.Ctx.Limit, mongoSyncLog.Ts)
	}
	return conn.Oplogs(conn.Ctx.Limit)
}
