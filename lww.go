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
	len()
	get(Element)
	set(Element, time.Time)
	list()
}

// Element define a set member. To make it possible to almost any type of data Element is defined as an empty interface.
// This means for if the element gets saved in the set and then retrieved, it needs type assertion.
//  e := w.Get()
//  client := e.(ClientType)
//  fmt.Println(client.name)
// Note: Element type must be usable as a hash key. Any comparable type can be used.
type Element interface{}
