# RabbitMQ 用法

## 引用库地址
[Amqp](https://github.com/streadway/amqp)

## 消费者
### config/dev.conf新增

``` json
{
    "amqpx": {
        "type" : 2,
        "uri": "amqp://[name]:[pwd]@[ip]:[port]/xx-service",
        "exchange": "xxw.amqp.direct.push_data",
        "queue": "push_queue",
        "routing_key": "direct.push_record",
        "consumer_tag": "xx-service1",
        "exchange_type":"direct"
    }
}
```

### global/config.go新增

```go
type Cfg struct {
	...
	Amqpx *amqpx.Amqpx `json:"amqpx"`
}
```

### 在适当位置注册方法
e.g. `task/task.go`

``` go
func init() {
	global.Config.Amqpx.ConsumerHandler = SaveData
}

func SaveData(body []byte) error {

    // do something
}
```

## 生产者
### config/dev.conf新增

``` json
{
    "amqpx": {
        "type" : 1,
        "uri": "amqp://[name]:[pwd]@[ip]:[port]/xx-service",
        "exchange": "xxw.amqp.direct.push_data",
        "routing_key": "direct.push_record",
        "exchange_type":"direct"
    }
}
```

### global/config.go新增

```go
type Cfg struct {
	...
	Amqpx *amqpx.Amqpx `json:"amqpx"`
}
```

### 在需要调用的地方

``` go
func DemoFun1() {
    ...
    err:=global.Config.Amqpx.Publish([body])
    if err!=nil{
        // do something
    }
    ...
}
```
