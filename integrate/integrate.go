/*
Package integrate defines expectations which LWW has from an underlying set.
It has a fucntion IntegrationTest that can be used by underlying sets which
implentent TimedSet to see if they are implementing a correct behaviour.

You need to create a test and pass your set to IntegrationTest as shown in the example.
*/
package integrate

import (
	"testing"
	"time"

	"github.com/kavehmz/crdt"
)

// IntegrationTest is a common set of test for underlyings to call.
// Passing this means LWW basic expectatinos from underlying sets is satisfied.
func IntegrationTest(add lww.TimedSet, remove lww.TimedSet, t *testing.T) {
	lww := lww.LWW{AddSet: add, RemoveSet: remove}
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
