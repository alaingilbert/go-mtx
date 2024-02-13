package mtx

import (
	"cmp"
	"sync"
)

var debug = false

func toPtr[T any](v T) *T { return &v }

type Locker[T any] interface {
	sync.Locker
	Get() T
	RWith(clb func(v T))
	RWithE(clb func(v T) error) error
	Replace(newVal T) (old T)
	Set(v T)
	Val() *T
	With(clb func(v *T))
	WithE(clb func(v *T) error) error
}

// Compile time checks to ensure type satisfies Locker interface
var _ Locker[any] = (*base[sync.Locker, any])(nil)

type base[M sync.Locker, T any] struct {
	m M
	v T
}

func newBase[M sync.Locker, T any](m M, v T) *base[M, T] {
	return &base[M, T]{m: m, v: v}
}

// Lock exposes the underlying sync.Mutex Lock function
func (m *base[M, T]) Lock() { m.m.Lock() }

// Unlock exposes the underlying sync.Mutex Unlock function
func (m *base[M, T]) Unlock() { m.m.Unlock() }

// Val gets the wrapped value by the mutex.
// WARNING: the caller must make sure the code that uses it is thread-safe
func (m *base[M, T]) Val() *T {
	return &m.v
}

// WithE provide a callback scope where the wrapped value can be safely used
func (m *base[M, T]) WithE(clb func(v *T) error) error {
	m.Lock()
	defer m.Unlock()
	return clb(&m.v)
}

// With same as WithE but do return an error
func (m *base[M, T]) With(clb func(v *T)) {
	_ = m.WithE(func(tx *T) error {
		clb(tx)
		return nil
	})
}

// RWithE provide a callback scope where the wrapped value can be safely used for Read only purposes
func (m *base[M, T]) RWithE(clb func(v T) error) error {
	if debug {
		println("base RWithE")
	}
	return m.WithE(func(v *T) error {
		return clb(*v)
	})
}

// RWith same as RWithE but do not return an error
func (m *base[M, T]) RWith(clb func(v T)) {
	_ = m.RWithE(func(tx T) error {
		clb(tx)
		return nil
	})
}

// Get safely gets the wrapped value
func (m *base[M, T]) Get() (out T) {
	m.RWith(func(v T) { out = v })
	return out
}

// Set a new value
func (m *base[M, T]) Set(newV T) {
	m.With(func(v *T) { *v = newV })
}

// Replace set a new value and return the old value
func (m *base[M, T]) Replace(newVal T) (old T) {
	m.With(func(v *T) {
		old = *v
		*v = newVal
	})
	return
}

//----------------------

// Mtx generic helper for sync.Mutex
type Mtx[T any] struct {
	*base[*sync.Mutex, T]
}

// NewMtx creates a new Mtx
func NewMtx[T any](v T) Mtx[T] {
	return Mtx[T]{newBase[*sync.Mutex, T](&sync.Mutex{}, v)}
}

// NewMtxPtr creates a new pointer to *Mtx
func NewMtxPtr[T any](v T) *Mtx[T] { return toPtr(NewMtx(v)) }

//----------------------

// RWMtx generic helper for sync.RWMutex
type RWMtx[T any] struct {
	*base[*sync.RWMutex, T]
}

// NewRWMtx creates a new RWMtx
func NewRWMtx[T any](v T) RWMtx[T] {
	return RWMtx[T]{newBase[*sync.RWMutex, T](&sync.RWMutex{}, v)}
}

// NewRWMtxPtr creates a new pointer to *RWMtx
func NewRWMtxPtr[T any](v T) *RWMtx[T] { return toPtr(NewRWMtx(v)) }

// RLock exposes the underlying sync.RWMutex RLock function
func (m *RWMtx[T]) RLock() { m.m.RLock() }

// RUnlock exposes the underlying sync.RWMutex RUnlock function
func (m *RWMtx[T]) RUnlock() { m.m.RUnlock() }

