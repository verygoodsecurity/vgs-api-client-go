package clients

import "os"

type ClientConfig interface {
	Get(name string) string
}

type EnvironmentClientConfig struct {
	fallback ClientConfig
}

func (e *EnvironmentClientConfig) WithFallback(fallback ClientConfig) *EnvironmentClientConfig {
	e.fallback = fallback
	return e
}

func (e *EnvironmentClientConfig) Get(name string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return e.fallback.Get(name)
}

type DynamicClientConfig struct {
	fallback       ClientConfig
	configurations map[string]string
}

func DynamicConfig() *DynamicClientConfig {
	return &DynamicClientConfig{
		configurations: make(map[string]string),
	}
}

func EnvironmentConfig() *EnvironmentClientConfig {
	return &EnvironmentClientConfig{}
}

func (d *DynamicClientConfig) Get(name string) string {
	result, ok := d.configurations[name]
	if ok {
		return result
	}
	return d.fallback.Get(name)
}

func (d *DynamicClientConfig) WithFallback(fallback ClientConfig) *DynamicClientConfig {
	d.fallback = fallback
	return d
}

func (d *DynamicClientConfig) AddParameter(name string, value string) *DynamicClientConfig {
	d.configurations[name] = value
	return d
}
