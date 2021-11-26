package clients

import "os"

type ClientConfig interface {
	Get(name string) string
}

type EnvironmentClientConfig struct{}

func (e *EnvironmentClientConfig) Get(name string) string {
	return os.Getenv(name)
}

type DynamicClientConfig struct {
	fallback       ClientConfig
	configurations map[string]string
}

func DynamicConfig() *DynamicClientConfig {
	return &DynamicClientConfig{}
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
