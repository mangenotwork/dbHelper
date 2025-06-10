package dbHelper

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

var RedisConn map[string]*redis.Client

func initRedisConn() {
	RedisConn = make(map[string]*redis.Client, len(Conf.RedisConf))
	for _, v := range Conf.RedisConf {
		m, err := redisConn(v)
		if err != nil {
			panic(err)
		}
		RedisConn[v.Tag] = m
	}
}

func GetRedisConn(tag string) *redis.Client {
	m, ok := RedisConn[tag]
	if !ok {
		panic("[Redis] 未init")
	}
	return m
}

func redisConn(conf *RedisConf) (*redis.Client, error) {

	options := &redis.Options{
		Addr:            fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password:        conf.Password,
		DB:              conf.DB,
		PoolSize:        conf.PoolSize,
		MinIdleConns:    conf.MinIdleConn,
		ConnMaxIdleTime: time.Duration(conf.ConnMaxIdleTime) * time.Second,
	}

	var (
		sshClient *ssh.Client
		err       error
	)

	if conf.IsSSH {
		sshConfig := &ssh.ClientConfig{
			User: conf.SSHUsername,
			Auth: []ssh.AuthMethod{
				ssh.Password(conf.SSHPassword),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         10 * time.Second,
		}
		sshClient, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", conf.SSHRemoteHost, conf.SSHRemotePort), sshConfig)
		if err != nil {
			ErrorF("Failed to dial SSH server: %v", err)
			return nil, err
		}
		options.Dialer = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return sshClient.Dial("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port))
		}

		// 禁用不适用于 SSH 隧道的超时设置, https://github.com/redis/go-redis/issues/2057
		// 提到 如果使用最新版本。#2176 已修复该问题。解决办法： #2176 (comment)
		// https://github.com/redis/go-redis/pull/2176
		options.ReadTimeout = -2
		options.WriteTimeout = -2
	}

	redisClient := redis.NewClient(options)
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		if sshClient != nil {
			_ = sshClient.Close()
		}
		return nil, err
	}

	Info("[Redis] connection successful:", pong)

	return redisClient, nil
}
