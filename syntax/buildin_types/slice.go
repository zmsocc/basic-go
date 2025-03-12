package main

import "fmt"

func Slice() {
	s1 := []int{1, 2, 3}
	fmt.Printf("s1: %v, len: %d, cap: %d\n", s1, len(s1), cap(s1))

	s2 := make([]int, 3, 4) // 直接初始化了三个元素，容量为4的切片
	fmt.Printf("s2: %v, len: %d, cap: %d\n", s2, len(s2), cap(s2))

	s2 = append(s2, 10) // 追加一个元素，没有扩容
	fmt.Printf("s2: %v, len=%d, cap=%d \n", s2, len(s2), cap(s2))

	s2 = append(s2, 100) // 再追加一个元素，扩容了
	fmt.Printf("s2: %v, len: %d, cap: %d \n", s2, len(s2), cap(s2))

	s3 := make([]int, 4)
	fmt.Printf("s3: %v, len=%d, cap=%d \n", s3, len(s3), cap(s3))

	fmt.Printf("s3[2]: %d", s3[2])
}

func SubSlice() {
	s1 := []int{2, 4, 6, 8, 10}
	s2 := s1[1:3] // cap 值为从1开始往后数到末尾，包含1
	fmt.Printf("s1: %v\ns2: %v, len=%d, cap=%d\n", s1, s2, len(s2), cap(s2))
}

func ShareSlice() {
	s1 := []int{2, 4, 6, 8, 10}
	s2 := s1[2:]
	fmt.Printf("s2: %v, len=%d, cap=%d\n", s2, len(s2), cap(s2))

	s2[0] = 100
	fmt.Printf("s2: %v, len: %d, cap: %d\n", s2, len(s2), cap(s2))
	fmt.Printf("s1: %v, len: %d, cap: %d\n", s1, len(s1), cap(s1))

	s2 = append(s2, 199)
	fmt.Printf("s2: %v, len: %d, cap: %d\n", s2, len(s2), cap(s2))
	fmt.Printf("s1: %v, len: %d, cap: %d\n", s1, len(s1), cap(s1))

	s2[1] = 19999
	fmt.Printf("s2: %v, len: %d, cap: %d\n", s2, len(s2), cap(s2))
	fmt.Printf("s1: %v, len: %d, cap: %d\n", s1, len(s1), cap(s1))
}
