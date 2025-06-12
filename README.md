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

```azure
minio:
  - tag: ""      # 标记,通过标记获得连接
    endpoint: "" # MinIO 服务器地址
    accessKeyId: "" # 访问密钥 ID
    accessKeySecret: "" # 秘密访问密钥
    useSSL: false # 是否使用 SSL
```

### 对象存储 MinIO 获取连接

- github.com/minio/minio-go/v7   获取的是 *minio.Client
- 
```azure
dbHelper.InitConf("./conf.yaml")
client := dbHelper.GetMinioClient("dev")
client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
...
```

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

### 常用辅助函数
```azure
dbHelper.ID() int64  // 生成雪花id
dbHelper.IDMd5() string // 生成雪花id MD5
dbHelper.GetMD5Encode(data string) string // MD5编码
dbHelper.NowTimestampStr() string // 当前时间戳字符串
dbHelper.FileMd5sum(fileName string) string // 文件MD5
dbHelper.GetAllFile(pathname string) ([]string, error) // 获取指定目录下的所有文件
dbHelper.RandomIntCaptcha(captchaLen int) string // 生成 captchaLen 位随机数，理论上会重复
dbHelper.DeepEqual(a, b interface{}) bool // DeepEqual 深度比较任意类型的两个变量的是否相等,类型一样值一样反回true, 如果元素都是nil，且类型相同，则它们是相等的; 如果它们是不同的类型，它们是不相等的
dbHelper.SliceContains[V comparable](a []V, v V) bool // 判断切片a中是否包含元素v
dbHelper.SliceDeduplicate[V comparable](a []V) []V // 去重
dbHelper.AnyToString(i interface{}) string // AnyToString any -> string
dbHelper.JsonToMap(str string) (map[string]interface{}, error) // JsonToMap json -> map
dbHelper.MapToJson(m interface{}) (string, error) // MapToJson map -> json
dbHelper.AnyToMap(data interface{}) map[string]interface{} // AnyToMap interface{} -> map[string]interface{}
dbHelper.AnyToInt(data interface{}) int // AnyToInt interface{} -> int
dbHelper.AnyToInt64(data interface{}) int64 // AnyToInt64 interface{} -> int64
dbHelper.AnyToFloat64(data interface{}) float64 // AnyToFloat64 interface{} -> float64
dbHelper.AnyToStrings(data interface{}) []string // AnyToStrings interface{} -> []string
dbHelper.AnyToJson(data interface{}) (string, error) // AnyToJson interface{} -> json string
dbHelper.AnyToJsonB(data interface{}) ([]byte, error) // AnyToJsonB interface{} -> json string
dbHelper.IntToHex(i int) string // IntToHex int -> hex
dbHelper.Int64ToHex(i int64) string // Int64ToHex int64 -> hex
dbHelper.HexToInt(s string) int // HexToInt hex -> int
dbHelper.HexToInt64(s string) int64 // HexToInt64 hex -> int
dbHelper.StrNumToInt64(str string) int64 // StrNumToInt64 string -> int64
dbHelper.StrNumToInt(str string) int // StrNumToInt string -> int
dbHelper.StrNumToInt32(str string) int32 // StrNumToInt32 string -> int32
dbHelper.StrNumToFloat64(str string) float64 // StrNumToFloat64 string -> float64
dbHelper.StrNumToFloat32(str string) float32 // StrNumToFloat32 string -> float32
dbHelper.Uint8ToStr(bs []uint8) string // Uint8ToStr []uint8 -> string
dbHelper.StrToByte(s string) []byte // StrToByte string -> []byte
dbHelper.ByteToStr(b []byte) string // ByteToStr []byte -> string
dbHelper.BoolToByte(b bool) []byte // BoolToByte bool -> []byte
dbHelper.ByteToBool(b []byte) bool // ByteToBool []byte -> bool
dbHelper.IntToByte(i int) []byte // IntToByte int -> []byte
dbHelper.ByteToInt(b []byte) int // ByteToInt []byte -> int
dbHelper.Int64ToByte(i int64) []byte // Int64ToByte int64 -> []byte
dbHelper.ByteToInt64(b []byte) int64 // ByteToInt64 []byte -> int64
dbHelper.Float32ToByte(f float32) []byte // Float32ToByte float32 -> []byte
dbHelper.Float32ToUint32(f float32) uint32 // Float32ToUint32 float32 -> uint32
dbHelper.ByteToFloat32(b []byte) float32 // ByteToFloat32 []byte -> float32
dbHelper.Float64ToByte(f float64) []byte // Float64ToByte float64 -> []byte
dbHelper.Float64ToUint64(f float64) uint64 // Float64ToUint64 float64 -> uint64
dbHelper.ByteToFloat64(b []byte) float64 // ByteToFloat64 []byte -> float64
dbHelper.StructToMap(obj interface{}) map[string]interface{} // StructToMap  struct -> map[string]interface{}
dbHelper.ByteToBit(b []byte) []uint8 // ByteToBit []byte -> []uint8 (bit)
dbHelper.BitToByte(b []uint8) []byte // BitToByte []uint8 -> []byte
dbHelper.ByteToBinaryString(data byte) (str string) // ByteToBinaryString  字节 -> 二进制字符串
dbHelper.MapStrToAny(m map[string]string) map[string]interface{} // MapStrToAny map[string]string -> map[string]interface{}
dbHelper.ByteToGBK(strBuf []byte) []byte // ByteToGBK   byte -> gbk byte
    
还有其余的没有那么常用的方法...

```


### excel操作相关辅助函数

todo...

### 图片处理相关辅助函数

todo...

# todo list
- excel操作相关的辅助函数
- 图片处理相关辅助函数,压缩,裁剪,水印,缩略图等
