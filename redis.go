package lww

import (
	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
)

/*RedisSet is a race free implementation of what WWL can use as udnerlying set.
This implementation uses redis ZSET.
ZSET in redis uses scores to sort the elements. Score is a IEEE 754 floating point number,
that is able to represent precisely integer numbers between -(2^53) and +(2^53) included.
That is between -9007199254740992 and 9007199254740992.
This will limit this sets precision to save element's action timestamp to 1 milli-seconds.
Notice that time.Time precision is 1 nano-seconds by defaults. For this lack of precision all
timestamps are rounded to nearest microsecond.
Using redis can also cause latency cause by network or socket communication.
*/
type RedisSet struct {
	// Conn is the redis connection to be used.
	Conn redis.Conn
	// AddSet sets which key will be used in redis for the set.
	SetKey string
	// Marshal function needs to convert the Element to string. Redis can only store and retrieve string values.
	Marshal func(Element) string
	// UnMarshal function needs to be able to convert a Marshalled string back to a readable structure for consumer of library.
	UnMarshal func(string) Element
	// LastState is an error type that will return the error state of last executed redis command. Add redis connection are not shareable this can be used after each command to know the last state.
	LastState error
}

func roundToMicro(t time.Time) int64 {
	return t.Round(time.Microsecond).UnixNano() / 1000
}

func (s *RedisSet) checkErr(err error) {
	if err != nil {
		s.LastState = err
		return
	}
	s.LastState = nil
}

func (s *RedisSet) init() {
	if s.Conn == nil {
		s.checkErr(errors.New("Conn must be set"))
		return
	}
	if s.Marshal == nil {
		s.checkErr(errors.New("Marshal must be set"))
		return
	}
	if s.UnMarshal == nil {
		s.checkErr(errors.New("UnMarshal must be set"))
		return
	}
	if s.SetKey == "" {
		s.checkErr(errors.New("SetKey must be set"))
		return
	}

	_, err := s.Conn.Do("DEL", s.SetKey)
	s.checkErr(err)

}

func (s *RedisSet) set(e Element, t time.Time) {
	_, err := s.Conn.Do("ZADD", s.SetKey, roundToMicro(t), s.Marshal(e))
	s.checkErr(err)
}

func (s *RedisSet) len() int {
	n, err := redis.Int(s.Conn.Do("ZCARD", s.SetKey))
	s.checkErr(err)
	return n
}

func (s *RedisSet) get(e Element) (val time.Time, ok bool) {
	n, err := redis.Int(s.Conn.Do("ZSCORE", s.SetKey, s.Marshal(e)))
	s.checkErr(err)
	if err == nil {
		ok = true
		val = time.Unix(0, 0).Add(time.Duration(n) * time.Microsecond)
	}
	return val, ok
}

func (s *RedisSet) list() []Element {
	var l []Element
	zs, err := redis.Strings(s.Conn.Do("ZRANGE", s.SetKey, 0, -1))
	s.checkErr(err)
	for _, v := range zs {
		l = append(l, s.UnMarshal(v))
	}
	return l
}
