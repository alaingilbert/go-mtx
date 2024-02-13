package mtx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

func TestMtx_Debug(t *testing.T) {
	debug = true
	defer func() { debug = false }()
	m := NewMap[string, int](nil)
	m.Len()
	m2 := NewRWMap[string, int](nil)
	m2.Len()
	assert.Equal(t, 1, 1)
}

func TestMtx_LockUnlock(t *testing.T) {
	m := NewMtx("old")
	m.Lock()
	val := m.Val()
	*val = "new"
	m.Unlock()
	assert.Equal(t, "new", m.Get())
}

func TestMtx_With(t *testing.T) {
	m := NewMtx("old")
	m.With(func(v *string) {
		*v = "new"
	})
	assert.Equal(t, "new", m.Get())
}

func TestMtx_RWith(t *testing.T) {
	m := NewMtx("old")
	m.RWith(func(v string) {
		assert.Equal(t, "old", v)
	})
}

func TestMtx_Set(t *testing.T) {
	m := NewMtx("old")
	assert.Equal(t, "old", m.Get())
	m.Set("new")
	assert.Equal(t, "new", m.Get())
}

func TestMtx_Replace(t *testing.T) {
	m := NewMtx("old")
	old := m.Replace("new")
	assert.Equal(t, "old", old)
	assert.Equal(t, "new", m.Get())
}

func TestMtx_Val(t *testing.T) {
	someString := "old"
	orig := &someString
	m := NewMtx(orig)
	val := m.Val()
	**val = "new"
	assert.Equal(t, "new", someString)
	assert.Equal(t, "new", **val)
	assert.Equal(t, "new", *orig)
}

func TestMtxPtr_Replace(t *testing.T) {
	m := NewMtxPtr("old")
	old := m.Replace("new")
	assert.Equal(t, "old", old)
	assert.Equal(t, "new", m.Get())
}

func TestRWMtx_RLockRUnlock(t *testing.T) {
	m := NewRWMtx("old")
	m.RLock()
	val := m.Val()
	assert.Equal(t, "old", *val)
	m.RUnlock()
}

func TestRWMtx_Replace(t *testing.T) {
	m := NewRWMtx("old")
	old := m.Replace("new")
	assert.Equal(t, "old", old)
	assert.Equal(t, "new", m.Get())
}

func TestRWMtxPtr_Replace(t *testing.T) {
	m := NewRWMtxPtr("old")
	old := m.Replace("new")
	assert.Equal(t, "old", old)
	assert.Equal(t, "new", m.Get())
}

func TestRWMtx_Val(t *testing.T) {
	someString := "old"
	orig := &someString
	m := NewRWMtx(orig)
	val := m.Val()
	**val = "new"
	assert.Equal(t, "new", **val)
	assert.Equal(t, "new", *orig)
}

func TestMap_GetKey(t *testing.T) {
	m := NewMap[string, int](nil)
	_, ok := m.GetKey("a")
	assert.False(t, ok)
	m.SetKey("a", 1)
	el, ok := m.GetKey("a")
	assert.True(t, ok)
	assert.Equal(t, 1, el)
}

func TestMap_GetKeyValue(t *testing.T) {
	m := NewMap[string, int](nil)
	_, _, ok := m.GetKeyValue("a")
	assert.False(t, ok)
	m.SetKey("a", 1)
	key, value, ok := m.GetKeyValue("a")
	assert.True(t, ok)
	assert.Equal(t, "a", key)
	assert.Equal(t, 1, value)
}

func TestMap_HasKey(t *testing.T) {
	m := NewMap[string, int](nil)
	assert.False(t, m.ContainsKey("a"))
	m.SetKey("a", 1)
	assert.True(t, m.ContainsKey("a"))
	m.Delete("a")
	assert.False(t, m.ContainsKey("a"))
}

func TestMap_Take(t *testing.T) {
	m := NewMap[string, int](nil)
	m.SetKey("a", 1)
	m.SetKey("b", 2)
	m.SetKey("c", 3)
	assert.Equal(t, 3, m.Len())
	assert.True(t, m.ContainsKey("b"))
	val, ok := m.Take("b")
	assert.True(t, ok)
	assert.Equal(t, 2, val)
	assert.False(t, m.ContainsKey("b"))
	_, ok = m.Take("b")
	assert.False(t, ok)
	assert.Equal(t, 2, m.Len())
}

