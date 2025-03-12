package main

const (
	StatusA = iota
	StatusB
	StatusC
	StatusD
)

const (
	DsyA = iota*12 + 13
	DsyB
	DsyC
)

const (
	A = iota << 1
	B // 0001 左移一位变为 0010
	C
	D
)

func main() {
	const a = 123 // 常量不能更改
}
