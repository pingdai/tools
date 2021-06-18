# tools
旨在用最简单的轮子造出最强大的兵器。
## 支持
- [x] 三方依赖项组件化
- [x] 服务框架自动生成
- [x] gorm常用增删改查语句自动生成
- [x] 自动注册pprof接口  

## 简介
通过最简单的命令来启动一个新服务，服务统一化、规范化.让业务开发者专注于业务发开。

### 公共环境变量
#### 环境标志`ENV_FLAG`
- 测试环境值为：test
- 预发布环境值为：pre
- (如果有的话可以加上)

## 用法
### 安装
```sh
go get -u -v github.com/pingdai/tools
cd $GOPATH/src/github.com/pingdai/tools
go install
```
会生成一个tools工具，当前已支持的功能有
- 自动生成代码框架
- gorm自动生成
- 默认使用logrus作为日志库

```sh
tools

Usage:
  tools [command]

Available Commands:
  gen         generators
  help        Help about any command
  new         new service

Flags:
  -h, --help   help for tools

Use "tools [command] --help" for more information about a command.

```
#### tools new
`tools new`命令主要是用来一键初始化项目，会直接创建生成一个开箱即用的项目，如：
```shell
➜ ./tools new -h
new service

Usage:
  tools new [flags]

Flags:
      --db-name string   with db name
  -h, --help             help for new

➜ ./tools new demo
2021/06/18 14:14:23 Generated file to /Users/pingd/Go/src/github.com/pingdai/tools/test/demo/doc.go (0 KiB)
2021/06/18 14:14:24 Generated file to /Users/pingd/Go/src/github.com/pingdai/tools/test/demo/global/config.go (0 KiB)
2021/06/18 14:14:25 Generated file to /Users/pingd/Go/src/github.com/pingdai/tools/test/demo/routes/root.go (0 KiB)
2021/06/18 14:14:26 Generated file to /Users/pingd/Go/src/github.com/pingdai/tools/test/demo/main.go (0 KiB)
2021/06/18 14:14:26 Generated file to /Users/pingd/Go/src/github.com/pingdai/tools/test/demo/constants/types/doc.go (0 KiB)
2021/06/18 14:14:26 Generated file to /Users/pingd/Go/src/github.com/pingdai/tools/test/demo/constants/errors/doc.go (0 KiB)
2021/06/18 14:14:26 Generated file to /Users/pingd/Go/src/github.com/pingdai/tools/test/demo/modules/doc.go (0 KiB)
2021/06/18 14:14:26 Generated file to /Users/pingd/Go/src/github.com/pingdai/tools/test/demo/config/dev.conf (0 KiB)

➜ cd demo && ls
config    constants doc.go    global    main.go   modules   routes

```
#### tools gen
`tools gen`命令主要是用来配合gorm，定义一个gorm的model后，会自动生成一系列CRUD接口
```shell
➜ ./tools gen model -h
generate gorm db model method

Usage:
  tools gen model [flags]

Flags:
  -h, --help                help for model
  -t, --table-name string   custom table name

➜ cat b.go
package b

//go:generate gen model -t t_b BB
// @def primary ID
// @def index i_n_a Name Age
// @def unique_index u_no No
type BB struct{
	ID uint64 `gorm:"id"`
	Name string `gorm:"name"`
	Age int `gorm:"age"`
	No int `gorm:"no"`
}

➜ ./tools gen model -t t_b BB
主键索引:{Name:ID Type:uint64 DbFieldName: IsEnable:false IsCreateTime:false}
普通索引:i_n_a
	{Name:Name Type:string DbFieldName: IsEnable:false IsCreateTime:false}
	{Name:Age Type:int DbFieldName: IsEnable:false IsCreateTime:false}
唯一索引:u_no
	{Name:No Type:int DbFieldName: IsEnable:false IsCreateTime:false}
Field:	&{Name:ID Type:uint64 DbFieldName: IsEnable:false IsCreateTime:false}
Field:	&{Name:Name Type:string DbFieldName: IsEnable:false IsCreateTime:false}
Field:	&{Name:Age Type:int DbFieldName: IsEnable:false IsCreateTime:false}
Field:	&{Name:No Type:int DbFieldName: IsEnable:false IsCreateTime:false}
2021/06/18 11:26:00 Generated file to /Users/pingd/Go/src/github.com/pingdai/tools/test/b_b__generated.go (3 KiB)

➜ cat b_b__generated.go
package b

import (
	"github.com/jinzhu/gorm"
)

type BBList []BB

func (bb BB) TableName() string {
	table_name := "t_b"
	return table_name
}

func (bbl *BBList) BatchFetchByIDList(db *gorm.DB, iDList []uint64) error {
	if len(iDList) == 0 {
		return nil
	}

	err := db.Table(BB{}.TableName()).Where(" in (?)", iDList).Find(bbl).Error

	return err
}

func (bbl *BBList) BatchFetchByNameList(db *gorm.DB, nameList []string) error {
	if len(nameList) == 0 {
		return nil
	}

	err := db.Table(BB{}.TableName()).Where(" in (?)", nameList).Find(bbl).Error

	return err
}

func (bbl *BBList) BatchFetchByNoList(db *gorm.DB, noList []int) error {
	if len(noList) == 0 {
		return nil
	}

	err := db.Table(BB{}.TableName()).Where(" in (?)", noList).Find(bbl).Error

	return err
}

func (bb *BB) Create(db *gorm.DB) error {

	err := db.Table(bb.TableName()).Create(bb).Error

	return err
}

func (bb *BB) DeleteByID(db *gorm.DB) (err error) {

	err = db.Table(bb.TableName()).Where(" = ?", bb.ID).Delete(bb).Error

	return err
}

func (bb *BB) DeleteByNo(db *gorm.DB) (err error) {

	err = db.Table(bb.TableName()).Where(" = ?", bb.No).Delete(bb).Error

	return err
}

func (bb *BB) FetchByID(db *gorm.DB) error {

	err := db.Table(bb.TableName()).Where(" = ?", bb.ID).Find(bb).Error

	return err
}

func (bb *BB) FetchByIDForUpdate(db *gorm.DB) error {

	err := db.Table(bb.TableName()).Where(" = ?", bb.ID).Set("gorm:query_option", "FOR UPDATE").Find(bb).Error

	return err
}

func (bbl *BBList) FetchByName(db *gorm.DB, name string) error {

	err := db.Table(BB{}.TableName()).Where(" = ?", name).Find(bbl).Error

	return err
}
...
```
备注：
- primary 声明为主键索引
- index 声明为普通索引
- unique_index 声明为唯一索引

#### Gorm接入及用法
Gorm当前已经集成到自动化生成工具中，可以自动生成，也可以手动注册，详细用法点[这里](./tools/gormx/README.md)

#### RabbitMQ接入及用法
[这里](./tools/amqpx/README.md)

#### Redis client 接入及用法
[这里](./tools/redisx/README.md)