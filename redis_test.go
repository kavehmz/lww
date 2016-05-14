package lww

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
)

func TestRedis_init(t *testing.T) {
	r, _ := redis.Dial("tcp", ":6379")
	s := RedisSet{}
	s.init()
	if s.LastState == nil {
		t.Error("No error for missing params")
	}

	s = RedisSet{Conn: r}
	s.init()
	if s.LastState == nil {
		t.Error("No error for missing params")
	}

	s = RedisSet{Conn: r, Marshal: func(e Element) string { return e.(string) }}
	s.init()
	if s.LastState == nil {
		t.Error("No error for missing params")
	}

	s = RedisSet{Conn: r, Marshal: func(e Element) string { return e.(string) }, UnMarshal: func(e string) Element { return e }}
	s.init()
	if s.LastState == nil {
		t.Error("No error for missing params")
	}

	s = RedisSet{Conn: r, Marshal: func(e Element) string { return e.(string) }, UnMarshal: func(e string) Element { return e }, SetKey: "TESTKEY"}
	s.init()
	if s.LastState != nil {
		t.Error("Error raised when all params are present and correct")
	}
}

func TestRedis(t *testing.T) {
	r, err := redis.Dial("tcp", "localhost:6379")
	defer r.Close()
	if err != nil {
		t.Error("Can't setup redis for tests", err)
	}
	s := RedisSet{Conn: r, Marshal: func(e Element) string { return e.(string) }, UnMarshal: func(e string) Element { return e }, SetKey: "TESTKEY"}

	if s.len() != 0 {
		t.Error("New set if not empty")
	}

	a := "data"
	ts := time.Now().Round(time.Microsecond)
	s.set(a, ts)
	if s.len() != 1 {
		t.Error("Adding element to set failed")
	}
	if ts0, ok := s.get(a); !ok || ts0 != ts {
		t.Error("Element is not saved corretly", ts0, ok, ts)
	}

	ts = ts.Add(time.Second * 10)
	s.set(a, ts)
	if ts0, ok := s.get(a); !ok || ts0 != ts {
		t.Error("Element is not updated corretly")
	}

	s.set("new data", ts)
	if ts0, ok := s.get(a); !ok || ts0 != ts {
		t.Error("New Element is not added corretly")
	}

	l := s.list()
	if len(l) != 2 {
		t.Error("List is not returning all elements of the set correctly")
	}
	if l[0] != "data" || l[1] != "new data" {
		t.Error("List elements are not correct")
	}
}
