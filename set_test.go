package lww

import (
	"reflect"
	"testing"
	"time"
)

func TestSet_init(t *testing.T) {
	s := Set{}
	s.init()
	if s.members == nil {
		t.Error("Set is not initialized correctly")
	}
}

type customType struct {
	name string
	age  int
}

func TestSet_add(t *testing.T) {
	s := Set{}
	s.init()
	a := customType{name: "John", age: 18}

	ts := time.Now()
	s.set(a, ts)
	if v, ok := s.members[a]; !ok || v != ts {
		t.Error("Element was not added correctly", ok, v, ts)
	}

	ts = ts.Add(time.Second * 10)
	s.set(a, ts)
	if v, ok := s.members[a]; !ok || v != ts {
		t.Error("Element was not changed correctly if timestamp is different", ok, v, ts)
	}
}

func TestSet_len(t *testing.T) {
	s := Set{}
	s.init()

	if s.len() != 0 {
		t.Error("len is wrong after init")
	}
	s.set(customType{name: "John", age: 18}, time.Now())
	s.set(customType{name: "Frank", age: 20}, time.Now())

	if s.len() != 2 {
		t.Error("len is wrong after add", s.len())
	}
}

func TestSet_get(t *testing.T) {
	s := Set{}
	s.init()

	if _, ok := s.get(customType{name: "John", age: 18}); ok {
		t.Error("After init get is finding elements")
	}

	ts := time.Now()
	s.set(customType{name: "John", age: 18}, ts)
	if v, ok := s.get(customType{name: "John", age: 18}); !ok || v != ts {
		t.Error("get is wrong after add", ok, v, ts)
	}
}

func TestSet_list(t *testing.T) {
	s := Set{}
	s.init()
	l := []customType{{name: "John", age: 18}, {name: "Frank", age: 20}}
	s.set(l[0], time.Now())
	s.set(l[1], time.Now())

	a := s.list()
	if reflect.DeepEqual(l, a) {
		t.Error("list did not return correct memeber", a)
	}
}

func BenchmarkSet_add_different(b *testing.B) {
	s := Set{}
	s.init()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.set(i, time.Now())
	}
}

func BenchmarkSet_add_same(b *testing.B) {
	s := Set{}
	s.init()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.set("same_element", time.Now())
	}
}

func BenchmarkSet_get(b *testing.B) {
	s := Set{}
	s.init()
	for i := 0; i < b.N; i++ {
		s.set(1, time.Now())
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.get(i)
	}
}