// RWithE provide a callback scope where the wrapped value can be safely used for Read only purposes
func (m *RWMtx[T]) RWithE(clb func(v T) error) error {
	if debug {
		println("RWMtx RWithE")
	}
	m.RLock()
	defer m.RUnlock()
	return clb(m.v)
}

// RWith same as RWithE but do not return an error
func (m *RWMtx[T]) RWith(clb func(v T)) {
	_ = m.RWithE(func(tx T) error {
		clb(tx)
		return nil
	})
}

//----------------------

func defaultMap[K cmp.Ordered, V any](v map[K]V) map[K]V {
	if v == nil {
		v = make(map[K]V)
	}
	return v
}

func NewMap[K cmp.Ordered, V any](v map[K]V) Map[K, V] {
	return Map[K, V]{newBaseMapPtr[K, V](NewMtxPtr(defaultMap(v)))}
}

func NewMapPtr[K cmp.Ordered, V any](v map[K]V) *Map[K, V] { return toPtr(NewMap[K, V](v)) }

//----------------------

func NewRWMap[K cmp.Ordered, V any](v map[K]V) Map[K, V] {
	return Map[K, V]{newBaseMapPtr[K, V](NewRWMtxPtr(defaultMap(v)))}
}

func NewRWMapPtr[K cmp.Ordered, V any](v map[K]V) *Map[K, V] { return toPtr(NewRWMap[K, V](v)) }

//----------------------

type IMap[K cmp.Ordered, V any] interface {
	Locker[map[K]V]
	Clone() (out map[K]V)
	DeleteKey(k K)
	Each(clb func(K, V))
	GetKey(k K) (out V, ok bool)
	HasKey(k K) (found bool)
	Keys() (out []K)
	Len() (out int)
	SetKey(k K, v V)
	TakeKey(k K) (out V, ok bool)
	Values() (out []V)
}

// Compile time checks to ensure type satisfies IMap interface
var _ IMap[int, int] = (*Map[int, int])(nil)

func newBaseMapPtr[K cmp.Ordered, V any](m Locker[map[K]V]) *Map[K, V] {
	return &Map[K, V]{m}
}

type Map[K cmp.Ordered, V any] struct {
	Locker[map[K]V]
}

func (m *Map[K, V]) SetKey(k K, v V) {
	m.With(func(m *map[K]V) { (*m)[k] = v })
}

func (m *Map[K, V]) GetKey(k K) (out V, ok bool) {
	m.RWith(func(mm map[K]V) { out, ok = mm[k] })
	return
}

func (m *Map[K, V]) HasKey(k K) (found bool) {
	m.RWith(func(mm map[K]V) { _, found = mm[k] })
	return
}

func (m *Map[K, V]) TakeKey(k K) (out V, ok bool) {
	m.With(func(m *map[K]V) {
		out, ok = (*m)[k]
		if ok {
			delete(*m, k)
		}
	})
	return
}

func (m *Map[K, V]) DeleteKey(k K) {
	m.With(func(m *map[K]V) { delete(*m, k) })
	return
}

func (m *Map[K, V]) Len() (out int) {
	m.RWith(func(mm map[K]V) { out = len(mm) })
	return
}

func (m *Map[K, V]) Each(clb func(K, V)) {
	m.RWith(func(mm map[K]V) {
		for k, v := range mm {
			clb(k, v)
		}
	})
}

func (m *Map[K, V]) Keys() (out []K) {
	out = make([]K, 0)
	m.RWith(func(mm map[K]V) {
		for k := range mm {
			out = append(out, k)
		}
	})
	return
}

func (m *Map[K, V]) Values() (out []V) {
	out = make([]V, 0)
	m.RWith(func(mm map[K]V) {
		for _, v := range mm {
			out = append(out, v)
		}
	})
	return
}

