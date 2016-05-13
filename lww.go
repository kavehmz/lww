/*
Package lww implements a Last-Writer-Wins (LWW) Element Set data structure.

In distributed computing, a conflict-free replicated data type (CRDT) is a type of specially-designed data structure used to achieve strong eventual consistency (SEC) and monotonicity (absence of rollbacks).

One type of data structure used in implementing CRDT is LWW-element-set.

LWW-element-set is a set that its elements have timestamp. Add and remove will save the timestamp along with data in two different sets for each element.

Queries over LWW-set will check both add and remove timestamps to decide about state of each element is being existed to removed from the list.
*/
package lww

import "time"

// TimedSet interface defines what is required for an underlying set for WWL.
type TimedSet interface {
	init()
	len() int
	get(Element) (time.Time, bool)
	set(Element, time.Time)
	list() []Element
}

// Element define a set member. To make it possible to almost any type of data Element is defined as an empty interface.
// This means for if the element gets saved in the set and then retrieved, it needs type assertion.
//  e := w.Get()
//  client := e.(ClientType)
//  fmt.Println(client.name)
// Note: Element type must be usable as a hash key. Any comparable type can be used.
type Element interface{}

// LWW type a Last-Writer-Wins (LWW) Element Set data structure.
type LWW struct {
	add    TimedSet
	remove TimedSet
}

// Init will initialize the underlying sets required for LWW.
// Internally it works on two sets named "add" and "remove".
func (lww *LWW) Init() {
	lww.add = &Set{}
	lww.remove = &Set{}
	lww.add.init()
	lww.remove.init()
}

// Add will add an Element to the add-set if it does not exists and updates its timestamp to
// great one between current one and new one.
func (lww *LWW) Add(e Element, t time.Time) {
	if val, ok := lww.add.get(e); !ok || t.UnixNano() > val.UnixNano() {
		lww.add.set(e, t)
	}
}

// Remove will add an Element to the remove-set if it does not exists and updates its timestamp to
// great one between current one and new one.
func (lww *LWW) Remove(e Element, t time.Time) {
	if val, ok := lww.remove.get(e); !ok || t.UnixNano() > val.UnixNano() {
		lww.remove.set(e, t)
	}
}

// Exists returns true if Element has a more recent record in add-set than in remove-set
func (lww *LWW) Exists(e Element) bool {
	a, aok := lww.add.get(e)
	r, rok := lww.remove.get(e)
	if !rok {
		return aok
	}
	return a.UnixNano() > r.UnixNano()
}

// Get returns slice of Elements that "Exist".
func (lww *LWW) Get() []Element {

	l := make([]Element, 0, lww.add.len())
	for _, e := range lww.add.list() {
		if lww.Exists(e) {
			l = append(l, e)
		}
	}
	return l
}
