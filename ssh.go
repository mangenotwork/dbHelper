package dbHelper

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type sshConfig struct {
	User       string // ssh 用户
	Password   string // 密码认证
	PrivateKey string // 密钥文件路径
	RemoteHost string // SSH服务器地址
	RemotePort int64  // SSH服务器端口
	LocalHost  string // 本地监听地址
	LocalPort  int64  // 本地监听端口
	TargetHost string // 目标服务器地址
	TargetPort int64  // 目标服务端口
}

func (s *sshConfig) getSshConn() error {
	// 配置SSH客户端
	config := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// 注意：生产环境中应该验证主机密钥
			// 这里为了简化示例，接受所有密钥
			return nil
		},
		Timeout: 10 * time.Second,
	}

	// 设置认证方法
	if s.Password != "" {
		config.Auth = append(config.Auth, ssh.Password(s.Password))
	}

	if s.PrivateKey != "" {
		key, err := os.ReadFile(s.PrivateKey)
		if err != nil {
			ErrorF("[ssh隧道]读取私钥文件失败: %v", err)
			return err
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			ErrorF("[ssh隧道]解析私钥失败: %v", err)
			return err
		}

		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	}

	if len(config.Auth) == 0 {
		return fmt.Errorf("[ssh隧道]必须指定密码或私钥文件进行认证")
	}

	// 连接到SSH服务器
	sshAddr := fmt.Sprintf("%s:%d", s.RemoteHost, s.RemotePort)
	client, err := ssh.Dial("tcp", sshAddr, config)
	if err != nil {
		ErrorF("[ssh隧道]连接到SSH服务器失败: %v", err)
		return err
	}

	defer func() {
		_ = client.Close()
	}()

	// 本地监听
	localAddr := fmt.Sprintf("127.0.0.1:%d", s.LocalPort)
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		ErrorF("[ssh隧道]本地监听失败: %v", err)
		return err
	}
	defer func() {
		_ = listener.Close()
	}()

	InfoF("[ssh隧道]本地监听已启动: %s", localAddr)
	InfoF("[ssh隧道]转发规则: 本地 %s -> 远程 %s:%d", localAddr, s.TargetHost, s.TargetPort)

	// 设置信号处理，优雅关闭
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		<-signalChan
		InfoF("[ssh隧道]接收到关闭信号，正在关闭 ssh 隧道...")
		_ = listener.Close()
		_ = client.Close()
		os.Exit(0)
	}()

	// 接受本地连接并转发
	for {

		localConn, err := listener.Accept()
		if err != nil {

			// 检查是否是因为监听关闭导致的错误
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				ErrorF("[ssh隧道]临时错误: %v", err)
				continue
			}

			ErrorF("[ssh隧道]接受连接失败: %v", err)

			break
		}

		// todo 优化，需要正确的执行关闭释放连接
		defer localConn.Close()

		// 处理每个连接
		go handleConnection(client, localConn, s.TargetHost, s.TargetPort)
	}

	return nil
}

// 处理每个连接的转发
func handleConnection(client *ssh.Client, localConn net.Conn, targetHost string, targetPort int64) {

	// 连接到远程目标
	targetAddr := fmt.Sprintf("%s:%d", targetHost, targetPort)
	remoteConn, err := client.Dial("tcp", targetAddr)
	if err != nil {
		ErrorF("[ssh隧道]连接到远程目标失败: %v", err)
		return
	}
	InfoF("[ssh隧道]新连接已建立: %s <-> %s", localConn.RemoteAddr(), targetAddr)

	// 双向数据转发
	go forward(localConn, remoteConn)
	go forward(remoteConn, localConn)
}

// 数据转发函数
func forward(src, dst net.Conn) {
	defer func() {
		_ = src.Close()
	}()
	defer func() {
		_ = dst.Close()
	}()

	buf := make([]byte, 64*1024)
	_, err := io.CopyBuffer(dst, src, buf)
	if err != nil {
		ErrorF("[ssh隧道]转发数据失败: %v", err)
	}
}

// GetFreePort 获取一个未被使用的端口
func getFreePort() (int64, error) {
	// 创建一个 TCP 监听器，使用 :0 让系统分配一个可用端口
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = listener.Close() // 确保监听器关闭，释放端口
	}()

	// 获取分配的端口
	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return 0, fmt.Errorf("无法获取 TCP 地址")
	}

	return int64(tcpAddr.Port), nil
}
