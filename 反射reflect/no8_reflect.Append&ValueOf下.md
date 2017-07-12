## func Append(s Value, x ...Value) Value {}

继续前面的话题，在了解了`reflect.ValueOf(newVal)`这个方法是干什么用的以后，我们再回过头来看`reflect.Append`
reflect.Append(s Value, x ...Value)带了两个参数，其实不只两个，这个`x ...Value`表示你可以输入无数个Value类型的x，用","分割。  
append()有一个限制，你传入的s必须是一个slice类型，而后面的x可以是string，也可以是int之类的变量类型。  
然后看一下append()代码的内部  
```go
func Append(s Value, x ...Value) Value {
	s.mustBe(Slice)
	s, i0, i1 := grow(s, len(x))
	for i, j := i0, 0; i < i1; i, j = i+1, j+1 {
		s.Index(i).Set(x[j])
	}
	return s
}
```
一下子感觉好简单，首先这个`s.mustBe(Slice)`，一看就知道是什么意思了。
它是以变量s的kind()属性和常量的Slice数字对比，如果都是23（代表Slice)类型，那么就通过验证。  
`s, i0, i1 := grow(s, len(x))`接下来的这一段，gorw()是一个比较长的函数。
```go
/*
 * 传入s切片，x的数量（就是一共传进来几个x)
 * 传出新的切片，两个整数
 */
func grow(s Value, extra int) (Value, int, int) {
	i0 := s.Len() //s的切片数量
	i1 := i0 + extra //要产生的切片的新数量
	if i1 < i0 { //新切片的数量不可能比老切片要少，所以要报错
		panic("reflect.Append: slice overflow")
	}
	m := s.Cap()  //cap()函数返回的是数组切片分配的空间大小。
  //如果空间大小足够，那么直接返回
	if i1 <= m {
		return s.Slice(0, i1), i0, i1
	}
  //空间不足的情况下，如果完全没有空间，m的空间就是extra的数量
	if m == 0 {
		m = extra
	} else {
		//m现在比新的切片大小要小的话
    for m < i1 {
      //如果i0，即老s的切片数量在1024以内，m+m？
			if i0 < 1024 {
				m += m
			} else {
        //超过1024，m+m/4?
				m += m / 4
			}
		}
	}
  //最后产生一个新的切片返回
	t := MakeSlice(s.Type(), i1, m)
	Copy(t, s)
	return t, i0, i1
}
```

要记住的是，返回的是s, i0, i1三个变量，s代表新的切片，io代表老切片大小，i1代表新切片的大小

我的代码中，i0是5，i1是7。

整个方法中最关键的是变量m的用处，一开始的时候，在`m := s.Cap()`时，m的值为5，此时i1的值为7  
m比i1小，因此进入了下面的代码块  
```go
if i0 < 1024 {
	m += m
} else {
	m += m / 4
}
```
i0现在的值为5，肯定小于1024，因此m=m+m变成了10,其实这一段非常费解，为什么在i0小于1024的情况下m是做翻倍，而在大于1024的情况下，是翻25%呢？一个想法，可能是空间优化。没关系，继续读下去。  
m现在是10，走`t := MakeSlice(s.Type(), i1, m)`的代码，传入的变量分别是：s.Type(),7,10。
s.Type()其实也是一段非常复杂的方法，不过这里不用太去纠结这个方法的用处，因为在MakeSlice方法中，主要是想用到s.Type().kind() 和s.Type().common()，kind()值是23，表示是slice类型。这主要是做数据校验用的。抛开这一层面的话，MakeSlice的方法就剩下两段了：
```go
s := sliceHeader{unsafe_NewArray(typ.Elem().(*rtype), cap), len, cap}
return Value{typ.common(), unsafe.Pointer(&s), flagIndir | flag(Slice)}
```
typ.common()是指到rtpye地址的指针，rtype就是变量的信息；flagIndir | flag(Slice)前面也有提到过，是一个128和23的或运算，最后变成151。

