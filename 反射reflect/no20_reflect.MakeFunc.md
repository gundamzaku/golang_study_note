## func MakeFunc(typ Type, fn func(args []Value) (results []Value)) Value {}

MakeFunc的代码不长

```go
func MakeFunc(typ Type, fn func(args []Value) (results []Value)) Value {
	if typ.Kind() != Func {
		panic("reflect: call of MakeFunc with non-Func type")
	}

	t := typ.common()
	ftyp := (*funcType)(unsafe.Pointer(t))

	// Indirect Go func value (dummy) to obtain
	// actual code address. (A Go func value is a pointer
	// to a C function pointer. https://golang.org/s/go11func.)
	dummy := makeFuncStub
	code := **(**uintptr)(unsafe.Pointer(&dummy))

	// makeFuncImpl contains a stack map for use by the runtime
	_, _, _, stack, _ := funcLayout(t, nil)

	impl := &makeFuncImpl{code: code, stack: stack, typ: ftyp, fn: fn}

	return Value{t, unsafe.Pointer(impl), flag(Func)}
}
```

typ.common()就是抽取的typ的rtpye。  
然后把这个指针转成了funcType类型。  
```go
type funcType struct {
	rtype    `reflect:"func"`//这个是tag
	inCount  uint16
	outCount uint16 // top bit is set if last input parameter is ...
}
```
```go
dummy := makeFuncStub
```
makeFuncStub到底做了什么？恐怕并不是这么好能理解的了，官方给了一个详细的文档  
https://golang.org/s/go11func

来阐述原理，同时，也有相应的注释说明  
```
// makeFuncStub is an assembly function that is the code half of
// the function returned from MakeFunc. It expects a *callReflectFunc
// as its context register, and its job is to invoke callReflect(ctxt, frame)
// where ctxt is the context register and frame is a pointer to the first
// word in the passed-in argument frame.
```
好吧……其实看不懂。估计一时半会也很难理解的了。暂时先跳过吧  
```go
_, _, _, stack, _ := funcLayout(t, nil)
```
funcLayout(t, nil)就是传统的上锁-读cache-没有就生成-解锁-返回的流程。  
考虑到这一块的代码过于复杂，就不在这里强行解读了，等以后有进一步的深入再行分析。  
```go
impl := &makeFuncImpl{code: code, stack: stack, typ: ftyp, fn: fn}
```
则是拼装成一个struct  
```
type makeFuncImpl struct {
	code  uintptr
	stack *bitVector
	typ   *funcType
	fn    func([]Value) []Value
}
```
最后以Value的类型返回。  

