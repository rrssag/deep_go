package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ValueType interface {
	int | int8 | int16 | int32 | int64
}

type Values[T ValueType] []T

// go test -v homework_test.go

type CircularQueue[T ValueType] struct {
	values          []T
	elementsInQueue int
	pointerStart    int
	pointerEnd      int
}

func NewCircularQueue[T ValueType](size T) *CircularQueue[T] {
	return &CircularQueue[T]{
		values:       make(Values[T], size),
		pointerStart: 0,
		pointerEnd:   -1,
	}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.Full() {
		return false
	}

	q.incEndPointer()
	q.values[q.pointerEnd] = value
	q.elementsInQueue++

	return true
}

func (q *CircularQueue[T]) Pop() bool {
	if q.Empty() {
		return false
	}

	q.values[q.pointerStart] = 0
	q.elementsInQueue--
	q.incStartPointer()

	return true
}

func (q *CircularQueue[T]) Front() T {
	if q.Empty() {
		return -1
	}

	return q.values[q.pointerStart]
}

func (q *CircularQueue[T]) Back() T {
	if q.Empty() {
		return -1
	}

	return q.values[q.pointerEnd]
}

func (q *CircularQueue[T]) Empty() bool {
	return q.elementsInQueue == 0
}

func (q *CircularQueue[T]) Full() bool {
	return q.elementsInQueue == cap(q.values)
}

func (q *CircularQueue[T]) incEndPointer() {
	if q.pointerEnd == cap(q.values)-1 {
		q.pointerEnd = 0
		return
	}

	q.pointerEnd++
}

func (q *CircularQueue[T]) incStartPointer() {
	if q.pointerStart == cap(q.values)-1 {
		q.pointerStart = 0
		return
	}

	q.pointerStart++
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue(queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
