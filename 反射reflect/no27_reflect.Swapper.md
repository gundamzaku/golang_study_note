## func Swapper(slice interface{}) func(i, j int) {}

从方法的名称上来看，应该是起到一个交换的作用。传入的值是slice，也就是说，必须是切片类型，否则就会报错。  
```go
if v.Kind() != Slice {
	panic(&ValueError{Method: "Swapper", Kind: v.Kind()})
}
```

而返回的，是一个func的类型，带两个传入值（int i 和 int j）
