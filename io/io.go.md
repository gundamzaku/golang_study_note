在io.go的代码段中，有一段都是用来定义接口的。  
1、type Reader interface {}  
2、type Writer interface {}  
3、type Closer interface {}  
4、type Seeker interface {}  
5、type ReadWriter interface {}  
6、type ReadCloser interface {}  
7、type WriteCloser interface {}  
8、type ReadWriteCloser interface {}  
9、type ReadSeeker interface {}  
10、type WriteSeeker interface {}  
11、type ReadWriteSeeker interface {}  
12、type ReaderFrom interface {}  
13、type WriterTo interface {}  
14、type ReaderAt interface {}  
15、type WriterAt interface {}  
16、type ByteReader interface {}  
17、type ByteScanner interface {}  
18、type ByteWriter interface {}  
19、type RuneReader interface {}  
20、type RuneScanner interface {}  
21、type stringWriter interface {}  
一共21个……这么多。目前来说，我并不知道他们的实际用处。  

接下来，又定义了几组对外的公有函数：   
1、func WriteString(w Writer, s string) (n int, err error) {}  
2、func ReadAtLeast(r Reader, buf []byte, min int) (n int, err error) {}  
3、func ReadFull(r Reader, buf []byte) (n int, err error) {}  
4、func CopyN(dst Writer, src Reader, n int64) (written int64, err error) {}  
5、func Copy(dst Writer, src Reader) (written int64, err error) {}  
6、func CopyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {}  
`下面这个是私有的`  
7、func copyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {}  
8、func LimitReader(r Reader, n int64) Reader { return &LimitedReader{r, n} }  
