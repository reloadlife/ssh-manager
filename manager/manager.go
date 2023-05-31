package manager

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	ssh_manager "go.mamad.dev/ssh-manager/internal/ssh-manager"
	"go.mamad.dev/ssh-manager/pkg/ssh"
)

func askForNewFile(path string, ctx context.Context) error {
	fmt.Print("Do you want to create a new SSH Profile ? (Y/n) ")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		if err.Error() == "unexpected newline" {
			fmt.Print("Please respond with Y or N\n")
			return askForNewFile(path, ctx)
		}
		return err
	}

	if response == "Y" || response == "y" {
		return askForSSHInformation(path, ctx)
	}

	fmt.Println("Okay, Closing Application.")
	log.Infof("SSH Manager Closed.")
	ctx.Done()
	return nil
}

func Application(ctx context.Context) error {
	path := "./files"
	files := ssh_manager.NewManager(path)
	filesList := files.GetFiles()
	if len(filesList) == 0 {
		fmt.Println("No SSH Files found in", path)
		return askForNewFile(path, ctx)
	}

	log.Infof("SSH Files found in %s", path)
	return nil
}

func askForSSHInformation(path string, ctx context.Context) error {
	fmt.Print("Enter SSH Host: ")
	var host string
	_, err := fmt.Scanln(&host)
	if err != nil {
		return err
	}

	fmt.Print("Enter SSH Port: ")
	var port string
	_, err = fmt.Scanln(&port)
	if err != nil {
		return err
	}

	fmt.Print("Enter SSH User: ")
	var user string
	_, err = fmt.Scanln(&user)
	if err != nil {
		return err
	}

	fmt.Print("Enter SSH Password: ")
	var pwd string
	_, err = fmt.Scanln(&pwd)
	if err != nil {
		return err
	}

	fmt.Print("Enter SSH Private Key: ")
	var pKey string
	_, err = fmt.Scan(&pKey)
	if err != nil {
		return err
	}

	return nil
}

func sshRun(ctx context.Context) error {
	host := "127.0.0.1"
	user := "mamad"
	pwd := "nopassword"
	pKey := []byte("-----BEGIN OPENSSH PRIVATE KEY-----\nb3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW\nQyNTUxOQAAACDIeYwm53rT57UmALA1UCRtymzgoRyuO4xtqD2S1q/LdwAAAJBa8Ac4WvAH\nOAAAAAtzc2gtZWQyNTUxOQAAACDIeYwm53rT57UmALA1UCRtymzgoRyuO4xtqD2S1q/Ldw\nAAAEBxKOawfwEQaIbTZfsvXcxDZpGR2MWiIr0s5hpzt3ZahMh5jCbnetPntSYAsDVQJG3K\nbOChHK47jG2oPZLWr8t3AAAAC21hbWFkQG1hbWFkAQI=\n-----END OPENSSH PRIVATE KEY-----\n")

	sshClient := ssh.NewSSH(host, "22", user, pwd, pKey)
	err := sshClient.Term(ctx)
	if err != nil {
		fmt.Println(err)
	}
	log.Infof("SSH Closed.")
	return err
}
