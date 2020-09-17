package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int                          // длина списка
	Front() *ListItem                  // первый элемент списка
	Back() *ListItem                   // последний элемент списка
	PushFront(v interface{}) *ListItem // добавить значение в начало
	PushBack(v interface{}) *ListItem  // добавить значение в конец
	Remove(i *ListItem)                // удалить элемент
	MoveToFront(i *ListItem)           // переместить элемент в начало
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	// Place your code here
	len   int
	front *ListItem
	back  *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	var head = &ListItem{
		Value: v,
		Next:  l.front,
		Prev:  nil,
	}

	return l.pushFrontPointer(head)
}

func (l *list) pushFrontPointer(head *ListItem) *ListItem {
	l.len = l.len + 1

	if l.front != nil {
		head.Next = l.front
		head.Prev = nil
		l.front.Prev = head
		l.front = head
		return head
	} else {
		l.front = head
		l.back = head
		return head
	}
}

func (l *list) PushBack(v interface{}) *ListItem {
	var back = &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.back,
	}

	l.len = l.len + 1

	if l.back != nil {
		back.Prev = l.back
		back.Next = nil
		l.back.Next = back
		l.back = back
		return back
	} else {
		l.front = back
		l.back = back
		return back
	}
}

func (l *list) Remove(i *ListItem) {
	next := i.Next
	prev := i.Prev

	l.len = l.len - 1

	if next != nil { // не последний
		next.Prev = prev
	} else {
		l.back = prev
	}

	if prev != nil { // не первый
		prev.Next = next
	} else {
		l.front = next
	}
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.pushFrontPointer(i)
}

func NewList() List {
	return &list{}
}
