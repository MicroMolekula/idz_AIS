package service

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

func pingIP(ip string) bool {
	cmd := exec.Command("ping", "-n", "1", "-w", "1", ip)
	err := cmd.Run()
	return err == nil
}

func ScanNetwork(startIP, endIP string) []string {
	var aliveHosts []string

	// Преобразуем начальный и конечный IP-адреса в числа
	start := ipToUint32(startIP)
	end := ipToUint32(endIP)

	if start > end {
		fmt.Println("Ошибка: начальный IP-адрес больше конечного")
		return aliveHosts
	}

	// Сканируем диапазон IP-адресов
	for i := start; i <= end; i++ {
		ip := uint32ToIP(i)
		if pingIP(ip) {
			aliveHosts = append(aliveHosts, ip)
			fmt.Println(ip, "ok")
		} else {
			fmt.Println(ip, "error")
		}
	}

	return aliveHosts
}

// Функция для преобразования IP-адреса в число
func ipToUint32(ip string) uint32 {
	var result uint32
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return 0
	}
	for i, part := range parts {
		value, err := strconv.Atoi(part)
		if err != nil || value < 0 || value > 255 {
			return 0
		}
		result |= uint32(value) << ((3 - i) * 8)
	}
	return result
}

// Функция для преобразования числа в IP-адрес
func uint32ToIP(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		(ip>>24)&0xFF,
		(ip>>16)&0xFF,
		(ip>>8)&0xFF,
		ip&0xFF,
	)
}

func ExecuteSSHCommand(host, user, password, pureCommand string, sudo bool) (string, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	command := pureCommand
	if sudo {
		command = fmt.Sprintf("echo %s | sudo --stdin %s", password, pureCommand)
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
