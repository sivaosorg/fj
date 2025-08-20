package config

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/sivaosorg/fj/pkg/errors"
)

// Config represents the configuration structure for the FJ library
type Config struct {
	// EnableTransformers controls whether transformations are enabled
	EnableTransformers bool `json:"enable_transformers"`

	// EnableEventPublishing controls whether events are published
	EnableEventPublishing bool `json:"enable_event_publishing"`

	// MaxParseDepth limits the maximum depth for JSON parsing
	MaxParseDepth int `json:"max_parse_depth"`

	// DefaultTimeout sets the default timeout for operations
	DefaultTimeout time.Duration `json:"default_timeout"`

	// CacheSize sets the size of internal caches
	CacheSize int `json:"cache_size"`

	// EnableValidation controls whether input validation is performed
	EnableValidation bool `json:"enable_validation"`

	// LogLevel sets the logging level (0=none, 1=error, 2=warn, 3=info, 4=debug)
	LogLevel int `json:"log_level"`

	// CustomTransformers allows registration of custom transformers at startup
	CustomTransformers map[string]string `json:"custom_transformers"`
}

// GetEnableTransformers implements the Config interface
func (c *Config) GetEnableTransformers() bool { return c.EnableTransformers }

// GetEnableEventPublishing implements the Config interface  
func (c *Config) GetEnableEventPublishing() bool { return c.EnableEventPublishing }

// GetMaxParseDepth implements the Config interface
func (c *Config) GetMaxParseDepth() int { return c.MaxParseDepth }

// GetDefaultTimeout implements the Config interface
func (c *Config) GetDefaultTimeout() time.Duration { return c.DefaultTimeout }

// GetCacheSize implements the Config interface
func (c *Config) GetCacheSize() int { return c.CacheSize }

// GetEnableValidation implements the Config interface
func (c *Config) GetEnableValidation() bool { return c.EnableValidation }

// GetLogLevel implements the Config interface
func (c *Config) GetLogLevel() int { return c.LogLevel }

// GetCustomTransformers implements the Config interface
func (c *Config) GetCustomTransformers() map[string]string { return c.CustomTransformers }

// DefaultConfig returns the default configuration for the FJ library
func DefaultConfig() *Config {
	return &Config{
		EnableTransformers:    true,
		EnableEventPublishing: false,
		MaxParseDepth:        64,
		DefaultTimeout:       30 * time.Second,
		CacheSize:           1000,
		EnableValidation:     true,
		LogLevel:            2, // Warn level
		CustomTransformers:   make(map[string]string),
	}
}

// Manager implements the ConfigurationManager interface using the Singleton pattern.
// It provides thread-safe access to configuration settings and supports
// loading/saving configuration from/to files.
type Manager struct {
	config *Config
	mutex  sync.RWMutex
	settings map[string]interface{}
}

// singleton instance
var (
	instance *Manager
	once     sync.Once
)

// GetInstance returns the singleton instance of the configuration manager
func GetInstance() *Manager {
	once.Do(func() {
		instance = &Manager{
			config:   DefaultConfig(),
			settings: make(map[string]interface{}),
		}
	})
	return instance
}

// GetConfig returns the current configuration
func (m *Manager) GetConfig() *Config {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	// Return a copy to prevent external modification
	configCopy := *m.config
	return &configCopy
}

// SetConfig updates the configuration
func (m *Manager) SetConfig(config *Config) error {
	if config == nil {
		return errors.NewConfigurationError("configuration cannot be nil")
	}

	// Validate configuration
	if err := m.validateConfig(config); err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.config = &(*config) // Create a copy
	return nil
}

// GetSetting gets a specific configuration setting
func (m *Manager) GetSetting(key string) (interface{}, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	// First check custom settings
	if value, exists := m.settings[key]; exists {
		return value, true
	}
	
	// Then check built-in config fields
	switch key {
	case "enable_transformers":
		return m.config.EnableTransformers, true
	case "enable_event_publishing":
		return m.config.EnableEventPublishing, true
	case "max_parse_depth":
		return m.config.MaxParseDepth, true
	case "default_timeout":
		return m.config.DefaultTimeout, true
	case "cache_size":
		return m.config.CacheSize, true
	case "enable_validation":
		return m.config.EnableValidation, true
	case "log_level":
		return m.config.LogLevel, true
	case "custom_transformers":
		return m.config.CustomTransformers, true
	default:
		return nil, false
	}
}

