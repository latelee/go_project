/*
ssh封装
说明：
1、运行远程脚本或命令，并返回运行结果（如可获取ls的结果）
2、支持上传目录或文件。如果目录不存在，则创建之。
3、支持下载目录或文件。机制同2。
4、无论上传或下载，目标都是目录，暂不实现文件到文件的传输。
  如：upload("a.txt", "tmp/") 而不是：upload("a.txt", "tmp/a.txt")
5、HostKey测试失败，暂不使用，即只使用账号密码登陆。

参考：
https://lifelmy.github.io/post/2022_04_11_go_ssh_scp/
https://github.com/golang/crypto/blob/master/ssh/example_test.go


*/

package com

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Cli struct {
	user       string
	pwd        string // TODO 目前是明文存储，如何加密？
	ip         string
	port       string
	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

func isFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

func NewSSHClient(user, pwd, ip, port string) Cli {
	return Cli{
		user: user,
		pwd:  pwd,
		ip:   ip,
		port: port,
	}
}

func (c *Cli) getConfig1() *ssh.ClientConfig {
	// Every client must provide a host key check.  Here is a
	// simple-minded parse of OpenSSH's known_hosts file
	host := c.port
	// 带环境变量的，需使用Getenv获取
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	if hostKey == nil {
		log.Fatalf("no hostkey for %s", host)
	}

	config := &ssh.ClientConfig{
		User: c.user, // os.Getenv("USER"),
		Auth: []ssh.AuthMethod{
			ssh.Password(c.pwd),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	return config
}

func (c *Cli) getConfig2() *ssh.ClientConfig {
	var hostKey ssh.PublicKey
	key, err := ioutil.ReadFile("/c/Users/Administrator/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: c.user,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	return config
}

// 不使用 HostKey， 使用密码
func (c *Cli) getConfig_nokey() *ssh.ClientConfig {
	config := &ssh.ClientConfig{
		User: c.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.pwd),
		},
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return config
}

func (c *Cli) Connect() error {
	config := c.getConfig_nokey() //c.getConfig1()
	client, err := ssh.Dial("tcp", c.ip+":"+c.port, config)
	if err != nil {
		return fmt.Errorf("connect server error: %w", err)
	}
	sftp, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("new sftp client error: %w", err)
	}

	c.sshClient = client
	c.sftpClient = sftp
	return nil
}

func (c Cli) Run(cmd string) (string, error) {
	if c.sshClient == nil {
		if err := c.Connect(); err != nil {
			return "", err
		}
	}

	session, err := c.sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("create new session error: %w", err)
	}
	defer session.Close()

	buf, err := session.CombinedOutput(cmd)
	return string(buf), err
}

func (c Cli) download_file(remoteFile, localPath string) (int, error) {
	source, err := c.sftpClient.Open(remoteFile)
	if err != nil {
		return -1, fmt.Errorf("sftp client open file error: %w", err)
	}
	defer source.Close()

	localFile := path.Join(localPath, path.Base(remoteFile))
	os.MkdirAll(localPath, os.ModePerm)
	target, err := os.OpenFile(localFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return -1, fmt.Errorf("open local file error: %w", err)
	}
	defer target.Close()

	n, err := io.Copy(target, source)
	if err != nil {
		return -1, fmt.Errorf("write file error: %w", err)
	}
	return int(n), nil
}

func (c Cli) DownloadFile(remotePath, localPath string) (int, error) {
	if c.sshClient == nil {
		if err := c.Connect(); err != nil {
			return -1, err
		}
	}

	// 如是文件直接下载
	if isFile(remotePath) {
		return c.download_file(remotePath, localPath)
	}

	// 如是目录，递归下载
	remoteFiles, err := c.sftpClient.ReadDir(remotePath)
	if err != nil {
		return -1, fmt.Errorf("read path failed: %w", err)
	}
	for _, item := range remoteFiles {
		remoteFilePath := path.Join(remotePath, item.Name())
		localFilePath := path.Join(localPath, item.Name())
		if item.IsDir() {
			err = os.MkdirAll(localFilePath, os.ModePerm)
			if err != nil {
				return -1, err
			}
			_, err = c.DownloadFile(remoteFilePath, localFilePath) // 递归本函数
			if err != nil {
				return -1, err
			}
		} else {
			_, err = c.download_file(path.Join(remotePath, item.Name()), localPath)
			if err != nil {
				return -1, err
			}
		}
	}

	return 0, nil
}

func (c Cli) upload_file(localFile, remotePath string) (int, error) {
	file, err := os.Open(localFile)
	if nil != err {
		return -1, fmt.Errorf("open local file failed: %w", err)
	}
	defer file.Close()

	remoteFileName := path.Base(localFile)
	c.sftpClient.MkdirAll(remotePath)
	ftpFile, err := c.sftpClient.Create(path.Join(remotePath, remoteFileName))
	if nil != err {
		return -1, fmt.Errorf("Create remote path failed: %w", err)
	}
	defer ftpFile.Close()

	fileByte, err := ioutil.ReadAll(file)
	if nil != err {
		return -1, fmt.Errorf("read local file failed: %w", err)
	}

	ftpFile.Write(fileByte)

	return 0, nil
}

func (c Cli) UploadFile(localPath, remotePath string) (int, error) {
	if c.sshClient == nil {
		if err := c.Connect(); err != nil {
			return -1, err
		}
	}

	// 如是文件直接上传
	if isFile(localPath) {
		return c.upload_file(localPath, remotePath)
	}

	// 如是目录，递归上传
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		return -1, fmt.Errorf("read path failed: %w", err)
	}
	for _, item := range localFiles {
		localFilePath := path.Join(localPath, item.Name())
		remoteFilePath := path.Join(remotePath, item.Name())
		if item.IsDir() {
			err = c.sftpClient.Mkdir(remoteFilePath)
			if err != nil {
				return -1, err
			}
			_, err = c.UploadFile(localFilePath, remoteFilePath) // 递归本函数
			if err != nil {
				return -1, err
			}
		} else {
			_, err = c.upload_file(path.Join(localPath, item.Name()), remotePath)
			if err != nil {
				return -1, err
			}
		}
	}

	return 0, nil
}
