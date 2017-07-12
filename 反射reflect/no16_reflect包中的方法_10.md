接下来，我们再看一个make方法：`func MakeMap(typ Type) Value {}`吧，看上去还是比较简单的。从名字上看，必然是生成一个map的集合了。  
先写一个例子。  
```go
func main()  {
	var val string
	m:=reflect.MakeMap(reflect.TypeOf(val))
	fmt.Println(m)
}
result:
panic: reflect.MakeMap of non-map type
```
报错了，这是怎么回事？  
错误出在MakeMap代码内的验证上面：  
```go
if typ.Kind() != Map {
	panic("reflect.MakeMap of non-map type")
}
```
传入的不是map类型？我不是要生成一个map么，怎么一定要先定义成map类型了。先不管他了，既然人家的源代码有这种限制，只能照做，我重新改一下代码。  
```go
func main()  {
	var val map[string]string
	m:=reflect.MakeMap(reflect.TypeOf(val))
	fmt.Println(m)
}

result:
map[]
```
感觉把`m:=reflect.MakeMap(reflect.TypeOf(val))`这一段去掉也没关系啊。。  

先试一下原生的办法：
```go
func main()  {
	var val map[string]string
	fmt.Println(val)
	val["a"] = "h"
	val["b"] = "e"
	val["c"] = "l"
	val["d"] = "l"
	val["e"] = "o"
	fmt.Println(val)
}
result:
panic: assignment to entry in nil map
```
报错，nil的map不能被定义，必须在内存里划分一块地址给它。  
在代码中插入一行：
```go
val = make(map[string]string)

result:
map[d:l e:o a:h b:e c:l]
```
正常了，那么，我把`val = make(map[string]string)`用reflect来替代看一下，可不可以成功呢？  
答案是<b>不可以</b>。

```go
func main()  {
	var val map[string]string
	m :=reflect.MakeMap(reflect.TypeOf(val))
	m.SetMapIndex(reflect.ValueOf("a"),reflect.ValueOf("h"))
	m.SetMapIndex(reflect.ValueOf("b"),reflect.ValueOf("i"))
	fmt.Println(m)
	/* 不能再这样赋值了
	m["a"] = "h"
	m["b"] = "e"
	m["c"] = "l"
	m["d"] = "l"
	m["e"] = "o"
	*/
}
```
其实是需要按他特定的格式做一下小小的转换才对。这样一来，就实现了与原生同样的效果，可是相对来说，也更麻烦。  
接着，我们看一下原生的Make()做了什么事情吧。  

```go
func make(Type, size IntegerType) Type
```
注释很长，大意上是这个make是通用的，slice, map, or chan都可以用。第一个传入的参数需为Type类型，何为Type类型？其实你只要传入`map[string]string`这样的定义即可，如此一来，上面那句`var val map[string]string`都可以不要，直接`val:=make(map[string]string)`就可以了。

注释中有一句话：
```
The size may be omitted, in which case a small starting size is allocated.
大小可被省略，在这其中只是产生一个被分配的起始大小。
我理解就是划分一个地址给程序，里面没有实质的东西吧。
```

好了，原代码就看到这里，官方没有提供具体的实现代码给我们，只是一个接口声明。所以我们还是掉转枪头，去看reflect包中的makeMap()  
MakeMap(typ Type)就非常简单的两行来实现  
```
m := makemap(typ.(*rtype))
return Value{typ.common(), m, flag(Map)}
```	

首先内部的一个makemap私有方法，我传入的是我之前创建的map变量的属性（rtype）  
悲催的是，makemap这个方法我同样没有找到实现，在value.go里面。不过还好，通过盘查，在runtime包中的hashmap.go文件里找到一段同名的方法。  
代码很长
```go
func makemap(t *maptype, hint int64, h *hmap, bucket unsafe.Pointer) *hmap {}
```
并且传入有四个参数，和我们用到的`makemap(typ.(*rtype))`似乎又不一样，暂时来说，现在我也不太清楚具体的操作，只能作罢。  

`m := makemap(typ.(*rtype))`而这个方法，大体上应该是重新定义了一个新的map地址，并指给了m。然后返回一个固定的map格式给到调用的变量。  

唉，前面两块内容，真的是这也看不懂，那也看不懂。一路跳过，我们还是看看SetMapIndex()这个方法吧。  
SetMapIndex就给map分配（key,value），要传入的也是key，value两个参数（都需要经过reflect.Typeof()）的转换。  
```go
func (v Value) SetMapIndex(key, val Value) {
	v.mustBe(Map)
	v.mustBeExported()
	key.mustBeExported()
	tt := (*mapType)(unsafe.Pointer(v.typ))
	key = key.assignTo("reflect.Value.SetMapIndex", tt.key, nil)
	var k unsafe.Pointer
	if key.flag&flagIndir != 0 {
		k = key.ptr
	} else {
		k = unsafe.Pointer(&key.ptr)
	}
	if val.typ == nil {
		mapdelete(v.typ, v.pointer(), k)
		return
	}
	val.mustBeExported()
	val = val.assignTo("reflect.Value.SetMapIndex", tt.elem, nil)
	var e unsafe.Pointer
	if val.flag&flagIndir != 0 {
		e = val.ptr
	} else {
		e = unsafe.Pointer(&val.ptr)
	}
	mapassign(v.typ, v.pointer(), k, e)
}
```
这个方法是继承自Value对象的，有两个验证：
```go
v.mustBeExported()
key.mustBeExported()
```
没想到这么快就看到这个不太明白的exported了，在代码中，他是以v/key的flag和常量flagRO进行与运算，如果算下来不是0，就表示不是exported，不允许被调用。 
```
flagStickyRO    flag = 1 << 5
flagEmbedRO     flag = 1 << 6
flagRO          flag = flagStickyRO | flagEmbedRO
表示看不懂@_@
```

