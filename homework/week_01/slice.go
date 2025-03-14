package main

import (
	"errors"
	"fmt"
)

// RemoveElement 删除切片中指定下标的元素
func RemoveElement(slice []int, idx int) []int {
	// 检查下标是否在有效范围
	if idx < 0 || idx >= len(slice) {
		return slice
	}
	// 使用切片的切片操作删除指定下标
	return append(slice[:idx], slice[idx+1:]...)
}

var ErrIndexOutOfRange = errors.New("下标超出范围")

// Delete 改造为泛型方法
// 如果下标不是合法的，返回 ErrIndexOutOfRange
func Delete[T Number](slice []T, idx int) ([]T, error) {
	if idx < 0 || idx >= len(slice) {
		return nil, fmt.Errorf("ekit: %w, 下标超出范围，长度 %d，下标 %d", ErrIndexOutOfRange, len(slice), idx)
	}
	for i := idx; i+1 < len(slice); i++ {
		slice[i] = slice[i+1]
	}

	return Shrink[T](slice[:len(slice)-1]), nil
}

type Number interface {
	~int | ~float64 | ~string
}

// Shrink 是缩容
func Shrink[T Number](src []T) []T {
	c, l := cap(src), len(src)
	n, changed := calCapacity(c, l)
	if !changed {
		return src
	}
	s := make([]T, 0, n)
	s = append(s, src...)
	return s
}

func calCapacity(c, l int) (int, bool) {
	// 容量 <=64 缩不缩都无所谓， 因为浪费内存也浪费不了多少
	// 可以考虑调大这个阈值， 或者调小这个阈值
	if c <= 64 {
		return c, false
	}
	// 如果容量大于 2048， 但是元素不足一半，
	// 降低为 0.625， 也就是 5/8
	// 也就是比一半多一点， 和正向扩容的 1.25 倍相呼应
	if c > 2048 && (c/l >= 2) {
		factor := 0.625
		return int(float32(c) * float32(factor)), true
	}
	// 如果在 2048 以内， 并且元素不足1/4， 那么直接缩减为一半
	if c <= 2048 && (c/l >= 4) {
		return c / 2, true
	}
	// 整个实现的核心是希望在后续少触发扩容的前提下，一次性释放尽可能多的内存
	return c, false
}

//func DeleteSlice(n int, arr []int) []int {
//	for i := 0; i < len(arr); i++ {
//		if i == n {
//			for j := i; j < len(arr)-1; j++ {
//				arr[j] = arr[j+1]
//			}
//			arr[len(arr)-1] = 0
//			fmt.Printf("删除特定下标后的切片为: %v \n", arr)
//			return arr
//		}
//	}
//	return nil
//}