那s是求出什么值呢？在传入的三个参数中，len和cap都很好理解，len是新生成的切片长度（7），cap是系统现在给你的长度(10)  
`unsafe_NewArray(typ.Elem().(*rtype), cap)`就有点费解了。特别是`typ.Elem().(*rtype)`这一段，
看针对typ.Elem()的注释  
```
// Elem returns a type's element type.
Elem返回类型的元素类型。
// It panics if the type's Kind is not Array, Chan, Map, Ptr, or Slice.
如果不是Array,Chan,Map,Ptr或Slice的话，就会出错。
```
这样一来就好理解了typ.Elem()返回的就是切片里面的数据的数据类型，现在我`println(typ.Elem().Kind())`看一下，返回是24，表示字符串类型。  
typ.Elem()是一个Type的interface，`(*rtype)`这种引用方式我就不理解了，在interface里面的方法里只有一个common()的对象是`*rtype`类型的，那会不会是它呢？不妨我直接替换在.common()试一下看看。  
`s := sliceHeader{unsafe_NewArray(typ.Elem().common(), cap), len, cap}`
成功运行了，难道是真的？可是我还是难以理解这种方式的写法，还是用一个小例子看一下吧。 

```go
package main

import "fmt"

type Type interface {
	Elem() Type
	common() *rtype
}

type rtype struct {
	size int
	name string
}
func (t *rtype) Elem() Type {
	var temp = new(rtype)
	temp.size = 100
	return temp
}
func (t *rtype)common() *rtype {
	t.size = 80 
	fmt.Println("hello")
	return t
}

func main() {

	var typ Type = new(rtype)
	rs:=typ.Elem()
	fmt.Println(rs.(*rtype))
}

result:&{100 }
```

可见，代码并没有走common()，我想错了。那么……在`rs.(*rtype)`中`(*rtype)`就是表示是这个对象它本身吧？

根据这个思路，我再写一段代码作为测试  

```go
package main

import "fmt"

type small interface {
	Test() int
	Set()
}

type club struct {
	val1 int
	val2 int
	val3 string
}

func (t *club) Set(){
	t.val1 = 1
	t.val2 = 2
	t.val3 = "hello"
}

func (t *club) Test() int{
	return 2
}

func main() {
	var s small = new(club)
	s.Set()
	fmt.Println(s.(*club))
}
result:&{1 2 hello}
```
现在应该很清晰了，他就是把`s.(*club)`里定义的变量全部返回来。  
在源代码里面`typ.Elem().(*rtype)`里，typ.Elem()本来就是返回一个Type类型。  
在Elem()里面  
```go
func (t *rtype) Elem() Type {
	switch t.Kind() {
	case Array:
		tt := (*arrayType)(unsafe.Pointer(t))
		return toType(tt.elem)
	case Chan:
		tt := (*chanType)(unsafe.Pointer(t))
		return toType(tt.elem)
	case Map:
		tt := (*mapType)(unsafe.Pointer(t))
		return toType(tt.elem)
	case Ptr:
		tt := (*ptrType)(unsafe.Pointer(t))
		return toType(tt.elem)
	case Slice:
		tt := (*sliceType)(unsafe.Pointer(t))
		return toType(tt.elem)
	}
	panic("reflect: Elem of invalid type")
}
```
我们可以直接定位到`tt := (*sliceType)(unsafe.Pointer(t))`这一段，t就是指向该切片变量的信息的地址  
```go
// sliceType represents a slice type.
type sliceType struct {
	rtype `reflect:"slice"`
	elem  *rtype // slice element type
}
```
sliceType存放的就是切片的变量信息，总结一下，就是指向了我前面生成的切片的变量的定义地址，得到了它的有关的信息。
最后，把`elem  *rtype`完整返回。

听起来绕口，其实多看几次，很容易就可以掌握了，最关键的一点是要搞清楚`rtpye`的作用。
