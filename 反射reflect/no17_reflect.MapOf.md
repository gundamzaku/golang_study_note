### reflect.MapOf(key, elem Type) Type {}

传入key,elem，均为Type类型，返回一个Type类型。

Source:
```go
func MapOf(key, elem Type) Type {
	ktyp := key.(*rtype)
	etyp := elem.(*rtype)

	if !ismapkey(ktyp) {
		panic("reflect.MapOf: invalid key type " + ktyp.String())
	}

	// Look in cache.
	ckey := cacheKey{Map, ktyp, etyp, 0}
	if mt := cacheGet(ckey); mt != nil {
		return mt
	}

	// Look in known types.
	s := "map[" + ktyp.String() + "]" + etyp.String()
	for _, tt := range typesByString(s) {
		mt := (*mapType)(unsafe.Pointer(tt))
		if mt.key == ktyp && mt.elem == etyp {
			return cachePut(ckey, tt)
		}
	}

	// Make a map type.
	var imap interface{} = (map[unsafe.Pointer]unsafe.Pointer)(nil)
	mt := **(**mapType)(unsafe.Pointer(&imap))
	mt.str = resolveReflectName(newName(s, "", "", false))
	mt.tflag = 0
	mt.hash = fnv1(etyp.hash, 'm', byte(ktyp.hash>>24), byte(ktyp.hash>>16), byte(ktyp.hash>>8), byte(ktyp.hash))
	mt.key = ktyp
	mt.elem = etyp
	mt.bucket = bucketOf(ktyp, etyp)
	if ktyp.size > maxKeySize {
		mt.keysize = uint8(ptrSize)
		mt.indirectkey = 1
	} else {
		mt.keysize = uint8(ktyp.size)
		mt.indirectkey = 0
	}
	if etyp.size > maxValSize {
		mt.valuesize = uint8(ptrSize)
		mt.indirectvalue = 1
	} else {
		mt.valuesize = uint8(etyp.size)
		mt.indirectvalue = 0
	}
	mt.bucketsize = uint16(mt.bucket.size)
	mt.reflexivekey = isReflexive(ktyp)
	mt.needkeyupdate = needKeyUpdate(ktyp)
	mt.ptrToThis = 0

	return cachePut(ckey, &mt.rtype)
}
```
代码很长，从执行结果来看
```go
func main()  {
	newkv:=reflect.MapOf(reflect.TypeOf("b"),reflect.TypeOf("i"))
	fmt.Println(newkv)
}

result:
map[string]string
```
就只是返回数据类型。看样子，和reflect.TypeOf()是一样的，所以在`m :=reflect.MakeMap(reflect.TypeOf(val))`中，一般也会以`newkv:=reflect.MapOf(reflect.TypeOf("b"),reflect.TypeOf("i"))`来填充吧。

```go
func main()  {
	val:=reflect.MapOf(reflect.TypeOf("c"),reflect.TypeOf("baga"))
	m :=reflect.MakeMap(val)
	m.SetMapIndex(reflect.ValueOf("a"),reflect.ValueOf("h"))
	m.SetMapIndex(reflect.ValueOf("b"),reflect.ValueOf("i"))
	fmt.Println(m)
}
```
以上就是替换的代码，值得注意的是，这种写法有很强的约束代码，他约束了你产生的代码必须是Key=string,val=string  
所以下面的代码是会报错的：  
```go
func main()  {
	val:=reflect.MapOf(reflect.TypeOf("c"),reflect.TypeOf(12))
	m :=reflect.MakeMap(val)
	m.SetMapIndex(reflect.ValueOf("a"),reflect.ValueOf("h"))
	m.SetMapIndex(reflect.ValueOf("b"),reflect.ValueOf("i"))
	fmt.Println(m)
}
```
因为我约束成了key=string,val=int了。  

在reflect包中，有很多这种xxOf()的方法，其作用基本上都是一致的，接着，看一下代码的实现原理吧，因为我仿佛看到了cache的身影。  

头两行
```go
ktyp := key.(*rtype)
etyp := elem.(*rtype)
```
获得了各自的rtype。  

```go
if !ismapkey(ktyp) {
	panic("reflect.MapOf: invalid key type " + ktyp.String())
}
```
ismapkey()这个方法非常奇怪，因为在type.go中并没有实现它，而是放在了runtime包中的hashmap.go文件中。
```go
func ismapkey(t *_type) bool {
	return t.alg.hash != nil
}
```
而且和以前看到的不一样，他是验证t.alg.hash是否为空。  
在我的代码中，得到的是一个地址：0x4bd250，基本上也算是通过了。
alg.hash，算法表中的hash，具体什么定义，恐怕只有系统知道了。这一段应该只是验证类型的合法性，毕竟，我就算是把这个key传一个结构体进来，也是合法的（虽然后面报错了）

```go
type selfVal struct {
	name string
	age int
}
func main()  {
	s:= selfVal{}
	val:=reflect.MapOf(reflect.TypeOf(s),reflect.TypeOf("a"))//这里并没有任何是错误的
}
```

接下来，又开始锁进程，然后去缓存去寻找缓存了（之前有讲过）  
```go
// Look in cache.
ckey := cacheKey{Map, ktyp, etyp, 0}
if mt := cacheGet(ckey); mt != nil {
	return mt
}
```
因为这段前面已经讲过，这里看到也简单了很多，map类型有n种组合，比如
key:string,val:int  
key:int,val:string  
等等，每种组合一种ckey，如果发现已经存在，就直接返回Type，lookupCache中的m也是一个集合(map)呢。  
cache不存在的时候，那么开始生成`s := "map[" + ktyp.String() + "]" + etyp.String()`  
后面的方法和前面讲的大抵上都雷同，主要是指流程上面。细节还是有不一样的。我可以对比一下早先看过的SliceOf()  

```go
+// Make a map type. -// Make a slice type.
var imap interface{} = (map[unsafe.Pointer]unsafe.Pointer)(nil)
mt := **(**mapType)(unsafe.Pointer(&imap))
mt.str = resolveReflectName(newName(s, "", "", false))
mt.tflag = 0
mt.hash = fnv1(etyp.hash, 'm', byte(ktyp.hash>>24), byte(ktyp.hash>>16), byte(ktyp.hash>>8), byte(ktyp.hash))
mt.key = ktyp
mt.elem = etyp
mt.bucket = bucketOf(ktyp, etyp)
if ktyp.size > maxKeySize {
	mt.keysize = uint8(ptrSize)
	mt.indirectkey = 1
} else {
	mt.keysize = uint8(ktyp.size)
	mt.indirectkey = 0
}
if etyp.size > maxValSize {
	mt.valuesize = uint8(ptrSize)
	mt.indirectvalue = 1
} else {
	mt.valuesize = uint8(etyp.size)
	mt.indirectvalue = 0
}
mt.bucketsize = uint16(mt.bucket.size)
mt.reflexivekey = isReflexive(ktyp)
mt.needkeyupdate = needKeyUpdate(ktyp)
mt.ptrToThis = 0
```
SliceOf()
```go

var islice interface{} = ([]unsafe.Pointer)(nil)
prototype := *(**sliceType)(unsafe.Pointer(&islice))
slice := *prototype
slice.tflag = 0
slice.str = resolveReflectName(newName(s, "", "", false))
slice.hash = fnv1(typ.hash, '[')
slice.elem = typ
slice.ptrToThis = 0
```
