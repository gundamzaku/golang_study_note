我们简单的学习一下proto的语法，这样以后就可以自己进行定义了。  

##  定义消息类型  

```go
syntax = "proto3";

message SearchRequest {
  string query = 1;
  int32 page_number = 2;
  int32 result_per_page = 3;
}
```

然后运行编译：

可以看到生成了一个demo.pb.go文件。  

```go
type SearchRequest struct {
	Query         string `protobuf:"bytes,1,opt,name=query" json:"query,omitempty"`
	PageNumber    int32  `protobuf:"varint,2,opt,name=page_number,json=pageNumber" json:"page_number,omitempty"`
	ResultPerPage int32  `protobuf:"varint,3,opt,name=result_per_page,json=resultPerPage" json:"result_per_page,omitempty"`
}
```

而从后面的方法来看，这应该就是一个用来new的类对象。  
```go
func (m *SearchRequest) Reset(){
  *m = SearchRequest{} 
}
```
### 指定字段类型
在上面的例子中，所有的字段都是<a href="https://developers.google.com/protocol-buffers/docs/proto3#scalar">标量类型</a>：两个整形和一个字符。标量类型的意思就是四大类型（布尔、浮点、整型、字符）。然而，你也能给字段用指定的合成类型，比如说<a href="https://developers.google.com/protocol-buffers/docs/proto3#enum">枚举</a>，等…… 

### 分配标签（tags）  
每个字段有一个惟一的数字标签，这是用来识别你在<a href="https://developers.google.com/protocol-buffers/docs/encoding">message二进制格式</a>下的字段的，在使用中没事别去改。注意，标签（tags）如果值在1-15以内的话，解码时按一个byte来算，包含识别数字和字段类型。具体的可以去看<a href="https://developers.google.com/protocol-buffers/docs/encoding#structure">Protocol Buffer Encoding</a>  
从16到2047的话，则要用2个byte来计算了，所以，对于你代码中需要频繁出现的message元素，你应该知道要怎么做了。  
同时，考虑到将来的话，对于你感觉可能会用的元素，预留一些1-15之内的标签。  

最小的标签是1，最大的是 229 - 1, 或者 536,870,911（一般是不会用到这么大的），19000 到 19999这个段内的数字不能被使用（被用做FieldDescriptor::kFirstReservedNumber 到 FieldDescriptor::kLastReservedNumber 的预定义了。）反正用了会报错，同样的，你也不能用之前保留的标签。

### 指定字段规则  

Message字段修饰符有以下几种：
* singular：一个设计良好的message有0或1个这样的字段（但是不能大于1）  
* repeated：这个字段在设计良好的message中被重复多次（包括0），重复值的顺序将被保存。  

在proto3，标量数值类型的repeated字段默认使用packed编码。  
关于packed编码，在 <a href="https://developers.google.com/protocol-buffers/docs/encoding#packed">Protocol Buffer Encoding</a> 中有详细介绍。

### 添加多个Message类型
没什么花头，就是一个.proto文件里面可以有多个message，其中有提到，如果你想定义一个回复的message格式来响应你的SearchResponse，你可以这么写：  
```go
message SearchRequest {
  string query = 1;
  int32 page_number = 2;
  int32 result_per_page = 3;
}

message SearchResponse {
 ...
}
```
### 添加注释  
这真不用多说了，双斜杠（//）

### 保留字段  
如果你通过整个删除或者注释字段来更新message类型，将来的用户在他们对此类型进行更新操作时可以重用标签号码（tag number）。如果他后来加载同样的.proto文件的旧版本时会引发很严重的问题。包括数据错误，隐私错误等。一种方式确认这种情况不会发生就是指定你删除的字段的字段标签tags（和/或者名字，在JSON序列中它也能引发问题）是被保留的。协议编译器（protocol buffer compiler）将会在将来任何一个用户尝试使用这些字段标识符的时候会进行警告。

```go
message Foo {
  reserved 2, 15, 9 to 11;
  reserved "foo", "bar";
}
```
注意你不能在同一个reserved 声明中混合名字和标签码。  

### 如何产生你的.proto文件  

可以跳过了，反正就是用 <a href="https://developers.google.com/protocol-buffers/docs/proto3#generating">protocol buffer compiler</a>编译的时候，不同的语言产生不同类型的文件。知道一下就行了。  

## 标量值类型  
这只是定义了各种不同语言的一些标准，太多了，没什么花头，具体可以直接看文档。  

## 默认值  

