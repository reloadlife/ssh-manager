package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
)

type SSH struct {
	privateKey []byte
	host       string
	port       string
	username   string
	password   string

	addr    string
	signer  ssh.Signer
	sshConf *ssh.ClientConfig

	conn *ssh.Client
}

func NewSSH(host, port, username, password string, privateKey []byte) *SSH {
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	return &SSH{
		host:       host,
		port:       port,
		addr:       host + ":" + port,
		username:   username,
		password:   password,
		privateKey: privateKey,
		signer:     signer,
	}
}

func (s *SSH) connect() {
	if s.sshConf == nil {
		s.sshConf = &ssh.ClientConfig{
			User:            s.username,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Auth: []ssh.AuthMethod{
				ssh.Password(s.password),
				ssh.PublicKeys(s.signer),
			},
		}
	}

	conn, err := ssh.Dial("tcp", s.addr, s.sshConf)
	if err != nil {
		fmt.Println(err.Error())
	}
	s.conn = conn
}

func (s *SSH) Connect() {
	if s.conn == nil {
		s.connect()
	}
}

func (s *SSH) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return errors.New("connection was not established")
}

func (s *SSH) Run(cmd string) (string, error) {
	if s.conn == nil {
		s.connect()
	}
	session, err := s.conn.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	out, err := session.Output(cmd)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (s *SSH) RunWithPty(cmd string) (string, error) {
	if s.conn == nil {
		s.connect()
	}
	session, err := s.conn.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return "", err
	}

	out, err := session.Output(cmd)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func (s *SSH) RunWithPtyAndPipe(cmd string) (string, error) {
	if s.conn == nil {
		s.connect()
	}
	session, err := s.conn.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return "", err
	}

	pipeReader, pipeWriter := io.Pipe()
	session.Stdout = pipeWriter
	session.Stderr = pipeWriter

	if err := session.Start(cmd); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	go func() {
		io.Copy(&buf, pipeReader)
	}()

	if err := session.Wait(); err != nil {
		return "", err
	}

	return buf.String(), nil
}
