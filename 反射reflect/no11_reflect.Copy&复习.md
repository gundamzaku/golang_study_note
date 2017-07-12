## func Copy(dst, src Value) int {}

整理一下现在的思绪，有必要再回顾一下之前的一些知识点。  

先看上一篇中无法执行的代码。

```go
var a = [5]int {1, 2, 3, 4, 5}
var b = [5]int {6, 7, 8, 9, 0}
c := reflect.Copy(reflect.ValueOf(a),reflect.ValueOf(b))
```

问题就出在reflect.ValueOf()上面，在重新看这个方法的时候，我们先追溯一下源头，Value到底是一种什么样的类型。  
Value是一个对象：  
```go
type Value struct {
  typ *rtype
  ptr unsafe.Pointer
  flag
}
func (v Value) pointer() unsafe.Pointer {}
func (v Value) Addr() Value {}
等等
```
在源代码中，有丰富的注释解决了这个结构体的作用。  
```
// Value is the reflection interface to a Go value.
//value是Go的一个反射型接口
// Not all methods apply to all kinds of values. Restrictions,
不是所有方法支付值的所有类型
……
// Its IsValid method returns false, its Kind method returns Invalid,
无效的方法返回false，它的种类方法返回无效。  
Using == on two Values does not compare the underlying values
两个Value类型不能用来进行比较
// To compare two Values, compare the results of the Interface method.
要想比较，需要使用接口方法
```

比较这一部分是真的吗？我先试一下看看。  
```go
var a int32 = 16
var b int32 = 16
fmt.Println(reflect.ValueOf(a)==reflect.ValueOf(b))  
```
果然返回false，两者已经不能进行比较了。
我们可以分别对比一下这两个a，b经过valueOf()已经分别变成了两个不同的对象。
而这个Value对象一共有58个方法，比较夸张。  
我打印几个常用的方法看一下  
```go
var a int32 = 16
va := reflect.ValueOf(a)
fmt.Println(va.String())
fmt.Println(va.Type())
fmt.Println(va.CanAddr())
result:
<int32 Value>
int32
false
```
这些方法就是一个变量常用的信息吧，也就是说我把一个变量转化成了Value类型以后，就可以调用查看这个变量的相关属性信息（甚至是操作）  
这里的va.CanAddr()就是前面我碰到的报错的原因。我的数组，转成了Value，va.CanAddr()返回的是False，导致无法进行Copy操作。原因不明  

既然如此，我就单独拎这一段代码的功能来看看，它做了什么事情。  

```go
// CanAddr reports whether the value's address can be obtained with Addr.
CanAddr报告了这个值的地址是否能以Addr的方式获得？
// Such values are called addressable. A value is addressable if it is
诸如此值被称为addressable，一个值可以addressable
// an element of a slice, an element of an addressable array,
假如他是一个slice元素，一个addressable(可访问)的数组,
// a field of an addressable struct, or the result of dereferencing a pointer.
一个可访问的结构的领域，或者一个非关联化的指针
// If CanAddr returns false, calling Addr will panic.
如果返回false，你调用Addr将会报错。
func (v Value) CanAddr() bool {
	return v.flag&flagAddr != 0
}
```
好了，看完了，表示没看懂，为什么我定义的数组CanAddr是False呢？  
我现在转成数组  
```go
var a [2]int32
a[0] = 1
a[1] = 2
va := reflect.ValueOf(a)
fmt.Println(va.String())
fmt.Println(va.Type())
fmt.Println(va.CanAddr())

result:
<[2]int32 Value>
[2]int32
false
```

确实CanAddr()是False。说明我这个不是一个an addressable array，那什么才是an addressable array，真是好抓狂。  

正当我陷入深深地绝望的时候，万能的Google赐给了我一段下面的代码（这可是我搜索了近4个小时才找到的一点曙光，千万不能错过）

```go
a := 100
va := reflect.ValueOf(a)
vp := reflect.ValueOf(&a).Elem()
fmt.Println(va.CanAddr(), va.CanSet())
fmt.Println(vp.CanAddr(), vp.CanSet())

resut:
false false
true true
```
竟然有CanAddr()为true了，`vp := reflect.ValueOf(&a).Elem()`这段代码给了我灵感，我试着将自己的代码也改成这样的形式。  

