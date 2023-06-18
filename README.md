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

在项目初始化时，可以从配置文件中设置日志级别和日志文件：
``` go
import _ "github.com/jilin7105/ebase"
```

在代码中，当需要打印日志时，可以使用logger：

```go
import "github.com/jilin7105/ebase/logger"

logger.Info("This is an info log.")
logger.Debug("This is a debug log.")
logger.Warn("This is a warning log.")
logger.Error("This is an error log.")

```

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
| 创建gRPC服务 |  |
| 创建任务服务 | ✅ |
| 创建Kafka服务 | ✅ |
| 添加测试 |  |
| 添加文档 |  |

