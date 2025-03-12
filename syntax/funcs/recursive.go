package main

// Recursive 递归使用不当，可能会出现StackOverflow， 默认栈大小是2KB
func Recursive(n int) {
	if n == 10 {
		return
	}
	Recursive(n + 1)
}