func TestMap_Delete(t *testing.T) {
	m := NewMap[string, int](nil)
	assert.Equal(t, 0, m.Len())
	m.Delete("a")
	m.SetKey("a", 1)
	assert.Equal(t, 1, m.Len())
	m.Delete("a")
	assert.Equal(t, 0, m.Len())
}

func TestMap_Values(t *testing.T) {
	m := NewMap[string, int](nil)
	assert.Equal(t, []int{}, m.Values())
	m.Set(map[string]int{"a": 1, "b": 2, "c": 3})
	values := m.Values()
	slices.Sort(values)
	assert.Equal(t, []int{1, 2, 3}, values)
}

func TestMap_Keys(t *testing.T) {
	m := NewMapPtr[string, int](nil)
	assert.Equal(t, []string{}, m.Keys())
	m.Set(map[string]int{"a": 1, "b": 2, "c": 3})
	keys := m.Keys()
	slices.Sort(keys)
	assert.Equal(t, []string{"a", "b", "c"}, keys)
}

func TestMap_Each(t *testing.T) {
	m := NewMap[string, int](nil)
	m.Set(map[string]int{"a": 1, "b": 2, "c": 3})
	arr := make([]string, 0)
	m.Each(func(k string, v int) {
		arr = append(arr, fmt.Sprintf("%s_%d", k, v))
	})
	slices.Sort(arr)
	assert.Equal(t, []string{"a_1", "b_2", "c_3"}, arr)
}

func TestMap_Clone(t *testing.T) {
	m := NewMap[string, int](nil)
	m.Set(map[string]int{"a": 1, "b": 2, "c": 3})
	clonedMap := m.Clone()
	assert.Equal(t, 1, clonedMap["a"])
}

func TestRWMap_GetKey(t *testing.T) {
	m := NewRWMap[string, int](nil)
	_, ok := m.GetKey("a")
	assert.False(t, ok)
	m.SetKey("a", 1)
	el, ok := m.GetKey("a")
	assert.True(t, ok)
	assert.Equal(t, 1, el)
}

func TestRWMap_HasKey(t *testing.T) {
	m := NewRWMap[string, int](nil)
	assert.False(t, m.ContainsKey("a"))
	m.SetKey("a", 1)
	assert.True(t, m.ContainsKey("a"))
	m.Delete("a")
	assert.False(t, m.ContainsKey("a"))
}

func TestRWMap_Take(t *testing.T) {
	m := NewRWMap[string, int](nil)
	m.SetKey("a", 1)
	m.SetKey("b", 2)
	m.SetKey("c", 3)
	assert.Equal(t, 3, m.Len())
	assert.True(t, m.ContainsKey("b"))
	val, ok := m.Take("b")
	assert.True(t, ok)
	assert.Equal(t, 2, val)
	assert.False(t, m.ContainsKey("b"))
	_, ok = m.Take("b")
	assert.False(t, ok)
	assert.Equal(t, 2, m.Len())
}

func TestRWMap_Delete(t *testing.T) {
	m := NewRWMap[string, int](nil)
	assert.Equal(t, 0, m.Len())
	m.Delete("a")
	m.SetKey("a", 1)
	assert.Equal(t, 1, m.Len())
	m.Delete("a")
	assert.Equal(t, 0, m.Len())
}

func TestRWMap_Values(t *testing.T) {
	m := NewRWMap[string, int](nil)
	assert.Equal(t, []int{}, m.Values())
	m.Set(map[string]int{"a": 1, "b": 2, "c": 3})
	values := m.Values()
	slices.Sort(values)
	assert.Equal(t, []int{1, 2, 3}, values)
}

func TestRWMap_Keys(t *testing.T) {
	m := NewRWMapPtr[string, int](nil)
	assert.Equal(t, []string{}, m.Keys())
	m.Set(map[string]int{"a": 1, "b": 2, "c": 3})
	keys := m.Keys()
	slices.Sort(keys)
	assert.Equal(t, []string{"a", "b", "c"}, keys)
}

func TestRWMap_Each(t *testing.T) {
	m := NewRWMap[string, int](nil)
	m.Set(map[string]int{"a": 1, "b": 2, "c": 3})
	arr := make([]string, 0)
	m.Each(func(k string, v int) {
		arr = append(arr, fmt.Sprintf("%s_%d", k, v))
	})
	slices.Sort(arr)
	assert.Equal(t, []string{"a_1", "b_2", "c_3"}, arr)
}

