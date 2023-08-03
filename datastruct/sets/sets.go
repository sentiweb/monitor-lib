// Very simple implementation of Set 
// Mostly used for checking if an item already exists in the set or not.
// For more complete implementation more complete package should be used (like gods)
package sets

type Set[T comparable] struct {
	items map[T]struct{}
}

var itemExists = struct{}{}

func New[T comparable](size int) *Set[T] {
	return &Set[T]{
		items: make(map[T]struct{}, size),
	}
}

// Add an item to the set, returns false if not added (already exists), true if new item
func (s *Set[T]) Add(item T) bool {
	_, ok := s.items[item]
	if(ok) {
		return false
	}
	s.items[item] = itemExists
	return true
}

func (s *Set[T]) Has(item T) bool {
	_, ok := s.items[item]
	return ok
}

func (s *Set[T]) Remove(item T) bool {
	_, ok := s.items[item]
	if(!ok) {
		return false
	}
	delete(s.items, item)
	return true
}

func (s *Set[T]) Empty() bool {
	return len(s.items) == 0
}

func (s *Set[T]) Size() int {
	return len(s.items)
}

func (s *Set[T]) Values() []T {
	v := make([]T, 0, len(s.items))
	for k := range s.items {
		v = append(v, k)
	}
	return v
}

func (s *Set[T]) AddAll(itemsToAdd []T) int  {
	added := 0
	for _, item := range itemsToAdd {
		s.Add(item)
		added += 1
	}
	return added
}

func (s *Set[T]) HasAny(itemArray []T) bool  {
	if(len(itemArray) == 0) {
		return false
	}
	for _, item := range itemArray {
		if(s.Has(item)) {
			return true
		}
	}
	return false
}

func (s *Set[T]) HasAll(itemArray []T) bool  {
	if len(s.items) == 0 {
		return false
	}
	for _, item := range itemArray {
		if(!s.Has(item)) {
			return false
		}
	}
	return true
}
