package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mamad.dev/ssh-manager/pkg/ssh"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := run(ctx); err != nil {
			log.Print(err)
		}
		cancel()
	}()

	select {
	case <-sig:
		cancel()
	case <-ctx.Done():
	}
}

func run(ctx context.Context) error {
	host := "127.0.0.1"
	user := "mamad"
	pwd := "kos"
	pKey := []byte("-----BEGIN OPENSSH PRIVATE KEY-----\nb3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW\nQyNTUxOQAAACDIeYwm53rT57UmALA1UCRtymzgoRyuO4xtqD2S1q/LdwAAAJBa8Ac4WvAH\nOAAAAAtzc2gtZWQyNTUxOQAAACDIeYwm53rT57UmALA1UCRtymzgoRyuO4xtqD2S1q/Ldw\nAAAEBxKOawfwEQaIbTZfsvXcxDZpGR2MWiIr0s5hpzt3ZahMh5jCbnetPntSYAsDVQJG3K\nbOChHK47jG2oPZLWr8t3AAAAC21hbWFkQG1hbWFkAQI=\n-----END OPENSSH PRIVATE KEY-----\n")

	sshClient := ssh.NewSSH(host, "22", user, pwd, pKey)
	err := sshClient.Term(ctx)
	if err != nil {
		fmt.Println(err)
	}
	log.Infof("SSH Closed.")
	return err
}
