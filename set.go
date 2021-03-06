package lww

import (
	"sync"
	"time"
)

/*Set is a race free implementation of what WWL can use as udnerlying set.
This implementation uses maps. To avoid race condition that comes by using maps
it is using a locking mechanism. Set is using separete Read/Write locks.
Map data structure have a practical performance of O(1) but locking instructions might make
this implementation sub optimal for write heavy solutions.

Note: Elements of set type must be usable as a hash key. Any comparable in Go type can be used.
*/
type Set struct {
	members map[interface{}]time.Time
	sync.RWMutex
}

//Init will do a one time setup for underlying set. It will be called from WLL.Init
func (s *Set) Init() {
	s.Lock()
	defer s.Unlock()
	s.members = make(map[interface{}]time.Time)
}

//Set adds an element to the set if it does not exists. It it exists Set will update the provided timestamp.
func (s *Set) Set(e interface{}, t time.Time) {
	s.Lock()
	if val, ok := s.members[e]; !ok || t.UnixNano() > val.UnixNano() {
		s.members[e] = t
	}
	s.Unlock()
}

//Len must return the number of members in the set
func (s *Set) Len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.members)
}

//Get returns timestmap of the element in the set if it exists and true. Otherwise it will return an empty timestamp and false.
func (s *Set) Get(e interface{}) (time.Time, bool) {
	s.RLock()
	defer s.RUnlock()
	val, ok := s.members[e]
	return val, ok
}

//List returns list of all elements in the set
func (s *Set) List() []interface{} {
	s.RLock()
	defer s.RUnlock()
	l := make([]interface{}, 0, s.Len())
	for k := range s.members {
		l = append(l, k)
	}
	return l
}
