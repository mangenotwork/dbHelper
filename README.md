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
    disablePrepared: false # 是否禁用预编译
    maxIdle: 0 # 最大空闲连接数， 设置0或不设置为默认值
    maxOpen: 0 # 最大连接数， 设置0或不设置为默认值
    maxLife: 0 # 连接最大存活时间 单位ms， 设置0或不设置为默认值
    maxIdleTime: 0 # 连接最大空闲时间 单位ms， 设置0或不设置为默认值
    isSSH: false # true:开启  false:关闭
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
    dbHelper.InitConf("./conf.yaml")
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

```azure
redis:
  - tag: "" # 标记,通过标记获得连接
    host: ""
    port: 6379
    db: 0
    password: ""
    poolSize: 10
    minIdleConn: 5  # 最小空闲连接数
    connMaxIdleTime: 60   # 连接处于空闲状态的最长时间 单位秒
    isSSH: false # true:开启  false:关闭
    sshUser: "" # ssh 账号
    sshPassword: "" # ssh 密码认证; 当SSHPrivateKey同时设置，优先使用密钥认证
    sshPrivateKey: "" # ssh 密钥文件路径
    sshRemoteHost: "" # ssh 服务器地址
    sshRemotePort: 22  # ssh 服务器端口
```

### redis 获取连接

redis 使用的是 github.com/go-redis/redis/v9 库，获取的连接是 *redis.Client

```azure
...
    dbHelper.InitConf("./conf.yaml")
    redisConn := dbHelper.GetRedisConn("btest")
	val, err := redisConn.Conn().Get(context.Background(), "key").Result()
	if err != nil {
		dbHelper.ErrorF("Failed to get key: %v", err)
	}
	dbHelper.Info("Key value:", val)
...
```

### mongoDB 配置

```azure
mongoDB:
  - tag: "" # 标记,通过标记获得连接
    host: ""
    port: 27017
    user: ""
    password: ""
    database: ""
    maxPoolSize: 0     # 连接池最大数 默认100
    maxConnIdleTime: 0 # 连接池中保持空闲的最长时间，单位秒 默认300
    connectTimeout: 0  # 连接超时，单位秒 默认30
    isSSH: false # true:开启  false:关闭
    sshUser: "" # ssh 账号
    sshPassword: "" # ssh 密码认证; 当SSHPrivateKey同时设置，优先使用密钥认证
    sshPrivateKey: "" # ssh 密钥文件路径
    sshRemoteHost: "" # ssh 服务器地址
    sshRemotePort: 22  # ssh 服务器端口
```

### mongoDB 获取连接

使用库 go.mongodb.org/mongo-driver/mongo 保存连接是 *mongo.Database

```azure
...
    dbHelper.InitConf("./conf.yaml")
	mongoConn := dbHelper.GetMongoDBConn("tag")
	collections, err := mongoConn.ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		dbHelper.ErrorF("List collections failed: %v", err)
	}
	dbHelper.Info("Collections:", collections)
...
```

### postgreSQL 配置

```azure
pgsql:
  - tag: "test" # 标记,通过标记获得连接
    user: ""
    password: ""
    host: ""
    port: 0
    database: ""
    maxIdle: 0 # 最大空闲连接数
    maxOpen: 0 # 最大连接数
    maxLife: 0 # 连接最大存活时间 单位ms
    isSSH: false # true:开启  false:关闭
    sshUser: "" # ssh 账号
    sshPassword: "" # ssh 密码认证; 当SSHPrivateKey同时设置，优先使用密钥认证
    sshPrivateKey: "" # ssh 密钥文件路径
    sshRemoteHost: "" # ssh 服务器地址
    sshRemotePort: "" # ssh 服务器端口
```

### postgreSQL 获取连接

使用的 _ "github.com/lib/pq" 库，保存连接是 *sql.DB

```azure
...
    dbHelper.InitConf("./conf.yaml")
	conn := dbHelper.GetPgsqlConn("tag")
	var data map[string]interface{}
	err := conn.QueryRow("select * from user limit 1").Scan(&data)
    if err != nil {
        dbHelper.Error(err)
    }
	dbHelper.Info(data)
...
```

### 对象存储 MinIO 配置

todo...

### 对象存储 MinIO 获取连接

todo...

### 阿里云对象存储 配置

```azure
aliYunOSS:
   - tag: ""      # 标记,通过标记获得连接
     endpoint: "" # OSS访问域名，如：oss-cn-hangzhou.aliyuncs.com
     accessKeyId: ""
     accessKeySecret: ""
     bucketName: ""
```

### 阿里云对象存储 获取连接

- github.com/aliyun/aliyun-oss-go-sdk/oss
- 
```azure
dbHelper.InitConf("./conf.yaml")
bucket := dbHelper.GetAliYunOSSClient("dev")
localFile := "./a.txt"
objectKey := "a.txt"
bucket.PutObjectFromFile(objectKey, localFile)
```

### 腾讯云对象存储 配置
用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
```azure
tencentCOS:
  - tag: "dev" # 标记,通过标记获得连接
    secretId: ""  # secret Id
    secretKey: ""  # secret Key
    bucketUrl: ""  # bucket url
```

### 腾讯云对象存储 获取连接
- github.com/tencentyun/cos-go-sdk-v5
```azure
dbHelper.InitConf("./conf.yaml")
devCos := dbHelper.GetTencentCOSClient("dev")
dbHelper.Info(dbHelper.TencentCOSCheckIsExist(devCos, "a.txt"))
testCos := dbHelper.GetTencentCOSClient("test")
dbHelper.Info(dbHelper.TencentCOSCheckIsExist(testCos, "a.txt"))
```

### excel操作相关辅助函数

todo...

### 图片处理相关辅助函数

todo...

# todo list
- [ok] 配置化   
- [ok] mysql 的连接支持ssh隧道
- [ok] 日志打印
- [ok] 对象存储 腾讯云
- [ok] 常用方法支持，uuid, md5, 字符串处理
- [ok] mongoDB 的连接支持ssh隧道
- [ok] redis 的连接支持ssh隧道
- [ok] postgreSQL 的连接支持ssh隧道
- [ok] 对象存储 阿里云
- 对象存储 MinIO
- excel操作相关的辅助函数
- 图片处理相关辅助函数,压缩,裁剪,水印,缩略图等
