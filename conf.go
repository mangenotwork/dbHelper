package dbHelper

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

var Conf conf

// InitConf 初始化读取配置文件并连接各个配置，输入配置文件路径，路径是当前工作目录的相对路径
func InitConf(path string) {
	workPath, _ := os.Getwd()
	configPath := filepath.Join(workPath, path)

	if !FileExists(configPath) {
		panic("未找到配置文件! " + configPath)
	}

	config, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic("读取配置失败! " + err.Error())
	}

	err = yaml.Unmarshal(config, &Conf)
	if err != nil {
		panic("读取配置失败! " + err.Error())
	}

	if len(Conf.MysqlConf) > 0 {
		initMysqlConn()
	}

	if len(Conf.TenCentCOS) > 0 {
		initTencentCOSClient()
	}

	if len(Conf.MongoDBConf) > 0 {
		initMongoDBConn()
	}

	if len(Conf.RedisConf) > 0 {
		initRedisConn()
	}

	if len(Conf.PgsqlConf) > 0 {
		initPgsqlConn()
	}

	if len(Conf.AliYunOSS) > 0 {
		initAliYunOSSClient()
	}

}

type conf struct {
	MysqlConf   []*MysqlConf   `yaml:"mysql"`
	TenCentCOS  []*TenCentCOS  `yaml:"tencentCOS"`
	MongoDBConf []*MongoDBConf `yaml:"mongoDB"`
	RedisConf   []*RedisConf   `yaml:"redis"`
	PgsqlConf   []*PgsqlConf   `yaml:"pgsql"`
	AliYunOSS   []*AliYunOSS   `yaml:"aliYunOSS"`
}

type MysqlConf struct {
	Tag             string `yaml:"tag"` // 标记,通过标记获得连接
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	Host            string `yaml:"host"`
	Port            int64  `yaml:"port"`
	Database        string `yaml:"database"`
	DisablePrepared bool   `yaml:"disablePrepared"` // 是否禁用预编译
	MaxIdle         int64  `yaml:"maxIdle"`         // 最大空闲连接数
	MaxOpen         int64  `yaml:"maxOpen"`         // 最大连接数
	MaxLifeTime     int64  `yaml:"maxLife"`         // 连接最大存活时间 单位ms
	MaxIdleTime     int64  `yaml:"maxIdleTime"`     // 连接最大空闲时间 单位ms
	IsSSH           bool   `yaml:"isSSH"`           // t:开启  f:关闭
	SSHUsername     string `yaml:"sshUser"`         // ssh 账号
	SSHPassword     string `yaml:"sshPassword"`     // ssh 密码认证; 当SSHPrivateKey同时设置，优先使用密钥认证
	SSHPrivateKey   string `yaml:"sshPrivateKey"`   // ssh 密钥文件路径
	SSHRemoteHost   string `yaml:"sshRemoteHost"`   // ssh 服务器地址
	SSHRemotePort   int64  `yaml:"sshRemotePort"`   // ssh 服务器端口
}

// TenCentCOS 腾讯对象存储连接配置
// 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
type TenCentCOS struct {
	Tag       string `yaml:"tag"`       // 标记,通过标记获得连接
	SecretId  string `yaml:"secretId"`  // secret Id
	SecretKey string `yaml:"secretKey"` // secret Key
	BucketURL string `yaml:"bucketUrl"` // bucket url
}

type MongoDBConf struct {
	Tag             string `yaml:"tag"` // 标记,通过标记获得连接
	Host            string `yaml:"host"`
	Port            int64  `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	MaxPoolSize     int64  `yaml:"maxPoolSize"`     // 连接池最大数 默认100
	MaxConnIdleTime int64  `yaml:"maxConnIdleTime"` // 连接池中保持空闲的最长时间，单位秒 默认300
	ConnectTimeout  int64  `yaml:"connectTimeout"`  // 连接超时，单位秒 默认30
	IsSSH           bool   `yaml:"isSSH"`           // t:开启  f:关闭
	SSHUsername     string `yaml:"sshUser"`         // ssh 账号
	SSHPassword     string `yaml:"sshPassword"`     // ssh 密码认证; 当SSHPrivateKey同时设置，优先使用密钥认证
	SSHPrivateKey   string `yaml:"sshPrivateKey"`   // ssh 密钥文件路径
	SSHRemoteHost   string `yaml:"sshRemoteHost"`   // ssh 服务器地址
	SSHRemotePort   int64  `yaml:"sshRemotePort"`   // ssh 服务器端口
}

type RedisConf struct {
	Tag             string `yaml:"tag"` // 标记,通过标记获得连接
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	DB              int    `yaml:"db"`
	Password        string `yaml:"password"`
	PoolSize        int    `yaml:"poolSize"`
	MinIdleConn     int    `yaml:"minIdleConn"`     // 最小空闲连接数
	ConnMaxIdleTime int    `yaml:"connMaxIdleTime"` // 连接处于空闲状态的最长时间 单位秒
	IsSSH           bool   `yaml:"isSSH"`           // t:开启  f:关闭
	SSHUsername     string `yaml:"sshUser"`         // ssh 账号
	SSHPassword     string `yaml:"sshPassword"`     // ssh 密码认证; 当SSHPrivateKey同时设置，优先使用密钥认证
	SSHPrivateKey   string `yaml:"sshPrivateKey"`   // ssh 密钥文件路径
	SSHRemoteHost   string `yaml:"sshRemoteHost"`   // ssh 服务器地址
	SSHRemotePort   int64  `yaml:"sshRemotePort"`   // ssh 服务器端口
}

type PgsqlConf struct {
	Tag           string `yaml:"tag"` // 标记,通过标记获得连接
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	Host          string `yaml:"host"`
	Port          int64  `yaml:"port"`
	Database      string `yaml:"database"`
	MaxIdle       int64  `yaml:"maxIdle"`       // 最大空闲连接数
	MaxOpen       int64  `yaml:"maxOpen"`       // 最大连接数
	MaxLifeTime   int64  `yaml:"maxLife"`       // 连接最大存活时间 单位ms
	IsSSH         bool   `yaml:"isSSH"`         // t:开启  f:关闭
	SSHUsername   string `yaml:"sshUser"`       // ssh 账号
	SSHPassword   string `yaml:"sshPassword"`   // ssh 密码认证; 当SSHPrivateKey同时设置，优先使用密钥认证
	SSHPrivateKey string `yaml:"sshPrivateKey"` // ssh 密钥文件路径
	SSHRemoteHost string `yaml:"sshRemoteHost"` // ssh 服务器地址
	SSHRemotePort int64  `yaml:"sshRemotePort"` // ssh 服务器端口
}

type AliYunOSS struct {
	Tag             string `yaml:"tag"`      // 标记,通过标记获得连接
	Endpoint        string `yaml:"endpoint"` // OSS访问域名，如：oss-cn-hangzhou.aliyuncs.com
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	BucketName      string `yaml:"bucketName"`
}
