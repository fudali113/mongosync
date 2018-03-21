package sync

import "testing"

func TestConn_Oplogs(t *testing.T) {
	conn := getConn()
	oplogs := conn.Oplogs(1)
	if len(oplogs) != 1 {
		t.Error("Conn.Oplogs has bug")
	}
}

func TestConn_MongoSyncLog(t *testing.T) {
	conn := getConn()
	log := conn.MongoSyncLog(conn.Ctx.Name)
	if log.Dst != "localhost:27017" {
		t.Error("Conn.MongoSyncLog has bug")
	}
}

func TestConn_GetNotDealOplogs(t *testing.T) {
	conn := getConn()
	oplogs := conn.GetNotDealOplogs()
	if len(oplogs) < 2 {
		t.Error("Conn.GetNotDealOplogs has bug")
	}
}

func getConn() *Conn {
	conn, err := Connection("192.168.30.249:28717")
	if err != nil {
		panic(err)
	}
	conn.Ctx = SyncCtx{Name: "test"}
	return conn
}
