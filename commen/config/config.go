package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	ConfOnce   sync.Once
	instance   *Config
	configPath = "commen/config/config.yaml" // 修改为相对路径，更符合项目结构
)

// Config 代表应用程序的配置结构体
type Config struct {
	Database   DatabaseConfig   `yaml:"database"`
	Redis      RedisConfig      `yaml:"redis"`
	RabbitMQ   RabbitMQConfig   `yaml:"rabbitmq"`
	Prometheus PrometheusConfig `yaml:"prometheus"`
	Logstash   LogstashConfig   `yaml:"logstash"`
	JWT        JWTConfig        `yaml:"jwt"`
	Server     ServerConfig     `yaml:"server"`
}

// DatabaseConfig 代表数据库配置
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

// RedisConfig 代表Redis配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// RabbitMQConfig 代表RabbitMQ配置
type RabbitMQConfig struct {
	Addr string `yaml:"addr"`
}

// PrometheusConfig 代表Prometheus配置
type PrometheusConfig struct {
	Port int `yaml:"port"`
}

// LogstashConfig 代表Logstash配置
type LogstashConfig struct {
	Addr string `yaml:"addr"`
}

// JWTConfig 代表JWT配置
type JWTConfig struct {
	SecretKey      string `yaml:"secret_key"`
	ExpirationTime string `yaml:"expiration_time"`
}

// ServerConfig 代表服务器配置
type ServerConfig struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"` // debug, release, test
}

// LoadConfig 从YAML文件中加载配置
func loadConfig() (*Config, error) {
	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist at path: %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &config, nil
}

// GetConfig 返回配置的单例实例
func GetConfig() *Config {
	ConfOnce.Do(func() {
		conf, err := loadConfig()
		if err != nil {
			log.Fatalf("error loading config: %v", err)
		}
		instance = conf
	})
	return instance
}

// GetDBConnectionString 获取数据库连接字符串
func (c *Config) GetDBConnectionString() string {
	db := c.Database
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", db.User, db.Password, db.Host, db.Port, db.Name)
}

// GetServerAddress 获取服务器地址
func (c *Config) GetServerAddress() string {
	return ":" + c.Server.Port
}

// Example usage
// func main() {
// 	config := GetConfig()
//
// 	fmt.Printf("Database Host: %s\n", config.Database.Host)
// 	fmt.Printf("Server Port: %s\n", config.Server.Port)
// 	fmt.Printf("JWT Secret Key: %s\n", config.JWT.SecretKey)
// }
