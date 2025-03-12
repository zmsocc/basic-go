package main

import "fmt"

func main() {
	l := []string{"zms", "z", "zzz", "hgg", "dog", "cat", "7.01"}
	fmt.Printf("删除指定下标后的切片为: %v \n", Delete(l, 1))
}
