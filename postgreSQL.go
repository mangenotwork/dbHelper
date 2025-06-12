package dbHelper

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

var PgsqlConn map[string]*sql.DB

func GetPgsqlConn(tag string) *sql.DB {
	m, ok := PgsqlConn[tag]
	if !ok {
		panic("[PostgreSQL] 未init")
	}
	return m
}

func initPgsqlConn() {
	PgsqlConn = make(map[string]*sql.DB, len(Conf.PgsqlConf))
	for _, v := range Conf.PgsqlConf {
		m, err := pgsqlConn(v)
		if err != nil {
			panic(err)
		}
		PgsqlConn[v.Tag] = m
	}
}

func pgsqlConn(conf *PgsqlConf) (*sql.DB, error) {

	var (
		db   *sql.DB
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
			TargetPort: int64(conf.Port),
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

	// 构建连接字符串
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, conf.User, conf.Password, conf.Database)

	// 打开数据库连接
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		Error(err)
	}
	defer func() {
		_ = db.Close()
	}()

	// 测试连接
	err = db.Ping()
	if err != nil {
		Error(err)
		return nil, err
	}

	// 设置连接池参数
	if conf.MaxOpen == 0 {
		conf.MaxOpen = 20
	}
	if conf.MaxIdle == 0 {
		conf.MaxIdle = 20
	}
	if conf.MaxLifeTime == 0 {
		conf.MaxLifeTime = 60 * 1000 //60s
	}
	db.SetMaxOpenConns(int(conf.MaxOpen))
	db.SetMaxIdleConns(int(conf.MaxIdle))
	db.SetConnMaxLifetime(time.Duration(conf.MaxLifeTime) * time.Millisecond)

	return db, nil
}
