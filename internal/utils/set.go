package utils

type Set[T comparable] map[T]struct{}

func (s *Set[T]) Reset() {
	clear(*s)
}

func (s *Set[T]) Add(v ...T) {
	for _, val := range v {
		(*s)[val] = struct{}{}
	}
}

func (s *Set[T]) Delete(v ...T) {
	for _, val := range v {
		delete(*s, val)
	}
}

// Get 不保证顺序
func (s *Set[T]) Get() []T {
	ret := make([]T, len(*s))
	i := 0
	for val := range *s {
		ret[i] = val
		i++
	}
	return ret
}

func (s *Set[T]) Ok(v T) bool {
	_, ok := (*s)[v]
	return ok
}

// Clean 能保证顺序
func (s *Set[T]) Clean(ts []T) []T {
	write := 0
	for read := range ts {
		if !s.Ok(ts[read]) {
			(*s)[ts[read]] = struct{}{}
			ts[write] = ts[read]
			write++
		}
	}
	return ts[:write]
}
