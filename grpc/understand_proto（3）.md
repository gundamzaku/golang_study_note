## Maps  
如果你想创建一个关联map作为你数据定义的一部分，protocol buffers提供了一种便捷的语法：

```go
map<key_type, value_type> map_field = N;
```

`key_type`可以是任意的整形或字符类型（除了浮点指针类型和bytes的任意可变类型）。注意枚举并不是一个有效的key_type。value_type的话，除了其它的map以外的任意类型都是可以的。  

举个例子，如果你想创建一个项目的map，每一个项目message均与字符串键值进行关联。你能如下进行定义：  

```go
map<string, Project> projects = 3;
```

* map字段不能重复  
* Wire格式排序和map数值的map迭代排序是未定义的，因此你不能依赖你的map组进行特定的排序。
* 当.proto产生文本格式时，maps已经按key进行了排序。数字键会按数字进行排序。
* 当从wire中解析或当合并时，如果发现有重复的map keys的话，那么用其中的最后一个key。当解析文本格式中的map时，如果有重复的keys的话，解析会失败。

生成map的API目前已经被所有支持proto3的语言可用了。你能在相应的语言的<a href="https://developers.google.com/protocol-buffers/docs/reference/overview">相关文档</a>中找到关于这个map API的内容。（见鬼去吧……）
