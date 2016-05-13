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
*/
type Set struct {
	members map[Element]time.Time
	sync.RWMutex
}

func (s *Set) init() {
	s.Lock()
	defer s.Unlock()
	s.members = make(map[Element]time.Time)
}

func (s *Set) set(e Element, t time.Time) {
	s.Lock()
	s.members[e] = t
	s.Unlock()
}

func (s *Set) len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.members)
}

func (s *Set) get(e Element) (time.Time, bool) {
	s.RLock()
	defer s.RUnlock()
	val, ok := s.members[e]
	return val, ok
}

func (s *Set) list() []Element {
	s.RLock()
	defer s.RUnlock()
	l := make([]Element, 0, s.len())
	for k := range s.members {
		l = append(l, k)
	}
	return l
}
