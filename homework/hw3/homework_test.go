package main

import (
	"reflect"
	"runtime"
	"slices"
	"sync/atomic"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type COWBuffer struct {
	data       []byte
	refCounter *int32
}

func NewCOWBuffer(data []byte) COWBuffer {
	buff := COWBuffer{
		data:       data,
		refCounter: new(int32),
	}

	runtime.SetFinalizer(&buff, func(b *COWBuffer) {
		b.Close()
	})

	return buff
}

func (b *COWBuffer) Clone() COWBuffer {
	atomic.AddInt32(b.refCounter, 1)

	newBuff := COWBuffer{
		data:       unsafe.Slice(unsafe.SliceData(b.data), len(b.data)),
		refCounter: b.refCounter,
	}

	runtime.SetFinalizer(&newBuff, func(b *COWBuffer) {
		b.Close()
	})

	return newBuff
}

func (b *COWBuffer) Close() {
	if atomic.LoadInt32(b.refCounter) > 0 {
		atomic.AddInt32(b.refCounter, -1)
	}

	b.data = nil
}

func (b *COWBuffer) Update(index int, value byte) bool {
	if b == nil || index < 0 || index >= len(b.data) {
		return false
	}

	if atomic.LoadInt32(b.refCounter) > 1 {
		atomic.AddInt32(b.refCounter, -1)
		*b = NewCOWBuffer(slices.Clone(b.data))
	}

	b.data[index] = value

	return true
}

func (b *COWBuffer) String() string {
	return unsafe.String(unsafe.SliceData(b.data), len(b.data))
}

func TestCOWBuffer(t *testing.T) {
	data := []byte{'a', 'b', 'c', 'd'}
	buffer := NewCOWBuffer(data)
	defer buffer.Close()

	copy1 := buffer.Clone()
	copy2 := buffer.Clone()

	assert.Equal(t, unsafe.SliceData(data), unsafe.SliceData(buffer.data))
	assert.Equal(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	assert.True(t, (*byte)(unsafe.SliceData(data)) == unsafe.StringData(buffer.String()))
	assert.True(t, (*byte)(unsafe.StringData(buffer.String())) == unsafe.StringData(copy1.String()))
	assert.True(t, (*byte)(unsafe.StringData(copy1.String())) == unsafe.StringData(copy2.String()))

	assert.True(t, buffer.Update(0, 'g'))
	assert.False(t, buffer.Update(-1, 'g'))
	assert.False(t, buffer.Update(4, 'g'))

	assert.True(t, reflect.DeepEqual([]byte{'g', 'b', 'c', 'd'}, buffer.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy1.data))
	assert.True(t, reflect.DeepEqual([]byte{'a', 'b', 'c', 'd'}, copy2.data))

	assert.NotEqual(t, unsafe.SliceData(buffer.data), unsafe.SliceData(copy1.data))
	assert.Equal(t, unsafe.SliceData(copy1.data), unsafe.SliceData(copy2.data))

	copy1.Close()

	previous := copy2.data
	copy2.Update(0, 'f')
	current := copy2.data

	// 1 reference - don't need to copy buffer during update
	assert.Equal(t, unsafe.SliceData(previous), unsafe.SliceData(current))

	copy2.Close()
}
