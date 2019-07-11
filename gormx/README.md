## Gorm 用法

### 环境变量
可以直接从系统环境变量`DB_MYSQL`中进行获取的
e.g.
```sh
user:password@/dbname?charset=utf8mb4&parseTime=True&loc=Local
```
或将值直接存放进配置文件中，见下`dev.conf`（*建议）

### config/dev.conf新增
``` json
{
    "gormx": {
        "max_open_conns": 20,
        "user": "**",
        "password": "**",
        "db_name": "**",
        "host": "**",
        "port": 3306
    }
}
```

### global/config.go新增

```go
type Cfg struct {
	...
	Gormx *gormx.Gormx `json:"gormx"`
}
```

### 在适当位置注册方法
申明表结构，直接使用连接即可，详细用法点[这里](https://jasperxu.github.io/gorm-zh/)