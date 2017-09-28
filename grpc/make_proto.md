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
