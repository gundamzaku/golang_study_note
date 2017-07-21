## func PtrTo(t Type) Type {}

PtrTo是一个超简单的方法。  
```go
// PtrTo returns the pointer type with element t.
// For example, if t represents type Foo, PtrTo(t) represents *Foo.
func PtrTo(t Type) Type {
	return t.(*rtype).ptrTo()
}
```
看注释就明白，也就是你传个变量过来，它把这个变量变成了指针形式。  
不过我仍然没有弄清楚这个的应用场景到底在什么地方，我照着官方的说明写了一个例子，大概已经能很好的诠释reflect.PtrTo（）的功能了。  
代码如下，主要是转成指针，然后再还原回来：  
```go
func main(){
	var s string
	s = "hello"

	ns:=reflect.PtrTo(reflect.TypeOf(s))
	fmt.Println(ns)

	v:=reflect.New(ns)

	var p *string
	var val string = "hello"
	p = &val

	v.Elem().Set(reflect.ValueOf(p))
	fmt.Println(*(*string)(unsafe.Pointer(v.Elem().Pointer())))
}
result:
*string
hello
```

## func StructOf(fields []StructField) Type {}
比起reflect.PtrTo（）来说，reflect.StructOf（）的源代码就很长了。  
看注释的第一句话，就表达了它的功能  
```go
StructOf returns the struct type containing fields.
也就是说，返回一个Struct类型。
```

由于源代码过长，这里也不打算去看他了（其实也是看不懂），先写一段代码了解一下它的功能，从传入参数来看，他是需要一个StructField类型的数组，这又是什么呢？  
在源代码中有声明，其实这是一个结构体，有一个Name，就表示这个结构体的名称。  
```go
// A StructField describes a single field in a struct.
type StructField struct {
	// Name is the field name.
	Name string
	// PkgPath is the package path that qualifies a lower case (unexported)
	// field name. It is empty for upper case (exported) field names.
	// See https://golang.org/ref/spec#Uniqueness_of_identifiers
	PkgPath string

	Type      Type      // field type
	Tag       StructTag // field tag string
	Offset    uintptr   // offset within struct, in bytes
	Index     []int     // index sequence for Type.FieldByIndex
	Anonymous bool      // is an embedded field
}
```
问题是我怎么产生这个结构体？  

其实也不难，我照着上面的语句说明，直接生成不行么？  
```go
func main(){
	var fields []reflect.StructField
	fields = append(fields,reflect.StructField{
		Name: "S",
		Tag:  "s",
		Type: reflect.TypeOf(""),
	})

	fields = append(fields,reflect.StructField{
		Name: "X",
		Tag:  "x",
		Type: reflect.TypeOf(uint64(0)),
	})
	fmt.Println(fields)
}
result:
[{S  string s 0 [] false} {X  uint64 x 0 [] false}]
```
这其实是一个官方的例子，不过已经确实可以生成一个StructField的结构。至于其它的字段，暂时也用不到它。  
接着我们就可以直接转化了。  
```go
newstruct:=reflect.StructOf(fields)
fmt.Println(newstruct)
result:
struct { S string "s"; X uint64 "x" }
```
一共两种类型，string和unit64。  
接着，我再改一下，让他赋值。  
```go
newstruct:=reflect.StructOf(fields)
val := reflect.New(newstruct).Elem()
val.FieldByName("S").Set(reflect.ValueOf("hello"))
val.FieldByName("X").Set(reflect.ValueOf(uint64(100)))
fmt.Print(val)

result:
{hello 100}
```
需要注意的一点就是，把Type转到Value类型的时候，要使用`reflect.New(newstruct)`方法。  

这个方法就这么完了？如果从粗的来看，差不多就是这么个用法，当然，还有很多细节的内容，比如说type StructField struct {}里面的几个变量。  
前面我们就用了  
Name表示变量的名称  
Type表示变量的类型   
Tag表示变量的标签  
此外还有  
PkgPath，其实就是表示这个变量是哪个包的  
Offset，表示变量在结构体中的顺序  
Index，表示变量用Type.FieldByIndex()查询时对应的索引值 
Anonymous，不是特别清楚，大概和变量的匿名有关系吧  

