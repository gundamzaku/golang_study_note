`原来准备看rpcx的，结果因为google被墙的关系，加上一些莫名其妙的问题，始终没有调顺，故暂时放弃。转看google官方的grpc`  
`详细内容可以参加官方网站(grpc.io)，这里只是一些归纳和学习笔记。`  

# 安装  

## 第一步 安装Go的gRpc包  
```go
$ go get google.golang.org/grpc
```

不过由于国情的关系，上面的明显是无法下载的，所以还是要用下面的地址：  
```go
$ go get github.com/grpc/grpc-go
```

下载好就算是完成了，非常方便，如果只是运行例子的话，直接使用里面自带的DEMO即可。  

在`google.golang.org包中grpc/examples/helloworld`目录内有greeter_client和greeter_server两个目录，分别对应的是客户端和服务端。先不用管具体是怎么用的，也不用管其它的目录文件。先将这两个文件拷贝出来，分别新建两个不同在项目，一个client端，一个server端。  

接着就是先运行server端了，然后运行client端，你会发现一切都已经成功。  

## 第二步 运用protobuf  

这个东西就有些费解了，其实可以把它简单的当成一个文件代码生成器，主要是用来生成一套固定标准的代码的。  
protobuf的相关内容在这里，有好几套标准，分别对应不同的语言，也就是说你如果想用XX语言来与gRpc相联，就要用这个语言的protobuf。  
`https://github.com/google/protobuf/releases`  

可是里面唯独没有Go语言的，其实是在`https://github.com/golang/protobuf`里，用`go get`的方式可以获得。  

这个代码主要是为了生成`protoc-gen-go.exe`的可执行文件的，我用`go get`的时候它自动帮我生成好了。  

要注意的是，这个文件是默认放在你的`GO_PATH`里设置的目录里面的bin目录下面的。这样才能让这个文件可以在命令行模式下面自动执行。如果不对，你可以自行设置一下，比如我，我的`GO_PATH`指定在D:\Go_code目录下面，所以文件也生成在了D:\Go_code\bin下面 
我把它COPY到了我的GO安装目录C:\Go\bin下面  

接着就可以在命令行下面运行protoc命令了。  
```go
protoc --go_out=. *.proto
```

这是修改输出的目录  

