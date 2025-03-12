package main

type LinkedList struct {
	head *node
	tail *node

	// 这个就是包外可访问
	Len int
}

func (l *LinkedList) Add(idx int, val any) error {
	//TODO implement me
	panic("implement me")
}

func (l *LinkedList) Append(val any) {
	//TODO implement me
	panic("implement me")
}

func (l *LinkedList) Delete(idx int) (any, error) {
	//TODO implement me
	panic("implement me")
}

//func (l LinkedList) Add(idx int, val any) {
//
//}
//
//func (l *LinkedList) AddV1(idx int, val any) {
//
//}

type node struct {
	prev *node
	next *node
}
