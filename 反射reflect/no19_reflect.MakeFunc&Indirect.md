## func MakeFunc(typ Type, fn func(args []Value) (results []Value)) Value {}

继续拿我上一节的例子进行改造。  

```go
func main()  {

	var s []reflect.Type
	var e []reflect.Type
	var vr []string
	s=append(s,reflect.TypeOf("h"),reflect.TypeOf(vr))
	e=append(e,reflect.TypeOf("h"))
	rs:= reflect.FuncOf(s,e,true)

	f:= func(in []reflect.Value)[]reflect.Value{
		var rs []reflect.Value
		rs = append(rs,reflect.ValueOf("hello"))
		return rs
	}
	newFunc := reflect.MakeFunc(rs,f)
	fmt.Println(newFunc)
}
result:
0x47afd0
```
reflect.MakeFunc()要传的参数非常简单，一个是你用FuncOf()生成出来的func变量的定义，一个就是传一个函数的运算本体。  

可是返回给我的是一串指针地址，并没有多少用处。怎么样才能让他跑起来？  

其实这里有一个小技巧，那就是reflect包下面的all_test.go文件，里面有非常多的test代码可以参考。比如makeFunc()。  
我发现在当中还需要额外的做一步操作：  
`reflect.ValueOf(&rs).Elem().Set(newFunc)`  
照搬，再执行  
直接报错  
```go
panic: reflect.Set: value of type func(string, ...string) string is not assignable to type reflect.Type
```

再经过无数次的尝试之后，最终，我还是完成了函数的调用。  

```go
var intSwap func(string, ...string) (string)
var valueOf reflect.Value = reflect.Indirect( reflect.ValueOf(&intSwap))
var v reflect.Value = reflect.MakeFunc(rs, f)
valueOf.Set(v)
fmt.Println(intSwap("a","b"))
```
要用这么一段代码，定义一个func的变量，名为intSwap。用reflect.Indirect()直接得到它的地址valueOf。生成函数v。
把这个v直接塞给valueOf，完成关联，就可以直接用了。  

```go
func Indirect(v Value) Value {
	if v.Kind() != Ptr {
		return v
	}
	return v.Elem()
}
```
reflect.Indirect()也是reflect包中公开的一个方法，它做了些什么呢？  
首先，传进来的如果不是指针，它就直接返回。  
如果是，它就返回v.Elem() ,这个Elem()返回的是intSwap所指向的内容，现在是空的。  
valueOf.Set(v)就是把v分配给valueOf了。  
两者关系对应上以后，这个新MAKE的方式便可以使用。  

可见还是相当复杂的。

主要的操作点还是在这个Elem()上面，思绪重新整理，我抽一段test例子中的代码再分析一下。  

```go
func main()  {
  fn := func(i int) int { return i }
	incr := func(in []reflect.Value) []reflect.Value {
		return []reflect.Value{reflect.ValueOf(int(in[0].Int() + 1))}
	}
	fv := reflect.MakeFunc(reflect.TypeOf(fn), incr)
	reflect.ValueOf(&fn).Elem().Set(fv)

	rs:=fn(2)
	fmt.Println(rs)
}
result:
3
```
这段代码相对来说已经非常干净了。他的作用也很简单，就是把传入的值+1并返回。  

第一步、先声明一个fn，这个fn定义了函数是个怎么样的类型（或者我称之为接口？）它描述了传入1个Int参数，并返回一个Int参数。  
第二步、声明了一个具体的实现方法incr，你可以当他是一个实现了接口的方法。  
第三步、进行创建，将fn和incr一共传入，表示对应。并产生一个fv的变量。  
第四步、或者我可以称之为绑定吧，把这个创建的fv塞给了fn的Elem。这个函数就算是创建完成了。  

fn表示了一个地址：0x489a60，&fn 0xc042004028,表示是指向了fn地址的地址。  
现在我把这个&fn用reflect.ValueOf(&fn)包一下，就可以得到他的Elem了。  
根据定义，这里的Elem其实返回的就是fn的地址，因为&fn是指向fn的地址的地址，reflect.ValueOf(&fn).Elem()直接返回指向&fn的指针，即fn的地址。  
`Elem returns the value that the interface v contains or that the pointer v points to.`
我们确认一下  
```go
fn := func(i int) int { return i }
fmt.Println(reflect.ValueOf(fn))
incr := func(in []reflect.Value) []reflect.Value {
	return []reflect.Value{reflect.ValueOf(int(in[0].Int() + 1))}
}
fv := reflect.MakeFunc(reflect.TypeOf(fn), incr)
//reflect.ValueOf(&fn).Elem().Set(fv)
fmt.Println(reflect.ValueOf(&fn).Elem())

result:
0x489910
0x489910
```
确实是同一个地址，这里有个小问题，我把`reflect.ValueOf(&fn).Elem().Set(fv)`的注释符去掉以后，地址会发生变化。
```
result:
0x489c30
0x46bf80
```
看来，问题还是出在这个Set上面。好在这个Set方法代码量不大，所以看一下。  
```go
func (v Value) Set(x Value) {
	println(v.Pointer())
	v.mustBeAssignable()
	x.mustBeExported() // do not let unexported x leak
	var target unsafe.Pointer
	if v.kind() == Interface {
		target = v.ptr
	}
	x = x.assignTo("reflect.Set", v.typ, target)
	if x.flag&flagIndir != 0 {
		typedmemmove(v.typ, v.ptr, x.ptr)
	} else {
		*(*unsafe.Pointer)(v.ptr) = x.ptr
	}
}
```
它继承了Type{}接口，我在代码里面`println(v.Pointer())`，打印指针的数字形式。然后在外面的代码里面打印`fmt.Println(reflect.ValueOf(&fn).Elem().Pointer())`，最后得到的结果都是4760240，所以地址是一样的。

v.kind()是19，代表func，不是interface，跳过。如此一来target就是空了，后面的`x.assignTo("reflect.Set", v.typ, target)`似乎没什么意义了。

x.flag&flagIndir的操作最后得出的结果是0。

直接执行`*(*unsafe.Pointer)(v.ptr) = x.ptr`  
果然是在这里发生了变化，fn的ptr，现在指到了x.ptr（也就是fv）上面。所以地址也变了。  

你可以当fn是个空的，fv是个满的，现在就是把满的灌入到空的里面，所以这个函数便可以运作了。

看到现在，最重要的MakeFunc()还没看啊……唉，算了。下次再看吧。。