```go
var b [2]int
b[0] = 2
va := reflect.ValueOf(&b).Elem()
fmt.Println(va.String())
fmt.Println(va.Type())
fmt.Println(va.CanAddr())
```	
成功了！va.CanAddr()的结果为true，也不管具体的原理，我迫不急待地拿上一篇报错的代码进行进一步的调试。  
```go
var a = [5]int {1, 2, 3, 4, 5}
var b = [5]int {6, 7, 8, 9, 0}
c := reflect.Copy(reflect.ValueOf(&a).Elem(),reflect.ValueOf(b))
fmt.Println(a)
fmt.Println(b)
fmt.Println(c)
```
成功地把b复制到了a身上！感觉google，一个困扰了一天的难道终于迎刃而解。  
现在就要看一下到底为什么要这样呢？  
&a是一个地址，程序把地址转换成了一个Value类型？还是点进源码里看一看究竟。  
```go
在这个方法中，我们把&a作为一个interface{}传入。
func unpackEface(i interface{}) Value {
	//传入之前的a实际地址是0xc042048270
	//i:(0x491e00,0xc042048270)这是转成interface{}以后的形式
	//下面这段代码可以将传入的i还原回int数组
	//var b [5]int;
	//b = *(*[5]int)(unsafe.Pointer(i.(*[5]int)))
	
	e := (*emptyInterface)(unsafe.Pointer(&i))//这段可以理解成指向了i的地址的指针
	t := e.typ//赋于了变量的属性,e.typ.String()显示为*[5]int,e指向了i,e.typ.Kind()显示为22
	if t == nil {
		return Value{}
	}
	f := flag(t.Kind())//现在它变成了Ptr的类型，数字为22，不是数组了

	if ifaceIndir(t) {
		f |= flagIndir
	}
	return Value{t, e.word, f}//这里返回的是一个指针
}
```
`注意`在我的调试过程中，出现了一个奇怪的问题，因为我一直将build -a 这个参数打开的，然而我在执行上面的代表的时候，意外的发现reflect.ValueOf(&a)会执行两次。我不知道是什么问题，写了一段程序测试一下。  

```go
func main()  {
	var a int = 10;
	fmt.Println(a)
	geti(a)
}

func geti(i int){
	fmt.Println(i)
}
正常
```

```go
func main()  {
	var a int = 10;
	fmt.Println(a)
	geti(&a)
}

func geti(i *int){
	fmt.Println(*i)
}

正常
```

```go
func main()  {
	var a int = 10;
	fmt.Println(a)
	geti(&a)
}

func geti(i *int){
	fmt.Println(i)
}

result:
(0x491b60,0xc0420381d0) //这是我在valueOf()方法中的打印
10
0xc0420381d0
```

很费解，惟一的区别就是正常的那次我是打印指针指向的变量了，不正常的这次我是打印指针地址，可是不知道为什么会去触发ValueOf()方法。在此存疑。

reflect.ValueOf(&a)得到的是一个指针->a的地址，这个指针怎么会有Elem()？
再回到Copy()方法里面。
dk := dst.kind() //输出为17，是指针。奇怪，怎么又变成了指针了？

我又回到Elem()方法上面，看了一眼它的注释，似乎都明白了……
```
// Elem returns the value that the interface v contains
Elem返回interface v包含的值
// or that the pointer v points to.
或者指针v所指向的地址（重点）
// It panics if v's Kind is not Interface or Ptr.
在v的类型不是interface或ptr的时候报错
// It returns the zero Value if v is nil.
当v是nil的时候返回空值
```
似乎是全明白了……  
再看一下里面的代码部分吧  
```go
case Ptr: //如果是指针
	ptr := v.ptr//指针所指的地址给指针，ptr这时候就指向了指向数组的地址
	if v.flag&flagIndir != 0 {//校验
		ptr = *(*unsafe.Pointer)(ptr)
	}
	// The returned value's address is v's value.
	if ptr == nil {//为空时
		return Value{}
	}
	tt := (*ptrType)(unsafe.Pointer(v.typ))//组合成ptrType
	typ := tt.elem //取到数组
	fl := v.flag&flagRO | flagIndir | flagAddr
	fl |= flag(typ.Kind())
	return Value{typ, ptr, fl}
}
```

差不多就这样吧，Copy()方法基本上就全部了解了。
