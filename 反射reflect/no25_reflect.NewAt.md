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
```

