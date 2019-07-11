# tools
旨在用最简单的轮子造出最强大的兵器。
## 支持
- [x] 三方依赖项组件化
- [x] 服务框架自动生成
- [x] gorm常用增删改查语句自动生成
- []  

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
#### Gorm接入及用法
Gorm当前已经集成到自动化生成工具中，可以自动生成，也可以手动注册，详细用法点[这里](./tools/gormx/README.md)

#### RabbitMQ接入及用法
[这里](./tools/amqpx/README.md)

#### Redis client 接入及用法
[这里](./tools/redisx/README.md)