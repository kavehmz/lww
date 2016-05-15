package lww

import (
	"testing"
	"time"
)

func TestSet_init(t *testing.T) {
	s := Set{}
	s.Init()
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
	s.Init()
	a := customType{name: "John", age: 18}

	ts := time.Now()
	s.Set(a, ts)
	if v, ok := s.members[a]; !ok || v != ts {
		t.Error("Element was not added correctly", ok, v, ts)
	}

	ts = ts.Add(time.Second * 10)
	s.Set(a, ts)
	if v, ok := s.members[a]; !ok || v != ts {
		t.Error("Element was not changed correctly if timestamp is different", ok, v, ts)
	}
}

func TestSet_len(t *testing.T) {
	s := Set{}
	s.Init()

	if s.Len() != 0 {
		t.Error("len is wrong after init")
	}
	s.Set(customType{name: "John", age: 18}, time.Now())
	s.Set(customType{name: "Frank", age: 20}, time.Now())

	if s.Len() != 2 {
		t.Error("len is wrong after add", s.Len())
	}
}

func TestSet_get(t *testing.T) {
	s := Set{}
	s.Init()

	if _, ok := s.Get(customType{name: "John", age: 18}); ok {
		t.Error("After init get is finding elements")
	}

	ts := time.Now()
	s.Set(customType{name: "John", age: 18}, ts)
	if v, ok := s.Get(customType{name: "John", age: 18}); !ok || v != ts {
		t.Error("get is wrong after add", ok, v, ts)
	}
}

func TestSet_list(t *testing.T) {
	s := Set{}
	s.Init()
	l := [2]customType{{name: "John", age: 18}, {name: "Frank", age: 20}}
	s.Set(l[0], time.Now())
	s.Set(l[1], time.Now())

	a := s.List()
	if a[0].(customType) != l[0] && a[0].(customType) != l[1] {
		t.Error("list did not return correct memeber", a[0].(customType), l[0], l[1])
	}
	if a[1].(customType) != l[0] && a[1].(customType) != l[1] {
		t.Error("list did not return correct memeber", a[1].(customType), l[0], l[1])
	}
}

func BenchmarkSet_add_different(b *testing.B) {
	s := Set{}
	s.Init()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Set(i, time.Now())
	}
}

func BenchmarkSet_add_same(b *testing.B) {
	s := Set{}
	s.Init()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Set("same_element", time.Now())
	}
}

func BenchmarkSet_get(b *testing.B) {
	s := Set{}
	s.Init()
	for i := 0; i < b.N; i++ {
		s.Set(i, time.Now())
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Get(i)
	}
}
