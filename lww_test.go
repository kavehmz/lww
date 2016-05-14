package lww

import (
	"fmt"
	"testing"
	"time"
)

func TestLWW_Init(t *testing.T) {
	lww := LWW{}
	lww.Init()

	if lww.AddSet == nil || lww.RemoveSet == nil {
		t.Error("LWW is not initialized correctly", lww.AddSet, lww.RemoveSet)
	}
}

type element struct {
	name string
	age  int
}

func TestLWW_AddExistRemove(t *testing.T) {
	lww := LWW{}
	lww.Init()
	e := element{name: "John", age: 18}
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

func TestLWW_Get(t *testing.T) {
	lww := LWW{}
	lww.Init()
	l := []customType{{name: "John", age: 18}, {name: "Betty", age: 22}}
	lr := customType{name: "Frank", age: 20}
	lww.Add(l[0], time.Now())
	lww.Add(l[1], time.Now())
	lww.Add(lr, time.Now())
	lww.Remove(lr, time.Now().Add(time.Second))

	a := lww.Get()
	if a[0].(customType) != l[0] && a[0].(customType) != l[1] {
		t.Error("list did not return correct memeber", a[0].(customType), l[0], l[1])
	}
	if a[1].(customType) != l[0] && a[1].(customType) != l[1] {
		t.Error("list did not return correct memeber", a[1].(customType), l[0], l[1])
	}

}

func BenchmarkLWW_Add_differnt(b *testing.B) {
	l := LWW{}
	l.Init()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Add(i, time.Now())
	}
}

func BenchmarkLWW_Add_same(b *testing.B) {
	l := LWW{}
	l.Init()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Add("same", time.Now())
	}
}

func BenchmarkLWW_Remove(b *testing.B) {
	l := LWW{}
	l.Init()
	for i := 0; i < b.N; i++ {
		l.Add(i, time.Now())
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Remove(i, time.Now())
	}
}

func BenchmarkLWW_Exists(b *testing.B) {
	l := LWW{}
	l.Init()
	for i := 0; i < b.N; i++ {
		l.Add(i, time.Now())
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Exists(i)
	}
}

func ExampleLWW() {
	l := LWW{}
	l.Init()
	e := "Any_Structure"
	l.Add(e, time.Now().UTC())
	l.Remove(e, time.Now().UTC().Add(time.Second))
	fmt.Println(l.Exists(e))
	l.Add(e, time.Now().UTC().Add(2*time.Second))
	fmt.Println(l.Exists(e))
	// Output:
	// false
	// true
}
