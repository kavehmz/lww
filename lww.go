/*
Package lww implements a Last-Writer-Wins (LWW) Element Set data structure.

In distributed computing, a conflict-free replicated data type (CRDT) is a type of specially-designed data structure used to achieve strong eventual consistency (SEC) and monotonicity (absence of rollbacks).

One type of data structure used in implementing CRDT is LWW-element-set.

LWW-element-set is a set that its elements have timestamp. Add and remove will save the timestamp along with data in two different sets for each element.

Queries over LWW-set will check both add and remove timestamps to decide about state of each element is being existed to removed from the list.

LWW

lww package implements LWW data structure in a modular way. It defines a TimedSet interface for underlying storage.

lww package includes two storage underlying.

Set

Set is one implementation of TimedSet. It uses Go maps to store data. It is a fast but volatile implementation.

Maps in theory have worse Big O of O(n) for different operations, but in practice they are almost reliable for O(1) as long as hash function and hash table implementations are good enough.

Set is the default underlying for LWW if no other TimedSet are attached to AddSet or RemoveSet.

  # This will use Set as its AddSet and RemoveSet
  lww := LWW{}

Maps are by nature vulnerable to concurrent access. To avoid race problems Set uses a sync.RWMutex as its locking mechanism.

RedisSet

RedisSet is another implementation of TimedSet included in lww package. It uses Redis Sorted Sets to store data.

Redis nature of atomic operations makes it immune to race problem and there is no need to any extra lock mechanism. But it introduces other complexities.

To keep the lww simple, handling of Redis connection for both AddSet and RemoveSet in case of RedisSet is passed to client.
It is practical as Redis setup can vary based on application and client might want handle complex connection handling.

Adding New underlying

To add a new underlying you need to implement the necessary methods in your structure. They are defined in TimedSet interface.

Assuming you do that and they work as expected you can initialize LWW like:

  add    := MyUnderlying{param: "for_add"}
  remove := MyUnderlying{param: "for_remove"}
  lww    := LWW{AddSet:add, RemoveSet:remove}

Note that in theory AddSet and RemoveSet can have different underlying attached.
This might be useful in applications which can predict higher magnitude of Adds compared to Removes. In that case application can implementation different types of TimedSet to optimize the setup

*/
package lww

import (
	"testing"
	"time"
)

// TimedSet interface defines what is required for an underlying set for WWL.
type TimedSet interface {
	//Init will do a one time setup for underlying set. It will be called from WLL.Init
	Init()
	//Len must return the number of members in the set
	Len() int
	//Get returns timestmap of the element in the set if it exists and true. Otherwise it will return an empty timestamp and false.
	Get(interface{}) (time.Time, bool)
	//Set adds an element to the set if it does not exists. It it exists Set will update the provided timestamp.
	Set(interface{}, time.Time)
	//List returns list of all elements in the set
	List() []interface{}
}

// LWW type a Last-Writer-Wins (LWW) Element Set data structure.
type LWW struct {
	// AddSet will store the state of elements added to the set. By default it is will be of type lww.Set.
	AddSet TimedSet
	// AddSet will store the state of elements removed from the set. By default it is will be of type lww.Set
	RemoveSet TimedSet
}

// Init will initialize the underlying sets required for LWW.
// Internally it works on two sets named "add" and "remove".
func (lww *LWW) Init() {
	if lww.AddSet == nil {
		lww.AddSet = &Set{}
	}
	if lww.RemoveSet == nil {
		lww.RemoveSet = &Set{}
	}
	lww.AddSet.Init()
	lww.RemoveSet.Init()
}

// Add will add an element to the add-set if it does not exists and updates its timestamp to
// great one between current one and new one.
func (lww *LWW) Add(e interface{}, t time.Time) {
	lww.AddSet.Set(e, t)
}

// Remove will add an element to the remove-set if it does not exists and updates its timestamp to
// great one between current one and new one.
func (lww *LWW) Remove(e interface{}, t time.Time) {
	if val, ok := lww.RemoveSet.Get(e); !ok || t.UnixNano() > val.UnixNano() {
		lww.RemoveSet.Set(e, t)
	}
}

// Exists returns true if element has a more recent record in add-set than in remove-set
func (lww *LWW) Exists(e interface{}) bool {
	a, aok := lww.AddSet.Get(e)
	r, rok := lww.RemoveSet.Get(e)
	if !rok {
		return aok
	}
	return a.UnixNano() > r.UnixNano()
}

// Get returns slice of elements that "Exist".
func (lww *LWW) Get() []interface{} {

	l := make([]interface{}, 0, lww.AddSet.Len())
	for _, e := range lww.AddSet.List() {
		if lww.Exists(e) {
			l = append(l, e)
		}
	}
	return l
}

// IntegrationTest is a common set of test for underlyings to call.
// Passing this means LWW basic expectatinos from underlying sets is satisfied.
func IntegrationTest(add TimedSet, remove TimedSet, t *testing.T) {
	lww := LWW{AddSet: add, RemoveSet: remove}
	lww.Init()
	e := "e1"
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
