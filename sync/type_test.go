package sync

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func Test_Oplog_Struct(t *testing.T) {
	test_data := `{"ts" : Timestamp(1521526122, 1),"t" : NumberLong(1),"h" : NumberLong(-8496801878749629116),"v" : 2,"op" : "u","ns" : "mofangdb_report.trait"}`
	oplog := Oplog{}
	err := bson.UnmarshalJSON([]byte(test_data), &oplog)
	if err != nil {
		t.Error("反序列化出错", err)
	}
}
