## type Pointer *ArbitraryType

在go的reflect里面，我看到很多地方都用了unsafe.point方法，那么，这个unsafe.point到底是什么用处呢？  
首先，我查看了一下unsafe包的源码，可惜里面除了声明，并没有实现的方法。  
在声明中有注释：  
```go
// Pointer represents a pointer to an arbitrary type.  There are four special operations
// available for type Pointer that are not available for other types.
//	1) A pointer value of any type can be converted to a Pointer.
//	2) A Pointer can be converted to a pointer value of any type.
//	3) A uintptr can be converted to a Pointer.
//	4) A Pointer can be converted to a uintptr.
// Pointer therefore allows a program to defeat the type system and read and write
// arbitrary memory. It should be used with extreme care.
type Pointer *ArbitraryType
```
用我拙劣的英文翻译了一下   
指针以任意一种类型描述成一个新的指针，一共有四种特别的操作  
1.一个任意类型的指针值能被换化成一个指针  
2.一个指针能被转换成一个任意类型的指针值  
3.一个uintptr能被转换成一个指针  
4.一个指针能被转换成一个uintptr  

`uintptr：专门用于存储内存地址的类型`  

听上去很复杂，这里又多了一个新的概念，叫uintptr。不过这些我们目前都可以忽略，记住一记即可：unsafe.point是实现指针类型的转换。  
下面就写点代码进行测试  

```go
func main()  {

	var num int = 800
	var i *int
	i = &num
	fmt.Println(i)

	var s *string
	s=(*string)(unsafe.Pointer(i))
	fmt.Println(s)
}
result:
0xc820072220
0xc820072220
```

两个相同的地址，似乎什么都没有发生？现在我们把地址转换成数值看一下吧。  

```go
fmt.Println(*i)
fmt.Println(*s)
result:
800
```

怎么回事？只出来一个值？`*s`呢？  
```go
var rs string = *s
fmt.Println(int(rs))
```
我换一个形式看一下吧  
`./main.go:23: cannot convert rs (type string) to type int`
直接报错，从错误描述上来看，rs确实已经被转换成了string类型，只是为什么他消失了？  

```go
var s *string
s=(*string)(unsafe.Pointer(i))
var temp string = "hello"
s = &temp
fmt.Println(*s)

result:
hello
```
从这段来看，大概是转换的过程中，值丢失了。  
`不过我有一点不明白，既然s和i都指向了同一个地址，为什么s和i的赋值并不冲突。`

既然用unsafe.point不能实现不同类型的转换，那么同类型的呢？我再试一下。  
```go
func main()  {
	var num int = 800
	var i *int
	i = &num
	fmt.Println(*i)

	var s *int64
	s=(*int64)(unsafe.Pointer(i))
	fmt.Println(*s)
}

result:
800
800
```
从int->int64成功  

```go
func main()  {

	var num string = "800"
	var i *string
	i = &num
	fmt.Println(*i)

	var s *int
	s=(*int)(unsafe.Pointer(i))
	fmt.Println(*s)
}

result:
800
1066152
```
从string->int成功但是这值是什么鬼？  
由此看来，unsafe.point确实是unsafe了。  

再回到之前reflect代码上面来看这个point  
```go
type emptyInterface struct {
	typ  int
	word unsafe.Pointer
}

func main()  {

	var num int = 800
	var i *int
	i = &num
	fmt.Println(*i)

	var s *emptyInterface
	s=(*emptyInterface)(unsafe.Pointer(i))
	fmt.Println(s.typ)
	fmt.Println(s.word)
}

result:
800
800
0xc82003ded8
```
这里是将int指针转成一个结构体。并非太难理解  

最后，我将笔记（1)中的代码回顾一下，发现  

```go
var num int = 800
var e interface{}
e = num

eface := *(*emptyInterface)(unsafe.Pointer(&e))
fmt.Println(*eface.typ.name)
```
有一个问题，就是之前的结构体中的代码里，int类型必须转成interface{}类型，才可以进行指针转换。  
否则无法与`emptyInterface`中的`rtype`结构对应上，而报错。  
这大概是因为`rtyp`的结构体是interface{}的内置对应属性，于我们程序开发来说，是无法更改的。  

unsafe中除了point以外，还有两个方法，目前在此次学习之后，仍然是一知半解。  
