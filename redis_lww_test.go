package lww

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
)

func TestRedisSet_integration(t *testing.T) {
	var ac redis.Conn
	var rc redis.Conn
	add := setupSet(t, &ac, "TESTADD")
	remove := setupSet(t, &rc, "TESTREMOVE")

	lww := LWW{AddSet: &add, RemoveSet: &remove}
	lww.Init()
	e := "e1"
	ts := time.Now()

	if lww.Exists(e) {
		t.Error("New LWW claims to containt an element")
	}

	lww.Add(e, ts)
	if !lww.Exists(e) {
		t.Error("Newly added element does not exists and it should")
	}

	ts = ts.Add(time.Second)
	lww.Remove(e, ts)
	if lww.Exists(e) {
		t.Error("An element which was remove with a more recent timestmap must be removed and is not")
	}

	ts = ts.Add(time.Second)
	lww.Add(e, ts)
	if !lww.Exists(e) {
		t.Error("An element which was remove and added again with a more recent timestamp does not exists")
	}
}
