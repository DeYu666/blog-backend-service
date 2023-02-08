# 进行中
- 今天将 转载到 csdn 的工作做完了。
- 

# 用到的库
- json 转 html：用了 blackfriday 


# 其他
有时间可以啃啃 net/http 的源码，感觉很有用，


# macOS 中编译 linux 中的可执行 golang 文件
## 编译为 linux 执行程序
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

## 编译为 Windows 执行程序
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go


# 环境变量
VIPER_CONFIG 指的是 config.yaml 文件的绝对路径
LOG_ROOT_DIR 指的是 存放日志的路径
DB_HOST 指的是 mysql 数据库所处的 ip 地址
GIN_MODE 指的是 gin 框架的 model ，如果是 debug 模式，控制台会输出很多日志，如果是 release，控制台会很干净

示例如下：

```
VIPER_CONFIG=/usr/src/blog_end_go/config.yaml 
LOG_ROOT_DIR=/data/log 
DB_HOST=172.17.0.2
GIN_MODE=release
```