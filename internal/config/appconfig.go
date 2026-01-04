package config

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/padremortius/go-template-echo/internal/storage"
	"github.com/padremortius/go-template-echo/pkgs/baseconfig"
	"github.com/padremortius/go-template-echo/pkgs/common"
	"github.com/padremortius/go-template-echo/pkgs/crontab"
	"github.com/padremortius/go-template-echo/pkgs/httpserver"
	"github.com/padremortius/go-template-echo/pkgs/svclogger"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

type App struct {
	AuthSrvPassword string `yaml:"-" json:"-" validate:"required"`
}

type Config struct {
	App                `yaml:"app" json:"app" validate:"required"`
	baseconfig.BaseApp `yaml:"baseApp" json:"baseApp" validate:"required"`
	Crontab            crontab.CronOpts   `yaml:"crontab" json:"crontab" validate:"required"`
	HTTP               httpserver.HTTP    `yaml:"http" json:"http" validate:"required"`
	Log                svclogger.Log      `yaml:"logger" json:"logger" validate:"required"`
	Storage            storage.StorageCfg `yaml:"storage" json:"storage" validate:"required"`
	baseconfig.Version `json:"version"`
}

func (c *Config) ReadBaseConfig() error {
	if err := cleanenv.ReadConfig("application.yml", c); err != nil {
		return err
	}
	return nil
}

// NewConfig initializes the configuration by reading environment variables
// and a YAML configuration file.
//
// It returns an error if there is an issue reading the environment variables
// or the configuration file.
func NewConfig() (*Config, error) {
	var cfg Config
	if err := cfg.ReadBaseConfig(); err != nil {
		return &Config{}, errors.New("NewConfig: " + err.Error())
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return &Config{}, err
	}

	if _, err := os.Stat(".env"); err == nil {
		if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
			return &Config{}, err
		}
	}

	appConfigName := fmt.Sprint(cfg.Name, "-", cfg.ProfileName, ".yml")

	if cfg.ProfileName == "dev" {
		if err := cleanenv.ReadConfig(appConfigName, &cfg); err != nil {
			return &Config{}, err
		}
	} else {
		configURL, _ := url.JoinPath(cfg.ConfSrvURI, appConfigName)

		data, err := common.GetFileByURL(configURL)
		if err != nil {
			return &Config{}, err
		}

		if err := cleanenv.ParseYAML(bytes.NewBuffer(data), &cfg); err != nil {
			return &Config{}, err
		}
	}

	if err := cfg.ReadPwd(); err != nil {
		return &Config{}, errors.New("Read password error: " + err.Error())
	}

	if err := cfg.validateConfig(); err != nil {
		return &Config{}, err
	}
	return &cfg, nil
}

func (c *Config) validateConfig() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return err
	}
	return nil
}

func (c *Config) ReadPwd() error {
	pwd, err := baseconfig.FillPwdMap(c.SecPath)
	if err != nil {
		return err
	}

	c.AuthSrvPassword = pwd["authSrvPassword"]

	return nil
}