func (m *Map[K, V]) Clone() (out map[K]V) {
	m.RWith(func(mm map[K]V) {
		out = make(map[K]V, len(mm))
		for k, v := range mm {
			out[k] = v
		}
	})
	return
}

//----------------------

func defaultSlice[T any](v []T) []T {
	if v == nil {
		v = make([]T, 0)
	}
	return v
}

func NewSlice[T any](v []T) Slice[T] {
	return Slice[T]{newBaseSlicePtr[T](NewMtxPtr(defaultSlice(v)))}
}

func NewSlicePtr[T any](v []T) *Slice[T] { return toPtr(NewSlice[T](v)) }

//----------------------

func NewRWSlice[T any](v []T) Slice[T] {
	return Slice[T]{newBaseSlicePtr[T](NewRWMtxPtr(defaultSlice(v)))}
}

func NewRWSlicePtr[T any](v []T) *Slice[T] { return toPtr(NewRWSlice[T](v)) }

//----------------------

type ISlice[T any] interface {
	Locker[[]T]
	Append(els ...T)
	Clone() (out []T)
	DeleteIdx(i int)
	Each(clb func(T))
	GetIdx(i int) (out T)
	Insert(i int, el T)
	Len() (out int)
	Pop() (out T)
	Shift() (out T)
	Unshift(el T)
}

// Compile time checks to ensure type satisfies ISlice interface
var _ ISlice[any] = (*Slice[any])(nil)

type Slice[V any] struct {
	Locker[[]V]
}

func newBaseSlicePtr[V any](m Locker[[]V]) *Slice[V] {
	return &Slice[V]{m}
}

func (s *Slice[T]) Each(clb func(T)) {
	s.RWith(func(v []T) {
		for _, e := range v {
			clb(e)
		}
	})
}

func (s *Slice[T]) Append(els ...T) {
	s.With(func(v *[]T) { *v = append(*v, els...) })
}

// Unshift insert new element at beginning of the slice
func (s *Slice[T]) Unshift(el T) {
	s.With(func(v *[]T) { *v = append([]T{el}, *v...) })
}

// Shift (pop front)
func (s *Slice[T]) Shift() (out T) {
	s.With(func(v *[]T) { out, *v = (*v)[0], (*v)[1:] })
	return
}

func (s *Slice[T]) Pop() (out T) {
	s.With(func(v *[]T) { out, *v = (*v)[len(*v)-1], (*v)[:len(*v)-1] })
	return
}

func (s *Slice[T]) Clone() (out []T) {
	s.RWith(func(v []T) {
		out = make([]T, len(v))
		copy(out, v)
	})
	return
}

func (s *Slice[T]) Len() (out int) {
	s.RWith(func(v []T) { out = len(v) })
	return
}

func (s *Slice[T]) GetIdx(i int) (out T) {
	s.RWith(func(v []T) { out = (v)[i] })
	return
}

func (s *Slice[T]) DeleteIdx(i int) {
	s.With(func(v *[]T) { *v = (*v)[:i+copy((*v)[i:], (*v)[i+1:])] })
}

func (s *Slice[T]) Insert(i int, el T) {
	s.With(func(v *[]T) {
		var zero T
		*v = append(*v, zero)
		copy((*v)[i+1:], (*v)[i:])
		(*v)[i] = el
	})
}

//----------------------

type RWUInt64[T ~uint64] struct {
	*RWMtx[T]
}

func NewRWUInt64[T ~uint64]() RWUInt64[T] {
	return RWUInt64[T]{NewRWMtxPtr[T](0)}
}

func NewRWUInt64Ptr[T ~uint64]() *RWUInt64[T] {
	return &RWUInt64[T]{NewRWMtxPtr[T](0)}
}

func (s *RWUInt64[T]) Incr(diff T) {
	s.With(func(v *T) { *v += diff })
}

func (s *RWUInt64[T]) Decr(diff T) {
	s.With(func(v *T) { *v -= diff })
}
