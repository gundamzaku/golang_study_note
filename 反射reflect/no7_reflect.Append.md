## func Append(s Value, x ...Value) Value {}

之前在reflect的tyoeof()方法上耗去了大量的时间，这让人有些沮丧。  

让我们看看reflect包中一共有多少个方法呢？  

```go
reflect.Append()
reflect.Select()
reflect.AppendSlice()
reflect.ArrayOf()
reflect.ChanOf()
reflect.Copy()
reflect.DeepEqual()
reflect.FuncOf()
reflect.Indirect()
reflect.MakeChan()
reflect.MakeFunc()
reflect.MakeMap()
reflect.MakeSlice()
reflect.MapOf()
reflect.New()
reflect.NewAt()
reflect.PtrTo()
reflect.SliceOf()
reflect.StructOf()
reflect.Swapper()
reflect.ValueOf()
reflect.Zero()
```

哇！一共有22个方法，加上之前的typeof()是23个。没办法，既然别人能写得出来，我们也应该可以看得下来不是。  

先从上向下看吧，第一个方法  
`reflect.Append()`  
完整个的方法是结构是`func Append(s Value, x ...Value) Value {}`  

先看一下它的注释：  
```
// Append appends the values x to a slice s and returns the resulting slice.
// As in Go, each x's value must be assignable to the slice's element type.
```
把值X追加到切片s，并返回切片结果。在Go里面，每一个x的值必须被指派为切片元素类型。  
大概的意思似乎明白了，就是切片的合并。  
这样的话我写一段代码先试一下。  
```go
func main()  {
	
	slice := make([]string, 5)
	slice[0] = "Mon"
	slice[1] = "Tues"
	slice[2] = "Wed"
	slice[3] = "Thur"
	slice[4] = "Fri"

	sliceExt := make([]string, 2)
	sliceExt[0] = "Sat"
	sliceExt[1] = "sum"
	reflect.Append(slice,sliceExt)
	fmt.Println(slice)
}
```
第一次尝试，直接挂  
```
.\mian.go:33: cannot use slice (type []string) as type reflect.Value in argument to reflect.Append
.\mian.go:33: cannot use sliceExt (type []string) as type reflect.Value in argument to reflect.Append
```
经过调整，我发现两件事情。首先，传入reflect.Append()的两个值必须是Value类型。其次，是切片附加的操作，传入的第一个值是Slice类型，第二个则是其它类型（如int，string)，如果一来，我把方法再改了一下。  
```
func main()  {
	/*
	slice := make([]string, 5)
	slice[0] = "Mon"
	slice[1] = "Tues"
	slice[2] = "Wed"
	slice[3] = "Thur"
	slice[4] = "Fri"
	*/
	//为了简化，我改用下面的方式声明一个切片
	slice := []string{"Mon","Tues","Wed","Thur","Fri"}

	var newVal string = "sat"
	newSlice := reflect.Append(reflect.ValueOf(slice),reflect.ValueOf(newVal))
	fmt.Println(newSlice)
}
```
从方法中可以看到我用了一个reflect.ValueOf()的方法，将切片类型做了转换，转成了reflect.Append()要求的value类型。  
reflect.ValueOf()做了些什么？只能再到源代码中去翻。  

在上面的那个方法中，我将slice变量传入reflect.ValueOf(),而在reflect.ValueOf()中，将slice转成了原始类型interface{} 

`*一切类型都可以变成interface{}类型`

接下来，在里面依次执行了escapes(i)和unpackEface(i)方法，并直接将unpackEface(i)方法作为返回值。  

看上去很简单。先看一下escapes(i)做了些什么吧？  

这里我将main方法改造成了最简单的一个方法来进行测试：  
```go
func main()  {
	var newVal string = "hello"
	rs := reflect.ValueOf(newVal)
	fmt.Println(rs)
}
result:hello(什么都没有改变)
```

```
func escapes(x interface{}) {
	if dummy.b {
		dummy.x = x
	}
}

var dummy struct {
	b bool
	x interface{}
}
```
第一次调用了escapes(i)的方法，很明显dummy.b是false，因此，它什么都没有做，可以跳过。  
`*dummy是复制品的意思，这里应该是一个克隆方法吧（猜）` 

