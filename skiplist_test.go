package skiplist

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeList(n int) *Skiplist {
	list := New()
	for i := 1; i <= n; i++ {
		key := fmt.Sprintf("key_%d", i)
		member := fmt.Sprintf("member_%d", i)
		list.Set(key, member)
	}
	return list
}

func TestSkipList_Get(t *testing.T) {
	list := New()
	val := []byte("test_val")

	list.Set("ec", val)
	list.Set("dc", 123)
	list.Set("ac", val)

	assert.Equal(t, val, list.Get("ec").Value())
	assert.Equal(t, val, list.Get("ac").value)
	assert.Equal(t, 123, list.Get("dc").value)
}

func TestSkipList_Skl(t *testing.T) {
	list := makeList(10000)

	// check list for items
	for i := 1; i < 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		member := fmt.Sprintf("member_%d", i)
		assert.Equal(t, member, list.Get(key).Value())
	}

	// delete list for items
	for i := 500; i < 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		list.Delete(key)
	}

	// check list for items
	for i := 1; i < 500; i++ {
		key := fmt.Sprintf("key_%d", i)
		member := fmt.Sprintf("member_%d", i)
		assert.Equal(t, member, list.Get(key).Value())
	}
}

func TestSkipList_Remove(t *testing.T) {
	list := New()
	val := []byte("test_val")

	list.Set("ec", val)
	list.Set("dc", 123)
	list.Set("ac", val)

	list.Delete("dc")
	list.Delete("ec")
	list.Delete("ac")

	assert.Nil(t, list.Get("ec"))
	assert.Nil(t, list.Get("ac"))
	assert.Nil(t, list.Get("dc"))

}

func TestSkipList_Keys(t *testing.T) {
	list := makeList(100)
	assert.Equal(t, 100, len(list.Keys()))
}

func TestSkipList_Update(t *testing.T) {
	list := New()
	key := "foo"
	val := []byte("test_val")

	assert.Equal(t, val, list.Set(key, val).Value())
	assert.Equal(t, 123, list.Set(key, 123).Value())
	assert.Equal(t, "foo", list.Set(key, "foo").Value())

	assert.Equal(t, "foo", list.Get(key).Value())
	assert.Equal(t, 1, len(list.Keys()))
}

func TestSkipList_Unique(t *testing.T) {
	list := makeList(10000)
	assert.Equal(t, 10000, len(list.Keys()))
	key := "key_500"
	for i := 1; i < 10; i++ {
		member := fmt.Sprintf("member_%d", i)
		assert.Equal(t, member, list.Set(key, member).Value())
	}
	assert.Equal(t, 10000, len(list.Keys()))
	assert.Equal(t, "member_9", list.Get(key).Value())
}

func TestSkipList_Foo(t *testing.T) {
	list := makeList(10000)
	assert.Equal(t, 10000, len(list.Keys()))
	key := "key_500"
	for i := 1; i < 10; i++ {
		member := fmt.Sprintf("member_%d", i)
		list.Set(key, member)
	}
	assert.Equal(t, 10000, len(list.Keys()))
	assert.Equal(t, "member_9", list.Get(key).Value())
}
