package com

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"testing"

	"golang.org/x/crypto/ssh"

	"github.com/pkg/sftp"
)

const sshTurstedKey = "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBHHpFhk/28mT/rtxi77vz/NY35GCHoZwi7cAImtQITsy0uQRWH2PmhNwXRzzaZTtDFChzjE9BZwT7CtSGDt549M="

// SSH Key-strings
func trustedHostKeyCallback(trustedKey string) ssh.HostKeyCallback {

	if trustedKey == "" {
		return func(_ string, _ net.Addr, k ssh.PublicKey) error {
			log.Printf("WARNING: SSH-key verification is *NOT* in effect: to fix, add this trustedKey: %q", keyString(k))
			return nil
		}
	}

	return func(_ string, _ net.Addr, k ssh.PublicKey) error {
		ks := keyString(k)
		if trustedKey != ks {
			return fmt.Errorf("SSH-key verification: expected %q but got %q", trustedKey, ks)
		}

		return nil
	}
}

func keyString(k ssh.PublicKey) string {
	return k.Type() + " " + base64.StdEncoding.EncodeToString(k.Marshal())
}

func TestSshSimple(t *testing.T) {

	username := "test"
	password := "123"

	ip := "127.0.0.1"
	port := "22"

	// username = "latelee"
	// password = "123456"
	// ip = "192.168.28.11"

	username = "aftp"
	password = "passworld"

	client := NewSSHClient(username, password, ip, port)

	// 1.运行远程命令
	// cmd := "ls"
	// backinfo, err := client.Run(cmd)
	// if err != nil {
	// 	fmt.Printf("failed to run shell,err=[%v]\n", err)
	// 	return
	// }
	// fmt.Printf("%v back info: \n[%v]\n", cmd, backinfo)

	// 2. 上传一文件
	filename := "foo.txt"
	WriteFile(filename, []byte("hello ssh\r\n"))

	// 可选上传目录或文件
	n, err := client.upload_file(filename, filename)
	if err != nil {
		fmt.Printf("upload failed: %v\n", err)
		return
	}
	// // 3. 显示该文件
	// cmd = "cat " + "/tmp/" + filename
	// backinfo, err = client.Run(cmd)
	// if err != nil {
	// 	fmt.Printf("run cmd faild: %v\n", err)
	// 	return
	// }
	// fmt.Printf("%v back info: \n[%v]\n", cmd, backinfo)

	// 4. 下载该文件到本地
	//n, err := client.DownloadFile("/tmp/"+mypath, "testdata_new")
	n, err = client.download_file(filename, "foo_new.txt")
	if err != nil {
		fmt.Printf("download failed: %v\n", err)
		return
	}
	fmt.Printf("download file[%v] ok, size=[%d]\n", filename, n)
}

func TestSshNew(t *testing.T) {

	username := "test"
	password := "123"

	ip := "127.0.0.1"
	port := "22"

	// username = "latelee"
	// password = "123456"
	// ip = "192.168.28.11"

	username = "aftp"
	password = "passworld"

	client := NewSSHClient(username, password, ip, port)

	// 1.运行远程命令
	cmd := "ls"
	backinfo, err := client.Run(cmd)
	if err != nil {
		fmt.Printf("failed to run shell,err=[%v]\n", err)
		return
	}
	fmt.Printf("%v back info: \n[%v]\n", cmd, backinfo)

	// 2. 上传一文件
	mypath := "testdata"
	filename := "foo.txt"
	WriteFile(path.Join(mypath, filename), []byte("hello ssh\r\n"))

	// 可选上传目录或文件
	n, err := client.UploadFile(mypath, "/tmp/testdata")
	if err != nil {
		fmt.Printf("upload failed: %v\n", err)
		return
	}
	// 3. 显示该文件
	cmd = "cat " + "/tmp/" + filename
	backinfo, err = client.Run(cmd)
	if err != nil {
		fmt.Printf("run cmd faild: %v\n", err)
		return
	}
	fmt.Printf("%v back info: \n[%v]\n", cmd, backinfo)

	// 4. 下载该文件到本地
	n, err = client.DownloadFile("/tmp/"+mypath, "testdata_new")
	if err != nil {
		fmt.Printf("download failed: %v\n", err)
		return
	}
	fmt.Printf("download file[%v] ok, size=[%d]\n", filename, n)
}

// Based on example server code from golang.org/x/crypto/ssh and server_standalone
func TestSftpServer(t *testing.T) {

	user := "aftp"
	passwd := "pass"
	var (
		readOnly    bool = true
		debugStderr bool = false
	)

	// flag.BoolVar(&readOnly, "R", false, "read-only server")
	// flag.BoolVar(&debugStderr, "e", false, "debug to stderr")
	// flag.Parse()

	debugStream := ioutil.Discard
	if debugStderr {
		debugStream = os.Stderr
	}

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			// Should use constant-time compare (or better, salt+hash) in
			// a production setting.
			fmt.Fprintf(debugStream, "Login: %s\n", c.User())
			if c.User() == user && string(pass) == passwd {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}

	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key", err)
	}

	config.AddHostKey(private)

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", "0.0.0.0:2022")
	if err != nil {
		log.Fatal("failed to listen for connection", err)
	}
	fmt.Printf("Listening on %v\n", listener.Addr())

	nConn, err := listener.Accept()
	if err != nil {
		log.Fatal("failed to accept incoming connection", err)
	}

	// Before use, a handshake must be performed on the incoming
	// net.Conn.
	_, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		log.Fatal("failed to handshake", err)
	}
	fmt.Fprintf(debugStream, "SSH server established\n")

	// The incoming Request channel must be serviced.
	go ssh.DiscardRequests(reqs)

	// Service the incoming Channel channel.
	for newChannel := range chans {
		// Channels have a type, depending on the application level
		// protocol intended. In the case of an SFTP session, this is "subsystem"
		// with a payload string of "<length=4>sftp"
		fmt.Fprintf(debugStream, "Incoming channel: %s\n", newChannel.ChannelType())
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			fmt.Fprintf(debugStream, "Unknown channel type: %s\n", newChannel.ChannelType())
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Fatal("could not accept channel.", err)
		}
		fmt.Fprintf(debugStream, "Channel accepted\n")

		// Sessions have out-of-band requests such as "shell",
		// "pty-req" and "env".  Here we handle only the
		// "subsystem" request.
		go func(in <-chan *ssh.Request) {
			for req := range in {
				fmt.Fprintf(debugStream, "Request: %v\n", req.Type)
				ok := false
				switch req.Type {
				case "subsystem":
					fmt.Fprintf(debugStream, "Subsystem: %s\n", req.Payload[4:])
					if string(req.Payload[4:]) == "sftp" {
						ok = true
					}
				}
				fmt.Fprintf(debugStream, " - accepted: %v\n", ok)
				req.Reply(ok, nil)
			}
		}(requests)

		serverOptions := []sftp.ServerOption{
			sftp.WithDebug(debugStream),
		}

		if readOnly {
			serverOptions = append(serverOptions, sftp.ReadOnly())
			fmt.Fprintf(debugStream, "Read-only server\n")
		} else {
			fmt.Fprintf(debugStream, "Read write server\n")
		}

		server, err := sftp.NewServer(
			channel,
			serverOptions...,
		)
		if err != nil {
			log.Fatal(err)
		}
		if err := server.Serve(); err == io.EOF {
			server.Close()
			log.Print("sftp client exited session.")
		} else if err != nil {
			log.Fatal("sftp server completed with error:", err)
		}
	}
}
