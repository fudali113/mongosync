package sync

import (
	"testing"
)

func TestConn_Oplogs(t *testing.T) {
	conn := getConn()
	oplogs := conn.Oplogs(1, ALL_OPS)
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
	oplogsResult := conn.GetNotDealOplogs()
	oplogs := oplogsResult.Oplogs
	if len(oplogs) < 2 {
		t.Error("Conn.GetNotDealOplogs has bug")
	}
	if oplogs[0].Ts > oplogs[1].Ts {
		t.Error("oplogs 排序不正确")
	}
}

func getConn() *Conn {
	conn, err := Connection("192.168.30.249:28717", &SyncCtx{Name: "test", OpStr: DefaultOpStr})
	if err != nil {
		panic(err)
	}
	return conn
}
