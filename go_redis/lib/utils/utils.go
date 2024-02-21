package utils

// BytesEquals 判断给定的两个字节数组是否相等
func BytesEquals(a []byte, b []byte) bool {
	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	size := len(a)
	for i := 0; i < size; i++ {
		aValue := a[i]
		bValue := b[i]
		if aValue != bValue {
			return false
		}
	}
	return true
}
