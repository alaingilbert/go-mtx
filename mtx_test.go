package mtx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRWMtxMap_HasKey(t *testing.T) {
	m := NewMap[string, int]()
	assert.False(t, m.HasKey("a"))
	m.SetKey("a", 1)
	assert.True(t, m.HasKey("a"))
	m.DeleteKey("a")
	assert.False(t, m.HasKey("a"))
}

func TestRWMtxMap_TakeKey(t *testing.T) {
	m := NewMap[string, int]()
	m.SetKey("a", 1)
	m.SetKey("b", 2)
	m.SetKey("c", 3)
	assert.Equal(t, 3, m.Len())
	assert.True(t, m.HasKey("b"))
	val, ok := m.TakeKey("b")
	assert.True(t, ok)
	assert.Equal(t, 2, val)
	assert.False(t, m.HasKey("b"))
	_, ok = m.TakeKey("b")
	assert.False(t, ok)
	assert.Equal(t, 2, m.Len())
}
