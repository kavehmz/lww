package integrate

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/kavehmz/lww"
)

func setupSet(t interface {
	Error(...interface{})
}, r *redis.Conn, key string) lww.RedisSet {
	c, _ := redis.Dial("tcp", "localhost:6379")
	_, err := c.Do("DEL", key)
	r = &c
	if err != nil {
		t.Error("Can't setup redis for tests", err)
	}
	s := lww.RedisSet{Conn: *r, Marshal: func(e interface{}) string { return e.(string) }, UnMarshal: func(e string) interface{} { return e }, SetKey: key}
	s.Init()
	return s
}

func TestRedisSet_integration(t *testing.T) {
	var ac redis.Conn
	var rc redis.Conn
	add := setupSet(t, &ac, "TESTADD")
	remove := setupSet(t, &rc, "TESTREMOVE")

	IntegrationTest(&add, &remove, t)
}

func Example() {
	var ac redis.Conn
	var rc redis.Conn
	var t testing.T
	add := setupSet(&t, &ac, "TESTADD")
	remove := setupSet(&t, &rc, "TESTREMOVE")

	IntegrationTest(&add, &remove, &t)
}
