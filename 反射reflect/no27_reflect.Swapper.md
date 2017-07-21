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

我试着用指针的形式去试了一下程序  
```go
func main(){

	var v1 string = "hello"
	var v2 string = "world"

	var s []*string
	s = append(s,&v1)
	s = append(s,&v2)
	fmt.Println(s)

	fc := reflect.Swapper(s)
	fc(0,1)
	fmt.Print(s)
}
result:
[0xc042008240 0xc042008250]
[0xc042008250 0xc042008240]
```

其实并没有什么区别，算了，靠瞎猜不如仔细看一下代码。  

在源代码中，`hasPtr := typ.kind&kindNoPointers == 0`  
也就是说，如果typ.kind&kindNoPointers不等于0，hasPtr就是false了。而kindNoPointers值为1 << 7，即128。typ.kind的值范围肯定在0-32逃不离了。那这样我就可以写一段代码来看一下哪几个类型符合条件。
首先，我先到Type.go文件里面，找到`type Type interface {}`的结构体，在里面加一个小小的方法`CommonTo() uint8`  
并且在下面实现它：
```go
func (t *rtype) CommonTo() uint8 { return t.kind }
```
这一段主要是方便我拿到rtype里面的kind属性，原本kind是在程序内部的私有变量，我是得不到的。  
```go
type element struct {
	a int
	b int16
	c string
	d bool
}
var e element
for k := 0; k < reflect.TypeOf(e).NumField(); k++ {
	//fmt.Println(reflect.TypeOf(e).Field(k).Type.CommonTo())
	fmt.Print(reflect.TypeOf(e).Field(k).Name)
	fmt.Print("--")
	fmt.Println(reflect.TypeOf(e).Field(k).Type.CommonTo()&128)
}
result:
a--128
b--128
c--0
d--128
```
结果一看，原来除了string，int和bool都是属于false的情况，这样后面就容易理解了。
```go
switch size {
case 8:
	is := *(*[]int64)(v.ptr)
	return func(i, j int) { is[i], is[j] = is[j], is[i] }
case 4:
	is := *(*[]int32)(v.ptr)
	return func(i, j int) { is[i], is[j] = is[j], is[i] }
case 2:
	is := *(*[]int16)(v.ptr)
	return func(i, j int) { is[i], is[j] = is[j], is[i] }
case 1:
	is := *(*[]int8)(v.ptr)
	return func(i, j int) { is[i], is[j] = is[j], is[i] }
}
```
这应该是针对整形的一些优化。bool其实也应该算是0和1的整形吧。  
源码在这里还有一层判断，就是size不在8，4，2，1的数字之中的时候，会跳到后面一层代码去进行转化。  
不过到了这一层，应该是一些极端情况了，在此也不再作敷述。
