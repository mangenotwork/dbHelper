# dbHelper
db助手，解决高效跨平台跨服务处理数据程序编写的帮助支持库，配置管理各种连接;
配置使用的yaml

# 使用

```azure
go get github.com/mangenotwork/dbHelper
```

### mysql 配置
```azure
mysql:
  - tag: "" # 标记,通过标记获得连接
    user: "root" # 用户名
    password: "" # 密码
    host: "" # mysql主机
    port: 3306 # mysql端口
    database: "" # 数据库名
    disablePrepared : false # 是否禁用预编译
    maxIdle: 0 # 最大空闲连接数， 设置0或不设置为默认值
    maxOpen: 0 # 最大连接数， 设置0或不设置为默认值
    maxLife: 0 # 连接最大存活时间 单位ms， 设置0或不设置为默认值
    maxIdleTime: 0 # 连接最大空闲时间 单位ms， 设置0或不设置为默认值
    isSSH: true # t:开启  f:关闭
    sshUser: "" # ssh 账号
    sshPassword: "" # ssh 密码认证; 当SSHPrivateKey同时设置，优先使用密钥认证
    sshPrivateKey: "" # ssh 密钥文件路径
    sshRemoteHost: "" # ssh 服务器地址
    sshRemotePort: 22  # ssh 服务器端口
```

### mysql 获取连接

mysql使用的是gorm.io/gorm库，获取的连接是*gorm.DB

```azure
...
	conn := dbHelper.GetMysqlConn("tag")
	var data map[string]interface{}
	err := conn.Raw("select * from user limit 1").Scan(&data).Error
	if err != nil {
		dbHelper.Error(err)
	}
	dbHelper.Info(data)
...
```

### redis 配置

todo...

### redis 获取连接

todo...

### mongoDB 配置

todo...

### mongoDB 获取连接

todo...


# todo list
- [ok] 配置化   
- [ok] mysql 的连接支持ssh隧道
- [ok] 日志打印
- 对象存储 腾讯云
- 常用方法支持，uuid, md5, 字符串处理
- redis 的连接支持ssh隧道
- postgreSQL 的连接支持ssh隧道
- mongoDB 的连接支持ssh隧道
- 对象存储 阿里云
- 对象存储 MinIO
