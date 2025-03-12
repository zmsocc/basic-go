package main

// RemoveElement 删除切片中指定下标的元素
func RemoveElement(slice []int, idx int) []int {
	// 检查下标是否在有效范围
	if idx < 0 || idx >= len(slice) {
		return slice
	}
	// 使用切片的切片操作删除指定下标
	return append(slice[:idx], slice[idx+1:]...)
}

// Delete 改造为泛型方法
func Delete[T Number](slice []T, idx int) []T {
	if idx < 0 || idx >= len(slice) {
		return slice
	}
	return append(slice[:idx], slice[idx+1:]...)
}

type Number interface {
	~int | ~float64 | ~string
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
