package cmdb

import (
	"DYCLOUD/model/cmdb"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

func GetDefaultPrivateKeyPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unable to get current user: %v", err)
	}
	return filepath.Join(usr.HomeDir, ".ssh", "id_rsa"), nil
}

func GenerateSSHKey(privateKeyPath string) error {
	// 生成 RSA 私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("unable to generate private key: %v", err)
	}

	// 编码私钥为 PEM 格式
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err := ioutil.WriteFile(privateKeyPath, privateKeyPEM, 0600); err != nil {
		return fmt.Errorf("unable to write private key: %v", err)
	}

	// 生成公钥
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("unable to generate public key: %v", err)
	}

	// 编码公钥为 OpenSSH 格式
	publicKeyPath := privateKeyPath + ".pub"
	publicKeyBytes := ssh.MarshalAuthorizedKey(publicKey)
	if err := ioutil.WriteFile(publicKeyPath, publicKeyBytes, 0644); err != nil {
		return fmt.Errorf("unable to write public key: %v", err)
	}

	fmt.Printf("SSH key generated successfully.\nPrivate key: %s\nPublic key: %s\n", privateKeyPath, publicKeyPath)
	return nil
}

func CanSSHWithoutPassword(host, port, username, privateKeyPath string) (bool, error) {
	key, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return false, fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return false, fmt.Errorf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	address := fmt.Sprintf("%s:%s", host, port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return false, err
	}
	defer client.Close()

	return true, nil
}

// createSSHClient 创建 SSH 客户端
func CreateSSHClient(host, port, username, privateKeyPath string) (*ssh.Client, error) {
	key, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("读取私钥文件失败: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("解析私钥文件失败: %v", err)
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	}

	address := fmt.Sprintf("%s:%s", host, port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, fmt.Errorf("SSH 连接失败: %v", err)
	}

	return client, nil
}

// validateSSHConnection 验证 SSH 连接
func ValidateSSHConnection(host, port, username, password string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 你可以根据需要调整 HostKeyCallback
		Timeout:         time.Second * 10,            // 设置超时时间
	}

	address := fmt.Sprintf("%s:%s", host, port)
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// populateHostInfo 从 SSH 连接中获取主机信息并填充到 CmdbHosts 结构体中
func PopulateHostInfo(client *ssh.Client, host *cmdb.CmdbHosts) error {
	var err error

	// 获取内存信息
	host.Memory, err = executeCommand(client, `free -h | awk '/^Mem:/ {print $2}' 2>/dev/null`)
	if err != nil {
		return fmt.Errorf("获取内存信息失败: %v", err)
	}

	// 获取磁盘容量
	host.DiskTotal, err = executeCommand(client, `df -h | awk '$NF=="/"{print $2}' 2>/dev/null`)
	if err != nil {
		return fmt.Errorf("获取磁盘容量失败: %v", err)
	}

	// 获取 CPU 信息
	host.CpuModel, err = executeCommand(client, `lscpu | grep 'Model name:' | awk -F ': ' '{print $2}' 2>/dev/null`)
	if err != nil {
		return fmt.Errorf("获取 CPU 信息失败: %v", err)
	}

	// 获取 CPU 数量
	host.CpuCount, err = executeCommand(client, "nproc 2>/dev/null")
	if err != nil {
		return fmt.Errorf("获取 CPU 数量失败: %v", err)
	}

	// 获取操作系统信息
	host.Os, err = executeCommand(client, "cat /etc/*release | grep -E 'PRETTY_NAME' | cut -d= -f2 || echo 'Not available' 2>/dev/null")
	if err != nil {
		return fmt.Errorf("获取操作系统信息失败: %v", err)
	}

	// 获取操作系统版本
	host.OsVersion, err = executeCommand(client, "uname -r 2>/dev/null")
	if err != nil {
		return fmt.Errorf("获取操作系统版本失败: %v", err)
	}

	// 获取公共 IP 地址
	host.PublicIP, err = executeCommand(client, "curl -s ifconfig.me 2>/dev/null")
	if err != nil {
		return fmt.Errorf("获取公共 IP 地址失败: %v", err)
	}

	// 获取私有 IP 地址
	host.PrivateIP, err = executeCommand(client, `hostname -I | awk '{print $1}' 2>/dev/null`)
	if err != nil {
		return fmt.Errorf("获取私有 IP 地址失败: %v", err)
	}

	// 获取磁盘信息并填充到字符串
	diskInfoOut, err := executeCommand(client, "lsblk -o NAME,SIZE | sed 's/├─//g; s/└─//g; s/ //g' | tail -n +2 2>/dev/null")
	if err != nil {
		return fmt.Errorf("获取磁盘信息失败: %v", err)
	}

	// 将磁盘信息处理为 "磁盘名:大小" 格式
	var diskInfoList []string
	lines := strings.Split(strings.TrimSpace(diskInfoOut), "\n")
	for _, line := range lines {
		parts := strings.Fields(line) // 按空格拆分
		if len(parts) >= 2 {
			diskName := parts[0]
			diskSize := parts[1]
			diskInfoList = append(diskInfoList, fmt.Sprintf("%s:%s", diskName, diskSize)) // 生成 "磁盘名:大小" 格式
		}
	}
	host.DiskInfo = strings.Join(diskInfoList, ",")
	fmt.Println(host.DiskInfo)
	return nil
}

// executeCommand
//
//	@Description:  通过 SSH 执行命令并返回输出
//	@param client
//	@param cmd
//	@return string
//	@return error
func executeCommand(client *ssh.Client, cmd string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// EnablePasswordlessSSH
//
//	@Description: 设置本机与远程主机的免密登录
//	@param client
//	@param username
//	@return error
func EnablePasswordlessSSH(client *ssh.Client, username string) error {
	// 读取本机的 SSH 公钥
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户主目录失败: %v", err)
	}

	pubKeyPath := homeDir + "/.ssh/id_rsa.pub"
	pubKey, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		// 如果没有找到公钥，尝试生成一个新的公私钥对
		if os.IsNotExist(err) {
			if err := generateSSHKeyPair(homeDir + "/.ssh/id_rsa"); err != nil {
				return fmt.Errorf("生成 SSH 密钥对失败: %v", err)
			}
			pubKey, err = ioutil.ReadFile(pubKeyPath)
			if err != nil {
				return fmt.Errorf("读取生成的 SSH 公钥失败: %v", err)
			}
		} else {
			return fmt.Errorf("读取 SSH 公钥失败: %v", err)
		}
	}

	// 将公钥添加到远程主机的 authorized_keys 文件中
	cmd := fmt.Sprintf("echo '%s' >> ~/.ssh/authorized_keys", strings.TrimSpace(string(pubKey)))
	if _, err := executeCommand(client, cmd); err != nil {
		return fmt.Errorf("将公钥添加到远程主机失败: %v", err)
	}

	return nil
}

// generateSSHKeyPair
//
//	@Description: 生成新的 SSH 公私钥对
//	@param privateKeyPath
//	@return error
func generateSSHKeyPair(privateKeyPath string) error {
	if err := os.MkdirAll(filepath.Dir(privateKeyPath), 0700); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	privateKeyCmd := fmt.Sprintf("ssh-keygen -t rsa -b 2048 -N '' -f %s", privateKeyPath)
	cmd := exec.Command("bash", "-c", privateKeyCmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("生成 SSH 密钥对失败: %v", err)
	}

	return nil
}
