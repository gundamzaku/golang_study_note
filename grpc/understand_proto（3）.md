## Maps  
如果你想创建一个关联map作为你数据定义的一部分，protocol buffers提供了一种便捷的语法：

```go
map<key_type, value_type> map_field = N;
```

`key_type`可以是任意的整形或字符类型（除了浮点指针类型和bytes的任意可变类型）。注意枚举并不是一个有效的key_type。value_type的话，除了其它的map以外的任意类型都是可以的。
