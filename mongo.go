package dbHelper

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var MongoDBConn map[string]*mongo.Database

func initMongoDBConn() {
	MongoDBConn = make(map[string]*mongo.Database, len(Conf.MongoDBConf))
	for _, v := range Conf.MongoDBConf {
		m, err := mongoDBConn(v)
		if err != nil {
			panic(err)
		}
		MongoDBConn[v.Tag] = m
	}
}

func GetMongoDBConn(tag string) *mongo.Database {
	m, ok := MongoDBConn[tag]
	if !ok {
		panic("[MongoDB] 未init")
	}
	return m
}

func mongoDBConn(conf *MongoDBConf) (*mongo.Database, error) {

	var (
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

	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			InfoF("[MongoDB log] %v", startedEvent.Command.String())
		},
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
			ErrorF("[MongoDB log] %v", failedEvent.Failure)
		},
	}

	if conf.MaxPoolSize == 0 {
		conf.MaxPoolSize = 100
	}

	if conf.MaxConnIdleTime == 0 {
		conf.MaxConnIdleTime = 300
	}

	if conf.ConnectTimeout == 0 {
		conf.ConnectTimeout = 30
	}

	o := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d", conf.User, conf.Password, host, port)).
		SetRetryWrites(true).                                                        // 重试写
		SetRetryReads(true).                                                         // 重试读
		SetMaxConnIdleTime(time.Duration(conf.MaxConnIdleTime) * time.Second).       // 连接池中保持空闲的最长时间
		SetHeartbeatInterval(10 * time.Second).                                      // 心跳时间  默认 10s
		SetServerSelectionTimeout(time.Duration(conf.ConnectTimeout) * time.Second). // 交互超时 默认30s
		SetConnectTimeout(time.Duration(conf.ConnectTimeout) * time.Second).         // 连接超时 默认30s
		SetMaxPoolSize(uint64(conf.MaxPoolSize)).                                    // 连接池最大数
		SetMonitor(monitor)                                                          // 监控日志

	db, err := mongo.Connect(context.TODO(), o)
	if err != nil {
		Error("connect error of database mongo:", err)
		return nil, err
	}

	if err := db.Ping(context.TODO(), nil); err != nil {
		Error("connect error of database mongo:", err)
		return nil, err
	}

	return db.Database(conf.Database), nil
}
