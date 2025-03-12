package main

import "fmt"

func ForLoop() {
	for i := 0; i < 10; i++ {
		println(i)
	}

	for i := 0; i < 10; {
		println(i)
		i++
	}
}

func Loop2() {
	i := 0
	for i < 10 {
		println(i)
		i++
	}
}

func ForArr() {
	arr := [3]int{1, 2, 3}
	for index, val := range arr {
		println("下标", index, "值", val)
	}
}

func ForMap() {
	m := map[string]int{
		"key1": 100,
		"key2": 102,
	}

	for i, j := range m {
		println(i, j)
	}
}

func LoopBug() {
	users := []User{
		{
			name: "大明",
		},
		{
			name: "小明",
		},
	}
	m := make(map[string]*User)
	for _, U := range users {
		m[U.name] = &U // 不要对迭代参数取地址
	}

	fmt.Printf("%v", m)
}

type User struct {
	name string
}

func LoopBreak() {
	i := 0
	for true {
		if i > 10 {
			break
		}
		i++
		println(i)
	}
}

func LoopContinue() {
	i := 0
	for i < 10 {
		if i%2 == 1 {
			continue
		}
		println(i)
		i++
	}
}
