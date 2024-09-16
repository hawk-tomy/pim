package list

import (
	"container/list"
	"fmt"
)

// Element is an element of a linked list. see list.Element
type Element[T any] struct {
	element *list.Element
	// The value stored with this element. same as list.Element.Value
	Value T
}

func newElement[T any](e *list.Element) *Element[T] {
	return &Element[T]{element: e, Value: e.Value.(T)}
}

// Next returns the next list element or nil. see list.Element.Next
func (e *Element[T]) Next() *Element[T] {
	p := e.element.Next()
	if p == nil {
		return nil
	}

	return newElement[T](p)
}

// Prev returns the previous list element or nil. see list.Element.Prev
func (e *Element[T]) Prev() *Element[T] {
	p := e.element.Prev()
	if p == nil {
		return nil
	}

	return newElement[T](p)
}

// List represents a doubly linked list.
// The zero value for List is an empty list ready to use.
// see list.List
type List[T any] struct {
	list *list.List
}

// Init initializes or clears list l. see list.List.Init
func (l *List[T]) Init() *List[T] {
	l.list = l.list.Init()
	return l
}

// New returns an initialized list. see list.List.New
func New[T any]() *List[T] {
	return &List[T]{list.New()}
}

// Len returns the number of elements of list l.
// The complexity is O(1).
// see list.List.Len
func (l *List[T]) Len() int {
	return l.list.Len()
}

// Front returns the first element of list l or nil if the list is empty. see list.List.Front
func (l *List[T]) Front() *Element[T] {
	e := l.list.Front()
	if e == nil {
		return nil
	}

	return newElement[T](e)
}

// Back returns the last element of list l or nil if the list is empty. see list.List.Back
func (l *List[T]) Back() *Element[T] {
	e := l.list.Back()
	if e == nil {
		return nil
	}

	return newElement[T](e)
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
// The element must not be nil.
// see list.List.Remove
func (l *List[T]) Remove(e *Element[T]) T {
	return l.list.Remove(e.element).(T)
}

// PushFront inserts a new element e with value v at the front of list l and returns e. see list.List.PushFront
func (l *List[T]) PushFront(v T) *Element[T] {
	e := l.list.PushFront(v)
	return newElement[T](e)
}

// PushBack inserts a new element e with value v at the back of list l and returns e. see list.List.PushBack
func (l *List[T]) PushBack(v T) *Element[T] {
	e := l.list.PushBack(v)
	return newElement[T](e)
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
// see list.List.InsertBefore
func (l *List[T]) InsertBefore(v T, e *Element[T]) *Element[T] {
	ret := l.list.InsertBefore(v, e.element)
	if ret == nil {
		return nil
	}

	return newElement[T](ret)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
// see list.List.InsertAfter
func (l *List[T]) InsertAfter(v T, e *Element[T]) *Element[T] {
	ret := l.list.InsertAfter(v, e.element)
	if ret == nil {
		return nil
	}

	return newElement[T](ret)
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
// see list.List.MoveToFront
func (l *List[T]) MoveToFront(e *Element[T]) {
	l.list.MoveToFront(e.element)
}

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
// see list.List.MoveToBack
func (l *List[T]) MoveToBack(e *Element[T]) {
	l.list.MoveToBack(e.element)
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
// see list.List.MoveBefore
func (l *List[T]) MoveBefore(e, mark *Element[T]) {
	l.list.MoveBefore(e.element, mark.element)
}

// MoveAfter moves element e to its new position after mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
// see list.List.MoveAfter
func (l *List[T]) MoveAfter(e, mark *Element[T]) {
	l.list.MoveAfter(e.element, mark.element)
}

// PushBackList inserts a copy of another list at the back of list l.
// The lists l and other may be the same. They must not be nil.
// see list.List.PushBackList
func (l *List[T]) PushBackList(other *List[T]) {
	l.list.PushBackList(other.list)
}

// PushFrontList inserts a copy of another list at the front of list l.
// The lists l and other may be the same. They must not be nil.
// see list.List.PushFrontList
func (l *List[T]) PushFrontList(other *List[T]) {
	l.list.PushFrontList(other.list)
}

func (l *List[T]) String() string {
	var ts []string
	for e := l.list.Front(); e != nil; e = e.Next() {
		if es, ok := e.Value.(fmt.Stringer); ok {
			ts = append(ts, es.String())
		} else {
			ts = append(ts, fmt.Sprintf("%v", e.Value))
		}

	}
	return fmt.Sprintf("%v", ts)
}
