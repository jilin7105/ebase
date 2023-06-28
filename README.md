# ebase

`ebase`是一个使用Go编写的微服务框架，可以通过配置文件切换不同的服务类型。

## 主要功能

- 支持切换服务类型，包括HTTP服务器，gRPC服务器，定时任务，Kafka消费服务。
- 提供心跳检测功能。 (还未提供)
- 提供服务注册功能。 (还未提供)

## 快速开始

要使用`ebase`，首先需要下载和安装Go。然后，可以使用`go get`命令下载并安装`ebase`：

```bash
go get github.com/jilin7105/ebase
```

在代码中，当需要打印日志时，可以使用logger：

```go
import "github.com/jilin7105/ebase/logger"

logger.Info("This is an info log.")
logger.Debug("This is a debug log.")
logger.Warn("This is a warning log.")
logger.Error("This is an error log.")

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


### 配置文件相关 

[配置文件示例，**仅支持yml格式**](https://github.com/jilin7105/ebase/tree/main/config.yml)
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
| 添加测试 |  |
| 添加文档 |  |

