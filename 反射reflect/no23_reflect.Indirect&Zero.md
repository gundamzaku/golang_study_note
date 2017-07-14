不知不觉已经写了21篇了，自己也有点学晕了都，到目前为止，我学了几个方法了？  
~~reflect.Append()~~  
~~reflect.AppendSlice()~~  
~~reflect.ChanOf()~~  
~~reflect.Copy()~~  
~~reflect.FuncOf()~~  
~~reflect.MakeChan()~~  
~~reflect.MakeFunc()~~  
~~reflect.MakeMap()~~  
~~reflect.MakeSlice()~~  
~~reflect.MapOf()~~  
~~reflect.SliceOf()~~  
~~reflect.ValueOf()~~  
~~reflect.Select()~~  
reflect.ArrayOf()  
reflect.DeepEqual()  
reflect.Indirect()  
reflect.New()  
reflect.NewAt()  
reflect.PtrTo()  
reflect.StructOf()  
reflect.Swapper()  
reflect.Zero()  

这么一看，竟然还有一大半啊。哭，还是加快速度吧，先看一下reflect.Indirect()，这个在之前的学习中有碰到过。  

reflect.Indirect()这个方法非常简单，看注释也是很简单的，就是如果是指针，就返回地址，如果不是指针，就传进来什么返回就是什么。  
```go
// Indirect returns the value that v points to.
// If v is a nil pointer, Indirect returns a zero Value.
// If v is not a pointer, Indirect returns v.
func Indirect(v Value) Value {
	if v.Kind() != Ptr {
		return v
	}
	return v.Elem()
}
```

写一段不是指针的例子，
```go
func main()  {
	var s string
	s = "hello world"
	ns:=reflect.Indirect(reflect.ValueOf(s))
	fmt.Println(ns)
}
result:
hello world
```
返封不动的返回来了。

再找一个指针的例子，
```go
func main()  {
	var p *string
	var s string
	s = "hello world"
	p = &s
	fmt.Println(*p)

	ns:=reflect.Indirect(reflect.ValueOf(p))
	fmt.Println(ns)
}
result:
hello world
hello world
```
看上去返回的结果没有任何区别，但是实际上它执行的是`return v.Elem()`

好了，这个方法结束，再找一个简单的方法`func Zero(typ Type) Value {}`
```go
func Zero(typ Type) Value {
	if typ == nil {
		panic("reflect: Zero(nil)")
	}
	t := typ.common()
	fl := flag(t.Kind())
	if ifaceIndir(t) {
		return Value{t, unsafe_New(typ.(*rtype)), fl | flagIndir}
	}
	return Value{t, nil, fl}
}
```
从代码来看，也是非常简单，传入一个Type类型，转换了一下flag，又传回去了。那到底它是干什么的？  
```go
func main()  {
	var s string
	s = "hello world"
	zs:=reflect.Zero(reflect.TypeOf(s))
	fmt.Println(zs)
}
```
仿佛什么都没有发生……嗯，确实什么都没有发生，连fmt.println()都没有打印出东西。这是什么情况？  

我再打印一下
```
fmt.Println(zs.String() == "")
```
返回的是true，我的值没了！  

这是怎么回事呢，问题估计是出在  
```go
if ifaceIndir(t) {
	return Value{t, unsafe_New(typ.(*rtype)), fl | flagIndir}
}
```
这一段代码上面。  

```go
// ifaceIndir reports whether t is stored indirectly in an interface value.
func ifaceIndir(t *rtype) bool {
	return t.kind&kindDirectIface == 0
}
```
ifaceIndir是干嘛的？报道t是否被直接以interface值进行存储的。我的string必然是interface{}的子变量了。  
t.kind=24
kindDirectIface=32
两者与，结果是0，就是控制传入的参数是否在Go规定的32个类型（数字形式）里面。
在的话就直接触发`return Value{t, unsafe_New(typ.(*rtype)), fl | flagIndir}`从这段代码可以了解到一点。  
typ其实存放的是变量的类型的属性定义。  
`unsafe_New(typ.(*rtype))`应该是指针，指向了具体的变量赋的值上面。
flag是变量的标志位了。
所以这里等于分配新的地址，老的没了。  

如果不在Go定义的变量类型里面怎么办？
```go
return Value{t, nil, fl}
```
不在的话就强制把地址设空了。  
总得来讲，reflect.Zero() 就是把数据给清空。
