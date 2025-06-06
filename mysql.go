package dbHelper

import (
	"context"
	_ "database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var MysqlConn map[string]*gorm.DB

func GetMysqlConn(tag string) *gorm.DB {
	m, ok := MysqlConn[tag]
	if !ok {
		panic("[Mysql] 未init")
	}
	return m
}

func initMysqlConn() {
	MysqlConn = make(map[string]*gorm.DB, len(Conf.MysqlConf))
	for _, v := range Conf.MysqlConf {
		m, err := mysqlConn(v)
		if err != nil {
			panic(err)
		}
		MysqlConn[v.Tag] = m
	}
}

func mysqlConn(conf *MysqlConf) (*gorm.DB, error) {
	var (
		orm  *gorm.DB
		err  error
		host = conf.Host
		port = conf.Port
	)

	if conf.IsSSH {
		sshConf := &sshConfig{
			User:       conf.SSHUsername,
			Password:   conf.SSHPassword,
			PrivateKey: conf.SSHPrivateKey,
			RemoteHost: conf.SSHRemoteHost,
			RemotePort: conf.SSHRemotePort,
			TargetHost: conf.Host,
			TargetPort: conf.Port,
		}
		sshConf.LocalHost = "127.0.0.1"
		sshConf.LocalPort, err = getFreePort()
		if err != nil {
			return nil, err
		}

		go func() {
			err = sshConf.getSshConn()
			if err != nil {
				Error("ssh connect err:", err)
				panic(err)
			}
		}()

		// todo 优化 ssh连接成功了通知
		time.Sleep(3 * time.Second)

		host = sshConf.LocalHost
		port = sshConf.LocalPort
	}

	if conf.Database == "" || conf.User == "" || conf.Password == "" || host == "" {
		panic("数据库配置信息获取失败")
	}

	str := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conf.User, conf.Password, host, port, conf.Database) + "?charset=utf8mb4&parseTime=true&loc=Local"
	if conf.DisablePrepared {
		str = str + "&interpolateParams=true"
	}

	orm, err = gorm.Open(mysql.Open(str), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: NewGormLogger(),
	})
	if err != nil {
		return nil, err
	}

	db, err := orm.DB()
	if err != nil {
		return nil, err
	}

	if conf.MaxIdle < 1 {
		conf.MaxIdle = 10
	}
	if conf.MaxOpen < 1 {
		conf.MaxOpen = 100
	}
	if conf.MaxLifeTime < 1 {
		conf.MaxLifeTime = 300000
	}
	if conf.MaxIdleTime < 1 {
		conf.MaxIdleTime = 300000
	}

	db.SetMaxIdleConns(int(conf.MaxIdle)) //空闲连接数
	db.SetMaxOpenConns(int(conf.MaxOpen)) //最大连接数
	db.SetConnMaxLifetime(time.Duration(conf.MaxLifeTime) * time.Millisecond)
	db.SetConnMaxIdleTime(time.Duration(conf.MaxIdleTime) * time.Millisecond) // 连接最大空闲时间

	return orm, err
}

type GormLogger struct {
	SlowThreshold time.Duration
	//Level         gormLogger.LogLevel
}

var _ gormLogger.Interface = (*GormLogger)(nil)

func NewGormLogger() *GormLogger {
	return &GormLogger{
		SlowThreshold: 200 * time.Millisecond, // 一般超过200毫秒就算慢查所以不使用配置进行更改
	}
}

var _ gormLogger.Interface = (*GormLogger)(nil)

func (l *GormLogger) LogMode(lev gormLogger.LogLevel) gormLogger.Interface {
	return &GormLogger{
		//Level: lev,
	}
}
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	InfoF(msg, data)
}
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	WarnF(msg, data)
}
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	ErrorF(msg, data)
}
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		ErrorFTimes(5, "[SQL-Error]\t| err = %v \t| rows= %v \t| %v \t| %v", err, rows, elapsed, sql)
	}
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		WarnFTimes(5, "[SQL-SlowLog]\t| rows= %v \t| %v \t| %v", rows, elapsed, sql)
	}
	InfoFTimes(5, "[SQL]\t| rows= %v \t| %v \t| %v", rows, elapsed, sql)
}
