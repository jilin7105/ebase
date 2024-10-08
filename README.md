# ebase

`ebase`是一个使用Go编写的微服务框架，可以通过配置文件切换不同的服务类型。

## 主要功能

- 支持切换服务类型，包括HTTP服务器，gRPC服务器，定时任务，Kafka消费服务。
- 提供心跳检测功能。 (简易实现)
- 提供服务注册功能。 (简易实现)

## 快速开始
### 方式 1
要使用`ebase`，首先需要下载和安装Go。然后，可以使用`go get`命令下载并安装`ebase`：

```bash
go get github.com/jilin7105/ebase
```
### 方式 2

```bash
# 下载并安装工具包
go install github.com/jilin7105/ebasetools 

# 创建一个名为my-api-project的新项目   --type=  "HTTP", "KAFKA", "GRPC", "TASK"  
# TASK => 定时任务项目 
# HTTP => api 项目
# KAFKA => kafka 消费者 项目
# KAFKA => GRPC  grpc 项目
ebasetools create --type=HTTP --name=my-api-project
```

给eb传入当前工作目录， 由于各种原因， 防止ebase获取不到工作目录，或者获取不准确
```go
	path, _ := os.Getwd()

	ebase.SetProjectPath(path)
```



在代码中，当需要打印日志时，可以使用logger：

```go
import "github.com/jilin7105/ebase/logger"

logger.Info("This is an info log.")
logger.Debug("This is a debug log.")
logger.Warn("This is a warning log.")
logger.Error("This is an error log.")

logger.Info("hello world; name:%s; age:%d", "zhangsan", 18)
/* [INFO] xxx/xxx.go:47 hello world; name:zhangsan; age:18 */

// zap - *zap.Logger
logger.Infoz("hello world", zap.String("name", "zhangsan"), zap.Int("age", 18))
/* 一行内,这里为方便展示将json格式化
{
    "level":"info",
    "time":"2024-10-01 12:00:00.000",
    "line":"/zap/zap_test.go:143",
    "func":"zap_test.TestFunc",
    "msg":"hello world",
    "name":"zhangsan",
    "age":18
}
*/

// 带语法糖的zap - *zap.SugaredLogger; 性能不敏感场景下使用
logger.Infow("hello world", "name", "zhangsan", "age", 18)
/* 格式同上 */
```

在代码中使用数据库
```go
import "github.com/jilin7105/ebase"
db := ebase.GetDB("db_name")
//此处返回gorm.db对象
if db == nil {
    panic("数据库不存在")
}
type Product struct {
    gorm.Model
    Code  string
    Price uint
}
var product Product
db.First(&product, 1)
```

在代码中使用redis
```go
import "github.com/jilin7105/ebase"
rdb := ebase.GetRedis("redis_name")
if rdb == nil {
    panic("redis 不存在")
}
val, err := rdb.Get(ctx, "score").Result()
if err != nil {
    log.Fatal(err)
}
w.Write([]byte(fmt.Sprintf("score的值: %v", val)))
```

在代码中使用 Kafka
```go
import (
    "github.com/Shopify/sarama"
	"github.com/jilin7105/ebase"
)
kp := *ebase.GetKafka("Producer_name")
if kp == nil {
    panic("KafkaProducer 不存在")
}
topic := "your topic"

msg := &sarama.ProducerMessage{
    Topic: "your topic",
	//newManualPartitioner: true  #是否手动分配分区
	//如果手动分区选择true ，需要手动设置分区 增加如下代码
	//Partition: int32(your_partition_number),  // 设置分区号
    Value: sarama.StringEncoder("Hello World"),
}



partition, offset, err := kp.SendMessage(msg)
if err != nil {
    log.Fatalln("Failed to send message:", err)
}

log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)

```

### 微服务
```yaml
micro : # 微服务相关配置  （非必须）
  is_reg : true  #是否有服务注册
  is_heart_push : true  #是否心跳推送
  heart_push_speed : 5  #心跳推送速度 单位 秒  如果不填写将只执行1次 ，默认用户方法内部处理心跳逻辑
```

都是简单实现， 服务检测到微服务相关配置后 
```go
	//增加心跳推送   未使用 heart_push_speed   go f()  形式执行一次
	//如果使用heart_push_speed   
	//原理 go fun(){ 
	//    for{
	//		f()    //你的方法
	//		time.sleep(heart_push_speed  )
    //    }       
	//}
	ebase.SetHeartbeatPush(func() error {
		log.Println("HeartbeatPush")
		return nil
	})

	//增加服务注册   go f()  形式执行一次
	ebase.SetRegfunc(func() error {
		log.Println("Regfunc")
		return nil
	})
```

### 内置简单 http 请求服务 配合链路追踪 
```go
import (

"github.com/jilin7105/ebase/helpfunc/EBHttpRequest"

"github.com/levigross/grequests"

)

// 使用 github.com/levigross/grequests包进行请求
// https://pkg.go.dev/github.com/levigross/grequests 可以查看文档
EBHttpRequest.Post(context, "http://127.0.0.1:9999/p", &grequests.RequestOptions{
			JSON: map[string]string{
				"name": "jilin",
			},
		})



EBHttpRequest.Get(context, "http://127.0.0.1:9999/f", nil)
```

### Gin 统一返回封装
```go
import "github.com/jilin7105/ebase/helpfunc/EBHttpRequest"
//返回分页数据
EbaseGinResponse.PageOK(c *gin.Context, result any, count int, pageIndex int, pageSize int, msg string)
//常规返回
EbaseGinResponse.OK(c *gin.Context, data any, msg string)
//返回错误信息
EbaseGinResponse.Error(c *gin.Context, code int, err error, msg string)
```

### 配置文件相关 

[配置文件示例，**仅支持yml格式**](https://github.com/jilin7105/ebase/tree/main/ex.config.yml)
```shell
 # 默认config.yml 
 # 可以通过-i 进行指定
 # 配置文件需要在项目根目录下
 go build && ./{你的执行文件名称}  -i  config-online.yml 
```
### 辅助函数
[**辅助函数文档**](https://github.com/jilin7105/ebase/tree/main/doc/helpfunc.md)

### [在代码中 http 服务使用](https://github.com/jilin7105/ebase/tree/main/examp/httpex)
### [在代码中 定时任务 服务使用](https://github.com/jilin7105/ebase/tree/main/examp/task)
### [在代码中 kafka消费 服务使用](https://github.com/jilin7105/ebase/tree/main/examp/kafka)



| 任务 | 完成 |
| --- | --- |
| 创建项目 | ✅ |
| 创建配置文件 | ✅ |
| 解析命令行参数 | ✅ |
| 加载配置文件 | ✅ |
| 初始化Eb结构体 | ✅ |
| 创建HTTP服务 | ✅ |
| 创建gRPC服务 | ✅ |
| 创建任务服务 | ✅ |
| 创建Kafka服务 | ✅ |
| 添加测试 | ✅ |
| 添加文档 | ✅ |