对于strings, 默认值为空.  
对于bytes, 默认值也是空.  
对于bools, 默认值是false.  
对于numeric类型, 默认值是0.  
对于enums，默认值是第一个定义的枚举值，且必须为0.  
对于message字段，字段不设置，它对语言依赖比较严格，具体的可以看各种语言的文档指导.  

其实也没有什么花头，这些都是常识。  

## 枚举类型  
没花头，原文的例子：  
```go
message SearchRequest {
  string query = 1;
  int32 page_number = 2;
  int32 result_per_page = 3;
  enum Corpus {
    UNIVERSAL = 0;
    WEB = 1;
    IMAGES = 2;
    LOCAL = 3;
    NEWS = 4;
    PRODUCTS = 5;
    VIDEO = 6;
  }
  Corpus corpus = 4;
}
```

转化以后可以看到：  
```go
const (
	SearchRequest_UNIVERSAL SearchRequest_Corpus = 0
	SearchRequest_WEB       SearchRequest_Corpus = 1
	SearchRequest_IMAGES    SearchRequest_Corpus = 2
	SearchRequest_LOCAL     SearchRequest_Corpus = 3
	SearchRequest_NEWS      SearchRequest_Corpus = 4
	SearchRequest_PRODUCTS  SearchRequest_Corpus = 5
	SearchRequest_VIDEO     SearchRequest_Corpus = 6
)

var SearchRequest_Corpus_name = map[int32]string{
	0: "UNIVERSAL",
	1: "WEB",
	2: "IMAGES",
	3: "LOCAL",
	4: "NEWS",
	5: "PRODUCTS",
	6: "VIDEO",
}
var SearchRequest_Corpus_value = map[string]int32{
	"UNIVERSAL": 0,
	"WEB":       1,
	"IMAGES":    2,
	"LOCAL":     3,
	"NEWS":      4,
	"PRODUCTS":  5,
	"VIDEO":     6,
}
```

产生了一套常量，两个map  

有两个要注意的地方：
一，必须有个零值，也就是`UNIVERSAL = 0;`，你改成别的数字，转化的时候就会报错`demo.proto: The first enum value must be zero in proto3.`
二，零值必须放在第一个，不然还是会报错。用proto2的语法来编译，第一个枚举值总是固定的。

枚举类型可以设置相同的值，必须加上一段代码`option allow_alias = true;`  
不加会报错`demo.proto: "SearchRequest.IMAGES" uses the same enum value as "SearchRequest.WEB". If this is intended, set 'option allow_alias = true;' to the enum definition.`  
加了就正常  
```go
enum Corpus {
    option allow_alias = true;
    UNIVERSAL = 0;
    WEB = 1;
    IMAGES = 1;
    LOCAL = 3;
    NEWS = 4;
    PRODUCTS = 5;
    VIDEO = 6;
}
```

基本上也就这些花头了，可以跳过了，如果碰到更深的问题，可以回过头来再看一下文档。  

## 使用其它的消息类型  

一种方式，用repet  
```go
message SearchResponse {
  repeated Result results = 1;
}

message Result {
  string url = 1;
  string title = 2;
  repeated string snippets = 3;
}
```
可以看到，生成的  
```go
type SearchResponse struct {
	Results []*Result `protobuf:"bytes,1,rep,name=results" json:"results,omitempty"`
}
```
里面包含了一个Result的结构体指针数组（好复杂）  

Result仍然在下面有定义，是一个结构体。  
值得注意的是，看到`Snippets []string`，原来repet的意思就是……嗯，就是产生数组吧。  
```go
type Result struct {
	Url      string   `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
	Title    string   `protobuf:"bytes,2,opt,name=title" json:"title,omitempty"`
	Snippets []string `protobuf:"bytes,3,rep,name=snippets" json:"snippets,omitempty"`
}
```

### 导入  
这个更简单了，就是：  
```go
import "myproject/other_protos.proto";
```
这样的话，可以把proto拆分成多个文件，不过一般这是要搞大项目才用到了。  

import还有一个属性：  
```import public "new.proto";```
没错，就是这个public，什么意思呢？  

说是有时候你会要把.proto文件移到一个新的地方去，文件移走了，那以前的代码的import就出问题了。  
所以你得在原来的地方再弄一个.proto文件做一个新的指向。  
这听起来似乎也没什么卵用么。  

```
举个例子  
new.proto新文件
所有的代码放在这里了

old.proto老文件
import public "new.proto";这里就重定向了
import "other.proto";

client.proto客户端文件
import "old.proto";
可以看到other.proto文件就没用了。（干嘛不直接删除了？）
```

### 使用proto2消息类型  
这个不用看了，没卵用，真有用到的时候再看也不迟。  

