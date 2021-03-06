package sync

import (
	"context"
	"testing"
)

func Test_isDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := isDone(ctx)
	if done {
		t.Error("done has bug  1  ")
	}
	cancel()
	done = isDone(ctx)
	if !done {
		t.Error("done has bug  2  ")
	}
}

func Test_matchNS(t *testing.T) {
	result := matchNS("abc.abc", "abc.abc")
	if !result {
		t.Error("matchNS hash bug")
	}
	result = matchNS("abc.abc", "abc.*")
	if !result {
		t.Error("matchNS hash bug")
	}
}

func Test_needSync(t *testing.T) {
	ctx := &SyncCtx{
		IncludeNS: []string{"adc.qq"},
	}
	ctx1 := &SyncCtx{
		ExcludeNS: []string{"adc.cc"},
	}
	ctx2 := &SyncCtx{
		IncludeNS: []string{},
		ExcludeNS: []string{},
	}
	ctx3 := &SyncCtx{
		IncludeNS:[]string{"abc.*"},
	}
	oplog := Oplog{Ns: "adc.qq"}
	oplog1 := Oplog{Ns: "abc.cc"}
	need := needSync(ctx, oplog)
	if !need {
		t.Error("needSync hash bug 1")
	}
	need = needSync(ctx, oplog1)
	if need {
		t.Error("needSync hash bug 2")
	}
	need = needSync(ctx1, oplog)
	if !need {
		t.Error("needSync hash bug 3")
	}
	need = needSync(ctx1, oplog1)
	if !need {
		t.Error("needSync hash bug 4")
	}
	need = needSync(ctx2, oplog)
	if !need {
		t.Error("needSync hash bug 5")
	}
	need = needSync(ctx2, oplog1)
	if !need {
		t.Error("needSync hash bug 6")
	}
	need = needSync(ctx3, oplog)
	if need {
		t.Error("needSync hash bug 7")
	}
	need = needSync(ctx3, oplog1)
	if !need {
		t.Error("needSync hash bug 8")
	}
}
