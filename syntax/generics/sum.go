package main

import "io"

// Sum 是一个泛型方法
func Sum[T Number](vals ...T) T {
	var res T
	for _, val := range vals {
		res = res + val
	}
	return res
}

// Number 是一个泛型约束
type Number interface {
	~int | int64 | float64 // ~int 是指int及其衍生类型
}

type Integer int

func ReleaseResource[R io.Closer](r R) {
	r.Close()
}
