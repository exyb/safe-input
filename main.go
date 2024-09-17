package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	. "github.com/exyb/safe-input/client"
	. "github.com/exyb/safe-input/qrcode"
)

func main() {
	// 解析命令行参数
	username := flag.String("u", "", "Username for authentication")
	isGenerate := flag.Bool("g", false, "generate qr code for authentication")
	isValidate := flag.Bool("c", false, "validate qr code for authentication")
	clientType := flag.String("t", "", "client type")

	flag.Parse()

	if *isGenerate {
		GeneratingQRCode(*username)
		return
	}
	if *username == "" {
		fmt.Println("Error: Username (-u) is required.")
		os.Exit(1)
	}
	if *clientType == "" && (*isGenerate || *isValidate) {
		fmt.Println("Error: ClientType (-t) is required.")
		os.Exit(1)
	}

    // 提示用户输入 OTP 密码
    fmt.Print("Enter OTP password: ")
    reader := bufio.NewReader(os.Stdin)
    otpPassword, _ := reader.ReadString('\n')
    otpPassword = otpPassword[:len(otpPassword)-1] // 去掉换行符

    valid := ValidateOTP(otpPassword, *username)
    if !valid {
        fmt.Println("Error: Invalid OTP.")
        os.Exit(1)
    }
    if *isValidate {
        return
    }


	// 从环境变量中读取密码
	redisHost := GetEnv("REDIS_HOST", "localhost")
	redisPort := GetEnv("REDIS_PORT", "6379")
	redisDb, _ := strconv.Atoi(GetEnv("REDIS_DB", "0"))

	mysqlHost := GetEnv("MYSQL_HOST", "127.0.0.1")
	mysqlPort := GetEnv("MYSQL_PORT", "3306")
	mysqlUser := GetEnv("MYSQL_USER", "demo")
	mysqlSchema := GetEnv("MYSQL_SCHEMA", "demo")

	switch strings.ToLower(*clientType) {
	case "redis":
		password := os.Getenv("REDIS_PASSWORD")
		StartRedisClient(password, redisHost, redisPort, redisDb)
	case "minio":
		StartMinioClient()
	case "mysql":
		password := GetEnv("MYSQL_PASSWORD", "demo")

		if password == "" {
			fmt.Println("Error: MYSQL_PASSWORD environment variable is not set.")
			os.Exit(1)
		}
		StartMySQLClient(password, mysqlUser, mysqlHost, mysqlPort, mysqlSchema)
	case "mysqlroot":
		password := os.Getenv("MYSQL_ROOT_PASSWORD")
		if password == "" {
			fmt.Println("Error: MYSQL_ROOT_PASSWORD environment variable is not set.")
			os.Exit(1)
		}
		StartMySQLClient(password, "root", mysqlHost, mysqlPort, mysqlSchema)
	default:
		fmt.Println("Invalid client name. Please choose from redis, mysql, mysqlroot, minio.")
		os.Exit(1)
	}

}

func GetEnv(name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}
