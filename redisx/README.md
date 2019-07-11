## Redis Client 接入

### 环境变量
无

### config/dev.conf新增
``` json
{
    "redisx": {
        "addr": "localhost:6379",
        "password": "", // 有就填写，没有的话这个字段也可以省掉
        "db":1  // 连接redis的DB，不设置的话默认使用0
    }
}
```

### global/config.go新增

```go
type Cfg struct {
	...
	Redisx *redisx.Redisx `json:"redisx"`
}
```

### 在适当位置注册方法
详细用法点[这里](https://github.com/go-redis/redis)