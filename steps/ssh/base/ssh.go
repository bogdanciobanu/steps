package base

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"

	envconf "github.com/caarlos0/env/v6"
	sdkExec "github.com/stackpulse/steps-sdk-go/exec"
)

const (
	PrivateKeyPath = "/key"
	LogFilePath    = "/log"
)

type SSHArgs struct {
	Username              string        `env:"USERNAME,required" envDefault:""`
	Hostname              string        `env:"HOSTNAME,required" envDefault:""`
	Command               string        `env:"COMMAND,required" envDefault:""`
	AWSSecretKey          string        `env:"AWS_SECRET_KEY" envDefault:""`
	AWSRegion             string        `env:"AWS_REGION" envDefault:""`
	PrivateKey            string        `env:"PRIVATE_KEY" envDefault:""`
	Password              string        `env:"PASSWORD" envDefault:""`
	StrictHostKeyChecking string        `env:"STRICT_HOST_KEY_CHECKING" envDefault:"no"`
	LogLevel              string        `env:"LOG_LEVEL" envDefault:"ERROR"`
	Port                  int           `env:"PORT" envDefault:"22"`
	ConnectionTimeout     time.Duration `env:"CONNECTION_TIMEOUT" envDefault:"30s"`
}

func ParseArgs() (SSHArgs, error) {
	result := SSHArgs{}
	err := envconf.Parse(&result)
	if err != nil {
		return result, err
	}

	if result.AWSSecretKey == "" && result.PrivateKey == "" && result.Password == "" {
		return result, fmt.Errorf("private key, aws secret, or password is required")
	}
	if result.AWSSecretKey != "" && result.AWSRegion == "" {
		return result, fmt.Errorf("aws region is required when specifing aws secret")
	}

	return result, nil

}

func (s *SSHArgs) BuildCommand() (string, []string, error) {
	if s.Password != "" {
		sshCmd, sshArgs := BuildPasswordCommand(s.Username, s.Hostname, s.Command, s.StrictHostKeyChecking, s.LogLevel, s.Port, s.ConnectionTimeout, s.Password)
		return sshCmd, sshArgs, nil
	}
	privateKey := s.PrivateKey
	if s.AWSSecretKey != "" {
		var err error
		privateKey, err = FetchAwsSecret(s.AWSSecretKey, s.AWSRegion)
		if err != nil {
			return "", nil, fmt.Errorf("fetchAwsSecret: %w", err)
		}
	}

	err := ioutil.WriteFile(PrivateKeyPath, []byte(privateKey), 0644)
	if err != nil {
		return "", nil, fmt.Errorf("write private key: %w", err)
	}

	// Restrict key.pem file capabilities (for ssh usage)
	output, _, err := sdkExec.Execute("chmod", []string{"600", PrivateKeyPath})
	if err != nil {
		return "", nil, fmt.Errorf("chmod private key: %w: %s", err, output)
	}

	sshCmd, sshArgs := BuildPrivateKeyCommand(s.Username, s.Hostname, s.Command, s.StrictHostKeyChecking, s.LogLevel, s.Port, s.ConnectionTimeout)
	return sshCmd, sshArgs, nil
}

func FetchAwsSecret(secretKey, region string) (string, error) {
	secretsService := secretsmanager.New(session.New(), aws.NewConfig().WithRegion(region))
	if secretsService == nil {
		return "", fmt.Errorf("initialize AWS secrets manager")
	}

	result, err := secretsService.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId:     &secretKey,
		VersionStage: aws.String("AWSCURRENT"),
	})
	if err != nil || result.SecretString == nil {
		return "", fmt.Errorf("fetch aws secret: %w", err)
	}

	return *result.SecretString, nil
}

// try to print ssh log, ignore errors.
func PrintLog() {
	content, _ := ioutil.ReadFile(LogFilePath)
	fmt.Println("--SSHLOG--")
	fmt.Printf("%s", string(content))
	fmt.Println("----------")
}

func PrintOutput(output []byte) {
	fmt.Println("--OUTPUT--")
	fmt.Printf("%s", string(output))
	fmt.Println("----------")
}

func BuildPasswordCommand(username, hostname, linuxCmd, StrictHostKeyChecking, LogLevel string, port int, connectionTimeout time.Duration, password string) (string, []string) {
	args := []string{}
	args = append(args, "-p", password)
	args = append(args, "ssh")
	args = append(args, "-o", fmt.Sprintf("StrictHostKeyChecking=%s", StrictHostKeyChecking))
	args = append(args, "-o", fmt.Sprintf("LogLevel=%s", LogLevel))
	args = append(args, fmt.Sprintf("%s@%s", username, hostname))
	args = append(args, "-p", strconv.Itoa(port))
	args = append(args, "-E", LogFilePath)
	args = append(args, fmt.Sprintf("-oConnectTimeout=%d", int(connectionTimeout.Seconds())))
	args = append(args, linuxCmd)

	return "sshpass", args
}

func BuildPrivateKeyCommand(username, hostname, linuxCmd, StrictHostKeyChecking, LogLevel string, port int, connectionTimeout time.Duration) (string, []string) {
	args := []string{"-o", fmt.Sprintf("StrictHostKeyChecking=%s", StrictHostKeyChecking)}
	args = append(args, "-o", fmt.Sprintf("LogLevel=%s", LogLevel))
	args = append(args, "-i", PrivateKeyPath)
	args = append(args, fmt.Sprintf("%s@%s", username, hostname))
	args = append(args, "-p", strconv.Itoa(port))
	args = append(args, "-E", LogFilePath)
	args = append(args, fmt.Sprintf("-oConnectTimeout=%d", int(connectionTimeout.Seconds())))
	args = append(args, linuxCmd)

	return "ssh", args
}
