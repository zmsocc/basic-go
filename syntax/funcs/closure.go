package main

import "fmt"

func Closure(name string) func() string {
	return func() string {
		return "hello" + name
	}
}

func Closure1() func() string {
	name := "大名"
	age := 18
	return func() string {
		return fmt.Sprintf("Hello, %s, %d", name, age)
	}
}

func Closure2() func() int {
	age := 0
	return func() int {
		age++
		return age
	}
}