那再看一下和unpackEface(i)方法吧。  
```go
func unpackEface(i interface{}) Value {
	e := (*emptyInterface)(unsafe.Pointer(&i))
	// NOTE: don't read e.word until we know whether it is really a pointer or not.
	t := e.typ
	if t == nil {
		return Value{}
	}
	f := flag(t.Kind())
	if ifaceIndir(t) {
		f |= flagIndir
	}
	return Value{t, e.word, f}
}
```
哇，看上去第一行就好眼熟啊，这是我们的老方法了，得到变量的信息。它主要为了获取e.typ和e.word两个数据。  
这里还有一个警告，让我不要读e.word直到我们确认他是否是指针。（什么鬼……）  
e.typ就是变量的信息主体了。再继续读下去，到flag(t.Kind())  
t.Kind()是得到变量的类型的数字编号，这里我得到的是24,即string  
flag其实是就是uintptr的别名，这里基本就是uintptr(24)的意思  
转了一下还是24啊。  
然后又是ifaceIndir(t)，看注释上讲：  
```
ifaceIndir reports whether t is stored indirectly in an interface value.
ifaceIndir报告t是否直接存在interface变量之中的  
```
输出结果是true  
ifaceIndir()方法中执行的是  
`return t.kind&kindDirectIface == 0`  
t.kind我知道的，是24，那kindDirectIface又是什么？   
在上面的常量定义中我看到：`kindDirectIface = 1 << 5`  
`变量数据类型(1<<5)即为1左移5位，二进制即为 0010 0000`  
日日哦，读到这里我已经完全不知道是什么意思了，t.kind=24,kindDirectIface=32，两个值做&运算，得到的是24&32即结果是0。  
这TM是与运算了吧，32的二进制是1000 0000  
0010 0000  
1000 0000  
现在我们可以温习一下与运算是什么了。  
```
0 & 0=0
0 & 1=0
1 & 0=0
1 & 1=1
```
由此可见，结果就是0，最后0==0，返回的肯定是true了。  
从代码的功能上来看，这一段主要是防止越界，在value.go里面不是定义了一组常量么  
```go
const (
	Invalid Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	String
	Struct
	UnsafePointer
)
```

虽然只用到了28位，估计是在go里面认为从0-32都是go的变量类型对应的数值，因此是不能超过32的。  
既然通过了，即么代码又来到了`f |= flagIndir`这一段，又是奇怪的运算符，真是讨厌！  
flagIndir是一个uintptr类型，而且也是一个讨厌的位运算数字 1 << 7(结果是：128)  
`|=`即为：按位或后赋值  
这么一弄，`f |= flagIndir`，最后f变成了152了。  
最后，以Value结构体返回结果  
`return Value{t, e.word, f}`
```go
type Value struct {
	typ *rtype
	ptr unsafe.Pointer
	flag
}
```

总结一下就是，将一个变量转换成Value类型，其实就是在这个变量的内部给它追加了一个flag，其它的都没有任何变化，至于为什么要加个flag，现在我也不懂，虽然在源码里面密密麻麻地写了一大串注释，可是还是看不太懂，感觉是要对着注释机翻了。  
```
// flag holds metadata about the value.
flag保持了关于变量的元数据
// The lowest bits are flag bits:
最低的比特是flag比特（什么鬼）
//	- flagStickyRO: obtained via unexported not embedded field, so read-only
//	- flagEmbedRO: obtained via unexported embedded field, so read-only
//	- flagIndir: val holds a pointer to the data
//	- flagAddr: v.CanAddr is true (implies flagIndir)
//	- flagMethod: v is a method value.
// The next five bits give the Kind of the value.
接下来的五个比特位可以得到值的类型
// This repeats typ.Kind() except for method values.
这重复了typ.Kind(),除去为方法值（翻译得狗屁不通）
// The remaining 23+ bits give a method number for method values.
剩下的23+的比特为方法的值提供了方法数字
// If flag.kind() != Func, code can assume that flagMethod is unset.
如果flag.kind()不是func，代码能分派flagMethod未设置（在我的例子里，flag.kind()仍然得到的是24，即String）
// If ifaceIndir(typ), code can assume that flagIndir is set.
```
算了，翻译得实在是太差了，就此打住，不丢人现眼了。  
