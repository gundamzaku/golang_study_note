继续前面的话题，在了解了`reflect.ValueOf(newVal)`这个方法是干什么用的以后，我们再回过头来看`reflect.Append`
reflect.Append(s Value, x ...Value)带了两个参数，其实不只两个，这个`x ...Value`表示你可以输入无数个Value类型的x，用","分割。  
append()有一个限制，你传入的s必须是一个slice类型，而后面的x可以是string，也可以是int之类的变量类型。  
然后看一下append()代码的内部  
```go
func Append(s Value, x ...Value) Value {
	s.mustBe(Slice)
	s, i0, i1 := grow(s, len(x))
	for i, j := i0, 0; i < i1; i, j = i+1, j+1 {
		s.Index(i).Set(x[j])
	}
	return s
}
```
一下子感觉好简单，首先这个`s.mustBe(Slice)`，一看就知道是什么意思了。
它是以变量s的kind()属性和常量的Slice数字对比，如果都是23（代表Slice)类型，那么就通过验证。  
`s, i0, i1 := grow(s, len(x))`接下来的这一段，gorw()是一个比较长的函数。
```
/*
 * 传入s切片，x的数量（就是一共传进来几个x)
 * 传出新的切片，两个整数
 */
func grow(s Value, extra int) (Value, int, int) {
	i0 := s.Len() //s的切片数量
	i1 := i0 + extra //要产生的切片的新数量
	if i1 < i0 { //新切片的数量不可能比老切片要少，所以要报错
		panic("reflect.Append: slice overflow")
	}
	m := s.Cap()  //cap()函数返回的是数组切片分配的空间大小。
  //如果空间大小足够，那么直接返回
	if i1 <= m {
		return s.Slice(0, i1), i0, i1
	}
  //空间不足的情况下，如果完全没有空间，m的空间就是extra的数量
	if m == 0 {
		m = extra
	} else {
		//m现在比新的切片大小要小的话
    for m < i1 {
      //如果i0，即老s的切片数量在1024以内，m+m？
			if i0 < 1024 {
				m += m
			} else {
        //超过1024，m+m/4?
				m += m / 4
			}
		}
	}
  //最后产生一个新的切片返回
	t := MakeSlice(s.Type(), i1, m)
	Copy(t, s)
	return t, i0, i1
}
```
