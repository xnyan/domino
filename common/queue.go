package common

import (
	"strconv"
)

type node struct {
	Val  interface{}
	next *node
}

func NewNode(val interface{}) *node {
	return &node{Val: val}
}

type Queue struct {
	head *node
	tail *node
	len  int

	itr *node // For iteration
}

func NewQueue() *Queue {
	q := &Queue{
		head: NewNode(nil),
		len:  0,
	}
	q.tail = q.head
	return q
}

func (q *Queue) Peek() (interface{}, bool) {
	if q.head.next == nil {
		return nil, false
	}

	return q.head.next.Val, true
}

func (q *Queue) Push(val interface{}) {
	n := NewNode(val)
	q.tail.next = n
	q.tail = n
	q.len++
}

// Returns (val, exists)
func (q *Queue) Pop() (interface{}, bool) {
	if q.head.next == nil {
		return nil, false
	}
	q.len--
	n := q.head.next
	q.head.next = n.next
	if q.head.next == nil {
		q.tail = q.head
	}

	return n.Val, true
}

func (q *Queue) Size() int {
	return q.len
}

//// Iterator interface
func (q *Queue) Iterator() Iterator {
	q.itr = q.head
	return q
}

func (q *Queue) HasNext() bool {
	return q.itr.next != nil
}

func (q *Queue) Next() (interface{}, bool) {
	if q.itr.next == nil {
		return nil, false
	}
	q.itr = q.itr.next
	return q.itr.Val, true
}

func (q *Queue) String() string {
	ret := ""
	cur := q.head.next
	for cur != nil {
		switch cur.Val.(type) {
		case int:
			ret += strconv.Itoa(cur.Val.(int)) + ";"
		case string:
			ret += cur.Val.(string) + ";"
		}

		cur = cur.next
	}
	return ret
}
