## func DeepEqual(x, y interface{}) bool {}

DeepEqual是我所学习的reflect包中的最后一个对外函数了，看本体代码很短
```go
func DeepEqual(x, y interface{}) bool {
	if x == nil || y == nil {
		return x == y
	}
	v1 := ValueOf(x)
	v2 := ValueOf(y)
	if v1.Type() != v2.Type() {
		return false
	}
	return deepValueEqual(v1, v2, make(map[visit]bool), 0)
}
```
然而它的注释却比代码长得多得多。
