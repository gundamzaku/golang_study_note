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
func (m *SearchRequest) Reset()                    { *m = SearchRequest{} }
```

message可以创建多个。  

## 标量值类型  
这只是定义了各种不同语言的一些标准，具体可以直接看文档。  

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

