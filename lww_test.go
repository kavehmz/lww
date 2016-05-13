package lww

import (
	"testing"
	"time"
)

func TestLWW_Init(t *testing.T) {
	lww := LWW{}
	lww.Init()

	if lww.add == nil || lww.remove == nil {
		t.Error("LWW is not initialized correctly", lww.add, lww.remove)
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

	lww.Add(e, ts)
	if ts0, ok := lww.add.get(e); !ok || ts0 != ts {
		t.Error("Add failed", ok, ts0)
	}

	ts = ts.Add(time.Second)
	lww.Add(e, ts)
	if ts0, ok := lww.add.get(e); !ok || ts0 != ts {
		t.Error("Add failed to update the element with newer timestamp", ok, ts0)
	}

	tsOld := ts.Add(-10 * time.Second)
	lww.Add(e, tsOld)
	if ts0, ok := lww.add.get(e); !ok || ts0 != ts {
		t.Error("Add failed to ignore the element with older timestamp", ok, ts0)
	}

	lww.Remove(e, ts)
	if ts0, ok := lww.remove.get(e); !ok || ts0 != ts {
		t.Error("Remove failed to add a new element", ok, ts0)
	}

	ts = ts.Add(time.Second)
	lww.Remove(e, ts)
	if ts0, ok := lww.remove.get(e); !ok || ts0 != ts {
		t.Error("Remove failed to update the element with newer timestamp", ok, ts0)
	}

	tsOld = ts.Add(-10 * time.Second)
	lww.Remove(e, tsOld)
	if ts0, ok := lww.remove.get(e); !ok || ts0 != ts {
		t.Error("Remove failed to ignore the element with older timestamp", ok, ts0)
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
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Remove(i, time.Now())
	}
}

func ExampleLWW() {
	l := LWW{}
	l.Init()
	e := "Any_Structure"
	l.Add(e, time.Now().UTC())
	l.Remove(e, time.Now().UTC().Add(time.Second))
}
