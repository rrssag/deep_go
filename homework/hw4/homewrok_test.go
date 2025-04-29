package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

// Узел бинарного дерева поиска
type BSTNode struct {
	Key   int
	Value int
	Left  *BSTNode
	Right *BSTNode
}

// Упорядоченный словарь
type OrderedMap struct {
	root *BSTNode
	size int
}

// Создать новый упорядоченный словарь
func NewOrderedMap() OrderedMap {
	return OrderedMap{nil, 0}
}

// Вставить элемент в словарь
func (m *OrderedMap) Insert(key, value int) {
	m.root = insert(m.root, key, value)
	m.size++
}

// Метод вставки узла в BST
func insert(root *BSTNode, key, value int) *BSTNode {
	if root == nil {
		return &BSTNode{key, value, nil, nil}
	}
	if key < root.Key {
		root.Left = insert(root.Left, key, value)
	} else if key > root.Key {
		root.Right = insert(root.Right, key, value)
	}

	return root
}

// Удалить элемент из словаря
func (m *OrderedMap) Erase(key int) {
	m.root = erase(m.root, key)
	// Уменьшаем размер, только если элемент был удален
	if !m.Contains(key) {
		m.size--
	}
}

// Метод удаления узла из BST
func erase(root *BSTNode, key int) *BSTNode {
	if root == nil {
		return nil
	}
	if key < root.Key {
		root.Left = erase(root.Left, key)
	} else if key > root.Key {
		root.Right = erase(root.Right, key)
	} else {
		if root.Left == nil {
			return root.Right
		} else if root.Right == nil {
			return root.Left
		}
		// Если узел имеет оба потомка, заменяем его минимальным узлом из правого поддерева
		root.Key, root.Value = findMin(root.Right)
		root.Right = deleteMin(root.Right)
	}
	return root
}

// Найти минимальный узел в правом поддереве
func findMin(node *BSTNode) (int, int) {
	for node.Left != nil {
		node = node.Left
	}
	return node.Key, node.Value
}

// Удалить минимальный узел из поддерева
func deleteMin(node *BSTNode) *BSTNode {
	if node.Left == nil {
		return node.Right
	}
	node.Left = deleteMin(node.Left)
	return node
}

// Проверить существование элемента в словаре
func (m *OrderedMap) Contains(key int) bool {
	return contains(m.root, key)
}

// Метод проверки наличия узла в BST
func contains(root *BSTNode, key int) bool {
	if root == nil {
		return false
	}
	if key < root.Key {
		return contains(root.Left, key)
	} else if key > root.Key {
		return contains(root.Right, key)
	}
	return true
}

// Получить количество элементов в словаре
func (m *OrderedMap) Size() int {
	return m.size
}

// Применить функцию к каждому элементу словаря от меньшего к большему
func (m *OrderedMap) ForEach(action func(int, int)) {
	forEach(m.root, action)
}

// Метод обхода дерева в порядке возрастания ключей (ин-порядок)
func forEach(root *BSTNode, action func(int, int)) {
	if root == nil {
		return
	}
	forEach(root.Left, action)
	action(root.Key, root.Value)
	forEach(root.Right, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
