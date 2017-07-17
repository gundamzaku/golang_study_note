## func ArrayOf(count int, elem Type) Type {}

从功能上来讲，ArrayOf()的注释已经讲得很清楚了。  
```
// ArrayOf returns the array type with the given count and element type.
// For example, if t represents int, ArrayOf(5, t) represents [5]int.
```
也就是说你传他一个数组的长度，和数组的类型，他就返回给你这么个数组变量。  
```go
func main()  {
	var t string
	rs:= reflect.ArrayOf(5,reflect.TypeOf(t))
	fmt.Println(rs.String())
}
result:
[5]string
```
不过功能虽然简单，实现的代码可是超长的。  
同样的，他要上锁和读取缓存。这听上去似乎和之前的MakeXX函数类似，其实这也算是一个MakeArray()的函数吧。  
```go
// Look in cache.
ckey := cacheKey{Array, typ, nil, uintptr(count)}
if array := cacheGet(ckey); array != nil {
	return array
}
很容易可以看到这个cache的机制，同样的，他会以数组类型ID+要生成的数组类型（比如string）+生成的数组大小来创建一个key。
```

确认缓存中不存在以后，他便会以Array的规则来Make一个Array出来，这段代码就不去理解了。代码一拉到底，最后仍旧是  
`return cachePut(ckey, &array.rtype)`的方法，将一个生成的数组Type返回。

不过这个返回的Type并非是一个Array类型，是不能直接拿来用的。  

```go
func main()  {
	var t string
	rs:= reflect.ArrayOf(2,reflect.TypeOf(t))
	fmt.Println(rs.String())
	rs[1] = "hello"
}
result:
invalid operation: rs[1] (type reflect.Type does not support indexing)
```

我要使用他，就必须进行转换，在reflect包中，提供了这么一个转换的方法reflect.New()
```go
func main()  {

	var t string

	rs:= reflect.ArrayOf(2,reflect.TypeOf(t))

	v := reflect.New(rs).Elem()
	fmt.Println(v.String())

	v.Index(1).Set(reflect.ValueOf("hello"))
	fmt.Println(v)
}
```

reflect.New()的方法非常简单  
```go
// New returns a Value representing a pointer to a new zero value
// for the specified type. That is, the returned Value's Type is PtrTo(typ).
func New(typ Type) Value {
	if typ == nil {
		panic("reflect: New(nil)")
	}
	ptr := unsafe_New(typ.(*rtype))
	fl := flag(Ptr)
	return Value{typ.common().ptrTo(), ptr, fl}
}
```
主要是在调用了unsafe_New()这个方法上面。这个方法的具体实现是在runtime包的malloc.go之中。  
```go
//go:linkname reflect_unsafe_New reflect.unsafe_New
func reflect_unsafe_New(typ *_type) unsafe.Pointer {
	return newobject(typ)
}
```
而newobject(typ)调用的又是下面这个方法。
```go
// Allocate an object of size bytes.
// Small objects are allocated from the per-P cache's free lists.
// Large objects (> 32 kB) are allocated straight from the heap.
func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer {}
```
看注释的意思是给对象分配一个比特大小的空间吧。小的对象从per-P缓存的空闲列中分配。大的目标（大于32KB）的则从heap中分配。  

不知道我是不是可以这么理解，我在ArrayOf（）以后，等于有了一个房产证，可是光凭房产证，我没房子可以进去，而New（）其实就是把这个房子给了我。有了这个房子，我就可以进去了。

New（）返回给我的是一个Value对象，我可以直接用ValueOf（）么？答案是否定的。

```go
func main()  {
	var t string
	rs:= reflect.ArrayOf(2,reflect.TypeOf(t))
	v:=reflect.ValueOf(rs)
	v.Index(1).Set(reflect.ValueOf("hello"))
	fmt.Println(v)
}
result:
panic: reflect: call of reflect.Value.Index on ptr Value
```
无法运作，深究一下，你会发现，两个方法返回的类型是不相同的。  
```go
func main()  {
	var t string
	rs:= reflect.ArrayOf(2,reflect.TypeOf(t))
	v:=reflect.ValueOf(rs)
	fmt.Println(v.String())
	v1 := reflect.New(rs).Elem()
	fmt.Println(v1.String())

}
result:
<*reflect.rtype Value>
<[2]string Value>
```
