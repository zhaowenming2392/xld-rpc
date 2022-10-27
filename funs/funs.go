/*
常用函数
*/
package funs

/*
三元表达式
a, b := 2, 3
max := If3(a > b, a, b).(int)
println(max)
*/
func If3(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
