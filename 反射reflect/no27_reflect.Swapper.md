## func Swapper(slice interface{}) func(i, j int) {}

从方法的名称上来看，应该是起到一个交换的作用。传入的值是slice，也就是说，必须是切片类型，否则就会报错。  
```go
if v.Kind() != Slice {
	panic(&ValueError{Method: "Swapper", Kind: v.Kind()})
}
```

而返回的，是一个func的类型，带两个传入值（int i 和 int j）

这么一看，我还是没搞明白到底是要交换什么东西？为什么传进去是切片传出来就变成了函数呢？还是先看一下注释吧。  

```go
// Swapper returns a function that swaps the elements in the provided slice.
Swapper返回一个函数，用来互换提供的切片的元素。
//
// Swapper panics if the provided interface is not a slice.
如果提供的不是切片就报错。
```

看说明是互换切片里面的元素？不多说，写个例子看一下。  
```go
func main(){
	var s []string
	s = append(s,"hello")
	s = append(s,"world")
	fmt.Println(s)

	fc := reflect.Swapper(s)
	fc(0,1)
	fmt.Print(s)
}
result:
[hello world]
[world hello]
```
这样一看，就一目了然了。这个方法就是用来给Slice交换位置的，fc(0,1)中的0和1，就是slice中的index，把第0个元素和第1个元素位置给交换一下。  
原来如此，相当简洁明了的功能。  

swapper这个方法是单独放在swapper.go的文件中的，代码不长，也不太清楚为什么会这么做。  
值得注意的是，在源代码中针对slice的长度似乎还有一点区别。  
`hasPtr := typ.kind&kindNoPointers == 0`  
这又是什么鬼东西，看描述是确认是不是指针。那么我再来试一下吧。用我上面的代码测试了一下，结果是hasPtr是true。  
那……什么时候不为true呢？
