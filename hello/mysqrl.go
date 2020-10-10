// Package hw1 随便写的一个库
package hw1

// Sqrt 自己写的牛顿法开根号
func Sqrt(x float64) float64 {
	z := 1.0
	for diff := 100.0; diff > 0.0001; {
		oldZ := z
		z -= (z*z - x) / (2 * x)
		if oldZ < z {
			diff = z - oldZ
		} else {
			diff = oldZ - z
		}
	}
	return z
}
