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
然而它的注释却比代码长得多得多。我一扫而过，只记住了一个单词“deeply equal”，深度相等。所以我也很清楚地明白这个方法的定义。  
粗看代码，首先要两个传入的参数不为空，接着两个参数的Type()必须相等，这样还不够，符合条件了。还要走deepValueEqual（）方法进行验证。  
这个deepValueEqual(v1, v2, make(map[visit]bool), 0)，不仅仅是把v1,v2两个参数传了进去，还代入了个map变量。  

deepValueEqual()就不同了，这才是整个方法的核心部分。  

首先，他还会再验证一次Type()
```go
if !v1.IsValid() || !v2.IsValid() {
	return v1.IsValid() == v2.IsValid()
}
if v1.Type() != v2.Type() {
	return false
}
```
此还，还根据v1.Kind()，针对不同的类型（Array、Slice、Interface、Ptr、Struct、Map、Func），还进行了递归的操作。  
至于代码部分，大至了解一下就行了。  

接下来还是写点例子看一下吧。
