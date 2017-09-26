# 安装

rpcx的安装是很简单的，只要执行：  
`go get -u -v github.com/smallnest/rpcx/...`  
即可。不过由于众所周知的原因，谷歌在国内是无法访问的，所以在安装的过程中，会有大量调用golang.com网站的安装包。所以按正常思路来说是无法安装成功的。  

因此，我们一种方法就是翻墙，非常简单。  
如果不能翻墙，那么就用第二种方法。Go放在github的备份。  

一个是`https://github.com/golang/crypto`,一个是`github.com/golang/net`  
可以在用 `go get github.com/golang/net` 之后，做一个软连接到需要用到的目录中去。  
<b>example</b>:
`mklink /d "D:\go_code\src\golang.org\x\crypto"  "D:\go_code\src\github.com/golang/crypto"`  
`mklink /d "D:\go_code\src\golang.org\x\net\crypto"  "D:\go_code\src\github.com/golang/crypto"`

之后再重复最上面的get即可。  
我在安装的时候报了个错，
```
# cd D:\go_code\src\golang.org\x\crypto; git config remote.origin.url
fatal: Invalid symlink 'D:/go_code/src/golang.org/x/crypto': Function not implemented package golang.org/x/crypto/pbkdf2: exit status 128
# cd D:\go_code\src\golang.org\x\crypto; git config remote.origin.url
fatal: Invalid symlink 'D:/go_code/src/golang.org/x/crypto': Function not implemented package golang.org/x/crypto/salsa20: exit status 128
# cd D:\go_code\src\golang.org\x\crypto; git config remote.origin.url
fatal: Invalid symlink 'D:/go_code/src/golang.org/x/crypto': Function not implemented package golang.org/x/crypto/tea: exit status 128
# cd D:\go_code\src\golang.org\x\crypto; git config remote.origin.url
fatal: Invalid symlink 'D:/go_code/src/golang.org/x/crypto': Function not implemented package golang.org/x/crypto/twofish: exit status 128
# cd D:\go_code\src\golang.org\x\crypto; git config remote.origin.url
fatal: Invalid symlink 'D:/go_code/src/golang.org/x/crypto': Function not implemented package golang.org/x/crypto/xtea: exit status 128
# cd D:\go_code\src\golang.org\x\net; git config remote.origin.url
fatal: Invalid symlink 'D:/go_code/src/golang.org/x/net': Function not implemented package golang.org/x/net/ipv4: exit status 128
```
目前暂时原因不明。  