func TestRWMap_InitialValue(t *testing.T) {
	m := NewRWMap(map[string]int{"a": 1, "b": 2, "c": 3})
	assert.Equal(t, 3, m.Len())
	val, _ := m.GetKey("b")
	assert.Equal(t, 2, val)
}

func TestRWMap_Clone(t *testing.T) {
	m := NewRWMap[string, int](nil)
	m.Set(map[string]int{"a": 1, "b": 2, "c": 3})
	clonedMap := m.Clone()
	assert.Equal(t, 1, clonedMap["a"])
}

func TestSlice(t *testing.T) {
	m := NewSlicePtr[int](nil)
	assert.Equal(t, 0, m.Len())
	m.Append(1, 2, 3)
	assert.Equal(t, 3, m.Len())
	assert.Equal(t, []int{1, 2, 3}, m.Get())
	val2 := m.Shift()
	assert.Equal(t, 1, val2)
	m.Unshift(4)
	assert.Equal(t, []int{4, 2, 3}, m.Get())
	val2 = m.Pop()
	assert.Equal(t, []int{4, 2}, m.Get())
	m.DeleteIdx(1)
	assert.Equal(t, []int{4}, m.Get())
	m.Append(5, 6, 7)
	assert.Equal(t, []int{4, 5, 6, 7}, m.Get())
	assert.Equal(t, 6, m.GetIdx(2))
	m.Insert(2, 8)
	assert.Equal(t, []int{4, 5, 8, 6, 7}, m.Get())
}

func TestRWSlice(t *testing.T) {
	m := NewRWSlice[int](nil)
	assert.Equal(t, 0, m.Len())
	m.Append(1, 2, 3)
	assert.Equal(t, 3, m.Len())
	assert.Equal(t, []int{1, 2, 3}, m.Get())
	val2 := m.Shift()
	assert.Equal(t, 1, val2)
	m.Unshift(4)
	assert.Equal(t, []int{4, 2, 3}, m.Get())
	val2 = m.Pop()
	assert.Equal(t, []int{4, 2}, m.Get())
	m.DeleteIdx(1)
	assert.Equal(t, []int{4}, m.Get())
	m.Append(5, 6, 7)
	assert.Equal(t, []int{4, 5, 6, 7}, m.Get())
	assert.Equal(t, 6, m.GetIdx(2))
	m.Insert(2, 8)
	assert.Equal(t, []int{4, 5, 8, 6, 7}, m.Get())
}

func TestSlice_InitialValue(t *testing.T) {
	m := NewSlice([]int{1, 2, 3})
	assert.Equal(t, []int{1, 2, 3}, m.Get())
}

func TestRWSlice_Clone(t *testing.T) {
	m := NewRWSlice[int](nil)
	m.Set([]int{1, 2, 3})
	clonedSlice := m.Clone()
	assert.Equal(t, []int{1, 2, 3}, clonedSlice)
}

func TestRWSlice_Each(t *testing.T) {
	m := NewRWSlicePtr[int](nil)
	m.Append(1, 2, 3)
	arr := make([]string, 0)
	m.Each(func(el int) {
		arr = append(arr, fmt.Sprintf("E%d", el))
	})
	assert.Equal(t, []string{"E1", "E2", "E3"}, arr)
}

func TestRWSlice_Filter(t *testing.T) {
	m := NewRWSlicePtr([]int{1, 2, 3, 4, 5, 6})
	out := m.Filter(func(el int) bool { return el%2 == 0 })
	assert.Equal(t, 3, len(out))
	assert.Equal(t, []int{2, 4, 6}, out)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, m.Get())
}

func TestRWUInt64(t *testing.T) {
	m := NewRWUInt64[uint64](0)
	assert.Equal(t, uint64(0), m.Get())
	m.Incr(10)
	assert.Equal(t, uint64(10), m.Get())
	m.Decr(5)
	assert.Equal(t, uint64(5), m.Get())

	mp := NewRWUInt64Ptr[uint64](0)
	assert.Equal(t, uint64(0), mp.Get())
	mp.Incr(10)
	assert.Equal(t, uint64(10), mp.Get())
	mp.Decr(5)
	assert.Equal(t, uint64(5), mp.Get())
}
