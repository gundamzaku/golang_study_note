## func NewAt(typ Type, p unsafe.Pointer) Value {}

比起reflect.New()来说，reflect.NewAt显得更加简单和纯粹，就是你传个Type类型进去，给你加个flag，再以Value类型返回。  

```go
func main()  {
	var t string = "world"
	var p string
	p = "hello"
	//返回一个值，该值表示指定类型的指针的值，使用p作为该指针。
	vp:= reflect.NewAt(reflect.TypeOf(t),unsafe.Pointer(&p))
	fmt.Println(vp.Elem())
	fmt.Println(p)
}
result:
hello
hello
```
这是我根据自己的理解来写的，然则一开始理解的是把指针p指到了t的地址，理论上打印p的时候，应该出来的是t的值。  
可是实际上这段代码里实现了的是，用了p的值，t的类型，拼成了一个指针变量。听起来有点绕口，实际上现在你把t改成string类型看一下。  
```
result:
4922681
hello
```
实际上p受到了t的类型的约束，会被强行转义。  
这是不是NewAt()的实际用法？我有点迷糊。于是在网上又找了一个例子，可是这段例子看下来，人就更晕了。  

```go
package main
import (
	"fmt"
	"reflect"
	"unsafe"
)

func main(){
	var a []int
	var b []int
	var value reflect.Value = reflect.ValueOf(&a)
	var value1 reflect.Value = reflect.ValueOf(&b)
	
	value = reflect.Indirect(value) //使指针指向内存地址
	value1 = reflect.Indirect(value1) //使指针指向内存地址
	
	value1 = reflect.NewAt(value.Type(), unsafe.Pointer(reflect.ValueOf(&b).Pointer())).Elem()
	b = append(b, 1)
	a = append(a, 2)
	fmt.Println(value.Pointer(), value1.Pointer())
	//>>282918952 282918960
	
	value.Set(value1)
	
	fmt.Println(value1.Kind(), value.Pointer(), value1.Pointer(), value.Interface(), value1.Interface(), a, b)
	//>>slice 282918960 282918960 [1] [1] [1] [1]
}
result:
825741050440 825741050432
slice 825741050432 825741050432 [1] [1] [1] [1]
```

我试着将这段例子改成了我之前的写法。  

```go
func main(){
	var a string
	var b string

	value_a := reflect.Indirect(reflect.ValueOf(&a)) //使指针指向内存地址
	value_b := reflect.Indirect(reflect.ValueOf(&b)) //使指针指向内存地址

	value_b = reflect.NewAt(value_a.Type(), unsafe.Pointer(reflect.ValueOf(&b).Pointer())).Elem()
	b = "hello"
	a = "world"
	fmt.Println(value_a, value_b)

	value_a.Set(value_b)

	fmt.Println(value_b.Kind(),value_a.Interface(), value_b.Interface(), a, b)
}
result:
world hello
string hello hello hello hello
```

它用到了```value_a.Set(value_b)```这个方法，可是又有什么意义呢？纯粹告诉我value_a和a是同一个地址么？所以值会一起变。  

比如我把`value_a.Set(value_b)`改成`value_b.Set(reflect.ValueOf("big clude"))`

```go
result:
world hello
string hello big clude hello big clude
```
看样子是改变了b的值，而b和value_b是指同一个地址。  
```go
val := MakeSlice(typ, 0, cap)
data := NewAt(ArrayOf(cap, typ), unsafe.Pointer(val.Pointer()))
```

最后，我找到一个官方的例子，这个例子是从test的源码中抠出来的，现在将它整理了一下，或者能更好地说明NewAt()到底是个怎么样的用法。  

```go
func main(){

	type Ptr struct{
		x string
	}

	var Tptr = reflect.TypeOf(Ptr{})
	typ:= reflect.SliceOf(Tptr)
	val := reflect.MakeSlice(typ, 10, 10)

	var listp Ptr
	listp.x = "hello"
	var p []Ptr
	p = append(p,listp)
	fmt.Println(reflect.ValueOf(p).String())

	data := reflect.NewAt(reflect.ArrayOf(10, typ), unsafe.Pointer(val.Pointer()))

	var s []int
	s = append(s,5,6)
	data.Elem().Index(1).Set(reflect.ValueOf(p))
	fmt.Println(data)
	fmt.Println(data.Elem().Index(1))
	fmt.Println(data.Elem().Index(1).Kind())
	fmt.Println(data.Elem().Index(1).Index(0))
	fmt.Println(data.Elem().Index(1).Index(0).FieldByName("x"))
}

result:
<[]main.Ptr Value>
&[[] [{hello}] [] [] [] [] [] [] [] []]
[{hello}]
slice
{hello}
hello
```
最关键的是`data := reflect.NewAt(reflect.ArrayOf(10, typ), unsafe.Pointer(val.Pointer()))`  
data产生了一个指针，这个指针指到一个数组（reflect.ArrayOf(10, typ)），而数组里面放的又是一个slice变量，slice变量里面又放着是一个ptr的结构体。  
绕了一圈。
