package cache

import (
	"fmt"
	"os"
	"time"
)

type IConfig interface {
	CacherConfig() ICacherConfig
}

type Config struct{}

func NewConfig() IConfig {
	return &Config{}
}

func (cfg *Config) CacherConfig() ICacherConfig {
	return NewCacherConfig()
}

type CacherConfig struct{}

func NewCacherConfig() *CacherConfig {
	return &CacherConfig{}
}

func (cfg *CacherConfig) Endpoint() string {
	redis := os.Getenv("REDIS_HOST")
	return fmt.Sprintf("%s:6379", redis)
}

func (cfg *CacherConfig) Password() string {
	return ""
}

func (cfg *CacherConfig) DB() int {
	return 0
}

func (cfg *CacherConfig) ConnectionSettings() ICacherConnectionSettings {
	return NewDefaultCacherConnectionSettings()
}

type PersisterConfig struct{}

func NewPersisterConfig() *PersisterConfig {
	return &PersisterConfig{}
}

func (cfg *PersisterConfig) Endpoint() string {
	return "127.0.0.1"
}

func (cfg *PersisterConfig) Port() string {
	return "3306"
}

func (cfg *PersisterConfig) DB() string {
	return "my_database"
}

func (cfg *PersisterConfig) Username() string {
	return "my_user"
}

func (cfg *PersisterConfig) Password() string {
	return "my_password"
}

func (cfg *PersisterConfig) Charset() string {
	return "utf8mb4"
}

// ICacherConfig is cacher configuration interface
type ICacherConfig interface {
	Endpoint() string
	Password() string
	DB() int
	ConnectionSettings() ICacherConnectionSettings
}

// ICacherConnectionSettings is connection settings for cacher
type ICacherConnectionSettings interface {
	PoolSize() int
	MinIdleConns() int
	MaxRetries() int
	MinRetryBackoff() time.Duration
	MaxRetryBackoff() time.Duration
	IdleTimeout() time.Duration
	IdleCheckFrequency() time.Duration
	PoolTimeout() time.Duration
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
}

// DefaultCacherConnectionSettings contains default connection settings, this intend to use as embed struct
type DefaultCacherConnectionSettings struct{}

func NewDefaultCacherConnectionSettings() ICacherConnectionSettings {
	return &DefaultCacherConnectionSettings{}
}

func (setting *DefaultCacherConnectionSettings) PoolSize() int {
	return 50
}

func (setting *DefaultCacherConnectionSettings) MinIdleConns() int {
	return 5
}

func (setting *DefaultCacherConnectionSettings) MaxRetries() int {
	return 3
}

func (setting *DefaultCacherConnectionSettings) MinRetryBackoff() time.Duration {
	return 10 * time.Millisecond
}

func (setting *DefaultCacherConnectionSettings) MaxRetryBackoff() time.Duration {
	return 500 * time.Millisecond
}

func (setting *DefaultCacherConnectionSettings) IdleTimeout() time.Duration {
	return 30 * time.Minute
}

func (setting *DefaultCacherConnectionSettings) IdleCheckFrequency() time.Duration {
	return time.Minute
}

func (setting *DefaultCacherConnectionSettings) PoolTimeout() time.Duration {
	return time.Minute
}

func (setting *DefaultCacherConnectionSettings) ReadTimeout() time.Duration {
	return time.Minute
}

func (setting *DefaultCacherConnectionSettings) WriteTimeout() time.Duration {
	return time.Minute
}
