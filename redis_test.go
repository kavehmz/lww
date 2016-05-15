package lww

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
)

func TestRedisSet_init(t *testing.T) {
	r, _ := redis.Dial("tcp", ":6379")
	s := RedisSet{}
	s.Init()
	if s.LastState == nil {
		t.Error("No error for missing params")
	}

	s = RedisSet{Conn: r}
	s.Init()
	if s.LastState == nil {
		t.Error("No error for missing params")
	}

	s = RedisSet{Conn: r, Marshal: func(e Element) string { return e.(string) }}
	s.Init()
	if s.LastState == nil {
		t.Error("No error for missing params")
	}

	s = RedisSet{Conn: r, Marshal: func(e Element) string { return e.(string) }, UnMarshal: func(e string) Element { return e }}
	s.Init()
	if s.LastState == nil {
		t.Error("No error for missing params")
	}

	s = RedisSet{Conn: r, Marshal: func(e Element) string { return e.(string) }, UnMarshal: func(e string) Element { return e }, SetKey: "TESTKEY"}
	s.Init()
	if s.LastState != nil {
		t.Error("Error raised when all params are present and correct")
	}
}

func setupSet(t interface {
	Error(...interface{})
}, r *redis.Conn, key string) RedisSet {
	c, err := redis.Dial("tcp", "localhost:6379")
	r = &c
	if err != nil {
		t.Error("Can't setup redis for tests", err)
	}
	s := RedisSet{Conn: *r, Marshal: func(e Element) string { return e.(string) }, UnMarshal: func(e string) Element { return e }, SetKey: key}
	s.Init()
	return s
}

func TestRedisSet(t *testing.T) {
	var r *redis.Conn
	s := setupSet(t, r, "TESTKEY")

	if s.Len() != 0 {
		t.Error("New set if not empty")
	}

	a := "data"
	ts := time.Now().Round(time.Microsecond)
	s.Set(a, ts)
	if s.Len() != 1 {
		t.Error("Adding element to set failed")
	}
	if ts0, ok := s.Get(a); !ok || ts0 != ts {
		t.Error("Element is not saved corretly", ts0, ok, ts)
	}

	ts = ts.Add(time.Second * 10)
	s.Set(a, ts)
	if ts0, ok := s.Get(a); !ok || ts0 != ts {
		t.Error("Element is not updated corretly")
	}

	s.Set("new data", ts)
	if ts0, ok := s.Get(a); !ok || ts0 != ts {
		t.Error("New Element is not added corretly")
	}

	l := s.List()
	if len(l) != 2 {
		t.Error("List is not returning all elements of the set correctly")
	}
	if l[0] != "data" || l[1] != "new data" {
		t.Error("List elements are not correct")
	}
}

func BenchmarkRedisSet_add_different(b *testing.B) {
	var r *redis.Conn
	s := setupSet(b, r, "TESTKEY")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Set(strconv.Itoa(i), time.Now())
	}
}

func BenchmarkRedisSet_add_same(b *testing.B) {
	var r *redis.Conn
	s := setupSet(b, r, "TESTKEY")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Set("test", time.Now().Add(time.Duration(i)*time.Microsecond))
	}
}

func BenchmarkRedisSet_get(b *testing.B) {
	var r *redis.Conn
	s := setupSet(b, r, "TESTKEY")
	for i := 0; i < b.N; i++ {
		s.Set(strconv.Itoa(i), time.Now())
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Get(strconv.Itoa(i))
	}
}

func ExampleRedisSet() {
	c, _ := redis.Dial("tcp", "localhost:6379")
	s := RedisSet{Conn: c, Marshal: func(e Element) string { return e.(string) }, UnMarshal: func(e string) Element { return e }, SetKey: "TESTKEY"}
	s.Init()
	s.Set("Data", time.Unix(1451606400, 0))
	ts, ok := s.Get("Data")
	fmt.Println(ok)
	fmt.Println(ts.Unix())
	// Output:
	// true
	// 1451606400
}