// SetSetting sets a specific configuration setting
func (m *Manager) SetSetting(key string, value interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	// Handle built-in config fields
	switch key {
	case "enable_transformers":
		if v, ok := value.(bool); ok {
			m.config.EnableTransformers = v
		} else {
			return errors.NewConfigurationError("enable_transformers must be a boolean")
		}
	case "enable_event_publishing":
		if v, ok := value.(bool); ok {
			m.config.EnableEventPublishing = v
		} else {
			return errors.NewConfigurationError("enable_event_publishing must be a boolean")
		}
	case "max_parse_depth":
		if v, ok := value.(int); ok {
			if v <= 0 {
				return errors.NewConfigurationError("max_parse_depth must be positive")
			}
			m.config.MaxParseDepth = v
		} else {
			return errors.NewConfigurationError("max_parse_depth must be an integer")
		}
	case "default_timeout":
		if v, ok := value.(time.Duration); ok {
			if v <= 0 {
				return errors.NewConfigurationError("default_timeout must be positive")
			}
			m.config.DefaultTimeout = v
		} else {
			return errors.NewConfigurationError("default_timeout must be a time.Duration")
		}
	case "cache_size":
		if v, ok := value.(int); ok {
			if v < 0 {
				return errors.NewConfigurationError("cache_size cannot be negative")
			}
			m.config.CacheSize = v
		} else {
			return errors.NewConfigurationError("cache_size must be an integer")
		}
	case "enable_validation":
		if v, ok := value.(bool); ok {
			m.config.EnableValidation = v
		} else {
			return errors.NewConfigurationError("enable_validation must be a boolean")
		}
	case "log_level":
		if v, ok := value.(int); ok {
			if v < 0 || v > 4 {
				return errors.NewConfigurationError("log_level must be between 0 and 4")
			}
			m.config.LogLevel = v
		} else {
			return errors.NewConfigurationError("log_level must be an integer")
		}
	default:
		// Store in custom settings
		m.settings[key] = value
	}
	
	return nil
}

// LoadFromFile loads configuration from a JSON file
func (m *Manager) LoadFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.NewConfigurationError("configuration file not found: " + filepath).
				WithCause(err)
		}
		return errors.NewIOError("failed to read configuration file: " + filepath, err)
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return errors.NewConfigurationError("invalid configuration file format").
			WithCause(err).
			WithContext("file", filepath)
	}
	
	return m.SetConfig(&config)
}

// SaveToFile saves configuration to a JSON file
func (m *Manager) SaveToFile(filepath string) error {
	config := m.GetConfig()
	
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return errors.NewConfigurationError("failed to marshal configuration").
			WithCause(err)
	}
	
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return errors.NewIOError("failed to write configuration file: " + filepath, err)
	}
	
	return nil
}

// Reset resets configuration to defaults
func (m *Manager) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.config = DefaultConfig()
	m.settings = make(map[string]interface{})
}

// validateConfig validates the configuration settings
func (m *Manager) validateConfig(config *Config) error {
	if config.MaxParseDepth <= 0 {
		return errors.NewConfigurationError("max_parse_depth must be positive")
	}
	
	if config.DefaultTimeout <= 0 {
		return errors.NewConfigurationError("default_timeout must be positive")
	}
	
	if config.CacheSize < 0 {
		return errors.NewConfigurationError("cache_size cannot be negative")
	}
	
	if config.LogLevel < 0 || config.LogLevel > 4 {
		return errors.NewConfigurationError("log_level must be between 0 and 4")
	}
	
	return nil
}

// IsTransformersEnabled is a convenience method to check if transformers are enabled
func IsTransformersEnabled() bool {
	manager := GetInstance()
	config := manager.GetConfig()
	return config.EnableTransformers
}

// IsEventPublishingEnabled is a convenience method to check if event publishing is enabled
func IsEventPublishingEnabled() bool {
	manager := GetInstance()
	config := manager.GetConfig()
	return config.EnableEventPublishing
}

// GetMaxParseDepth is a convenience method to get the maximum parse depth
func GetMaxParseDepth() int {
	manager := GetInstance()
	config := manager.GetConfig()
	return config.MaxParseDepth
}

// GetDefaultTimeout is a convenience method to get the default timeout
func GetDefaultTimeout() time.Duration {
	manager := GetInstance()
	config := manager.GetConfig()
	return config.DefaultTimeout
}

// GetCacheSize is a convenience method to get the cache size
func GetCacheSize() int {
	manager := GetInstance()
	config := manager.GetConfig()
	return config.CacheSize
}

// IsValidationEnabled is a convenience method to check if validation is enabled
func IsValidationEnabled() bool {
	manager := GetInstance()
	config := manager.GetConfig()
	return config.EnableValidation
}

// GetLogLevel is a convenience method to get the log level
func GetLogLevel() int {
	manager := GetInstance()
	config := manager.GetConfig()
	return config.LogLevel
}