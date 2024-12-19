package pc

import (
	"idz_ais/internal/service"
)

type PC struct {
	IpHost   string
	UserName string
	Password string
}

// "echo  | sudo --stdin shutdown -h now"
func NewPC(ipHost, userName, password string) *PC {
	return &PC{
		IpHost:   ipHost,
		UserName: userName,
		Password: password,
	}
}

func (p *PC) GetUserInfo() string {
	command := "whoami"
	result, err := service.ExecuteSSHCommand(p.IpHost, p.UserName, p.Password, command, false)
	if err != nil {
		return "Ошибка получение информации о пользователе"
	}
	return "Текущий пользователь: " + result
}

func (p *PC) GetProcessInfo() string {
	command := "ps aux"
	result, err := service.ExecuteSSHCommand(p.IpHost, p.UserName, p.Password, command, false)
	if err != nil {
		return "Ошибка получение процессов"
	}
	return result
}

func (p *PC) ShutdownHost() string {
	command := "shutdown -h now"
	_, err := service.ExecuteSSHCommand(p.IpHost, p.UserName, p.Password, command, true)
	if err != nil {
		return "Ошибка выключения пк"
	}
	return "ПК выключается"
}

func (p *PC) RebootHost() string {
	command := "reboot"
	_, err := service.ExecuteSSHCommand(p.IpHost, p.UserName, p.Password, command, true)
	if err != nil {
		return "Ошибка перезагрузки"
	}
	return "ПК перезагружается"
}
