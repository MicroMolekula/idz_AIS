package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os/exec"
)

func main() {
	result, err := executeSSHCommand("", "", "", "echo  | sudo --stdin shutdown -h now")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	//fmt.Println(scanNetwork("192.168.0"))
}

func pingIP(ip string) bool {
	cmd := exec.Command("ping", "-c", "1", "-w", "1", ip)
	err := cmd.Run()
	return err == nil
}

func scanNetwork(subnet string) []string {
	var aliveHosts []string

	for i := 1; i <= 254; i++ {
		ip := fmt.Sprintf("%s.%d", subnet, i)
		if pingIP(ip) {
			aliveHosts = append(aliveHosts, ip)
			fmt.Println(ip, "ok")
		} else {
			fmt.Println(ip, "error")
		}
	}

	return aliveHosts
}

func executeSSHCommand(host, user, password, command string) (string, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return "", err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.Output(command)
	if err != nil {
		return "", err
	}

	return string(output), nil
}
