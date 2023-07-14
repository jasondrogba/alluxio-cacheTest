package startTest

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"runtime"
)

// 设置常量
const (
	//清空命令
	testCMD = "sudo su alluxio -c \"cd /opt/alluxio && touch test\""
	freeCMD = "sudo su alluxio -c \"cd /opt/alluxio && ./bin/alluxio fs free /\""
	//停止命令
	stopCMD = "sudo su alluxio -c \"cd /opt/alluxio && ./bin/alluxio-stop.sh all\""
	//格式化命令
	formatCMD = "sudo su alluxio -c \"cd /opt/alluxio && ./bin/alluxio format\""
	//启动命令
	startCMD = "sudo su alluxio -c \"cd /opt/alluxio && ./bin/alluxio-start.sh all\""
	//动态切换cache eviction policy
	cacheCMD = "sudo su alluxio -c \"cd /opt/alluxio && ./bin/alluxio fsadmin updateConf alluxio.worker.block.annotator.dynamic.sort=REPLICA\""
	port     = "22"
)

func Starttest(hostname string, policy string) {
	config := SetupSSH()
	//cmd := fmt.Sprintf("sudo su alluxio -c \"cd /opt/alluxio && ./bin/alluxio fs load /%d.txt --local flag\"")
	//multiSSH(hostname, port, config, freeCMD)

	if policy == "LRU" {
		fmt.Println("cache policy: LRU")
		multiSSH(hostname, port, config, freeCMD)
		multiSSH(hostname, port, config, stopCMD)
		multiSSH(hostname, port, config, formatCMD)
		multiSSH(hostname, port, config, startCMD)
	}
	if policy == "REPLICA" {
		fmt.Println("cache policy: REPLICA")
		//multiSSH(hostname, port, config, freeCMD)
		//multiSSH(hostname, port, config, stopCMD)
		//multiSSH(hostname, port, config, formatCMD)
		//multiSSH(hostname, port, config, startCMD)
		multiSSH(hostname, port, config, cacheCMD)
	}

}

func SetupSSH() *ssh.ClientConfig {
	// Read the private key file for the SSH connection
	PrivateKeyPath := ""
	if runtime.GOOS == "linux" {
		fmt.Println("Detected Linux system")
		PrivateKeyPath = "/home/ec2-user/.ssh/id_rsa"

	} else if runtime.GOOS == "darwin" {
		fmt.Println("Detected macOS system")
		PrivateKeyPath = "/Users/sunbury/.ssh/id_rsa"

	} else {
		fmt.Println("Unknown system")
	}
	privateKeyBytes, err := os.ReadFile(PrivateKeyPath)
	if err != nil {
		fmt.Println("Failed to read private key file:", err)
		os.Exit(1)
	}
	privateKey, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		fmt.Println("Failed to parse private key:", err)
		os.Exit(1)
	}

	// Set up the SSH configuration
	config := &ssh.ClientConfig{
		User: "ec2-user",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privateKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return config
}

func multiSSH(instance string, port string, config *ssh.ClientConfig, cmd string) {
	conn, err := ssh.Dial("tcp", instance+":"+port, config)
	if err != nil {
		fmt.Println("Failed to establish SSH connection:", err)
		os.Exit(1)
	}
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil {
		fmt.Println("Failed to create session:", err)
		os.Exit(1)
	}
	defer session.Close()

	//cmd := fmt.Sprintf("sudo su alluxio -c \"cd /opt/alluxio && ./bin/alluxio fs load /%d.txt --local flag\"")
	output, err := session.Output(cmd)
	if err != nil {
		fmt.Println("Failed to run command:", err)
		os.Exit(1)
	}
	fmt.Print(string(output))

}
