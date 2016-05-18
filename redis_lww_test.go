package lww

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

func TestRedisSet_integration(t *testing.T) {
	var ac redis.Conn
	var rc redis.Conn
	add := setupSet(t, &ac, "TESTADD")
	remove := setupSet(t, &rc, "TESTREMOVE")

	IntegrationTest(&add, &remove, t)
}
