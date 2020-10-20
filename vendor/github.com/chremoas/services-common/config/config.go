package config

import (
	"fmt"
	"github.com/micro/go-micro"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"

	// Import the remote config driver
	_ "github.com/spf13/viper/remote"
)

type Config interface {
	Load(filename string) error
	NewConnectionString() (string, error)
	NewService(version, defaultName string) (micro.Service, error)
	AuthServiceName() (string, error)
	LookupService(serviceType string, serviceName string) (serviceFullName string)
}

type Configuration struct {
	initialized bool
	Namespace   string
	Database    struct {
		Driver         string
		Protocol       string
		Host           string
		Port           uint
		Database       string
		Username       string
		Password       string
		Options        string
		MaxConnections int `yaml:"maxConnections"`
	}
	Redis struct {
		Host     string
		Port     uint
		Password string
		Database int
	}
	Bot struct {
		BotToken        string   `yaml:"botToken"`
		DiscordServerId string   `yaml:"discordServerId"`
		BotRole         string   `yaml:"botRole"`
		IgnoredRoles    []string `yaml:"ignoredRoles"`
	}
	OAuth struct {
		ClientId         string `yaml:"clientId"`
		ClientSecret     string `yaml:"clientSecret"`
		CallBackProtocol string `yaml:"callBackProtocol"`
		CallBackHost     string `yaml:"callBackHost"`
		CallBackUrl      string `yaml:"callBackUrl"`
	} `yaml:"oauth"`
	Net struct {
		ListenHost string `yaml:"listenHost"`
		ListenPort int    `yaml:"listenPort"`
	}
	Discord struct {
		InviteUrl string `yaml:"inviteUrl"`
	} `yaml:"discord"`
	Registry struct {
		Hostname         string `yaml:"hostname"`
		Port             int    `yaml:"port"`
		RegisterTTL      int    `yaml:"registerTtl"`
		RegisterInterval int    `yaml:"registerInterval"`
	} `yaml:"registry"`
	Inputs []string `yaml:"inputs"`
	Chat   struct {
		Slack struct {
			Debug bool   `yaml:"debug"`
			Token string `yaml:"token"`
		} `yaml:"slack"`
		Discord struct {
			Token     string   `yaml:"token"`
			WhiteList []string `yaml:"whiteList"`
			Prefix    string   `yaml:"prefix"`
		} `yaml:"discord"`
	} `yaml:"chat"`
	Extensions map[interface{}]interface{} `yaml:"extensions"`
}

func (c *Configuration) Load(filename string) error {
	var fileRead, remoteRead bool
	var fileReadErr, remoteReadErr error

	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	sugar.Info("LEEEEEEEROY JENKIIIIIIINS")

	configNameSpace := os.Getenv("CONFIG_NAMESPACE")
	if configNameSpace == "" {
		configNameSpace = "default"
	}

	configType := os.Getenv("CONFIG_TYPE")
	if configType == "" {
		configType = "yaml"
	}

	viper.SetConfigFile(filename)

	if fileReadErr = viper.ReadInConfig(); fileReadErr == nil {
		sugar.Info("Successfully read local config file")
		fileRead = true
	}

	if err := viper.BindEnv("consul"); err == nil {
		consul := viper.Get("consul")

		if consul != nil {
			// TODO: This is very rigid. Let's find a better way.
			configPath := fmt.Sprintf("/%s/config", configNameSpace)
			sugar.Infof("Using %s Config: %s", configType, configPath)
			err := viper.AddRemoteProvider("consul", consul.(string), configPath)
			if err == nil {
				viper.SetConfigType(configType) // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop"

				if remoteReadErr = viper.ReadRemoteConfig(); remoteReadErr == nil {
					sugar.Info("Successfully read remote config")
					remoteRead = true
				}
			} else {
				sugar.Info(err.Error())
			}
		}
	}

	if !fileRead && !remoteRead {
		return fmt.Errorf("unable to read config:\n\tfile=%v\n\tremote=%v|n", fileReadErr, remoteReadErr)
	}

	if err := viper.Unmarshal(&c); err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}

	// Let's set a default namespace because a lot of people don't care what it actually is
	if c.Namespace == "" {
		c.Namespace = "com.aba-eve"
	}

	c.initialized = true

	return nil
}

func (c *Configuration) IsInitialized() bool {
	return c.initialized
}

func (c *Configuration) LookupService(serviceType string, serviceName string) (serviceFullName string) {
	return fmt.Sprintf("%s.%s.%s", c.Namespace, serviceType, serviceName)
}
