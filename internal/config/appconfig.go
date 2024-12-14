package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type (
	App struct {
		AuthSrvPassword string `yaml:"-" json:"-" validate:"required"`
	}

	pwdData map[string]string
)

func ReadPwd() error {
	fname := fmt.Sprint("./", Cfg.BaseApp.Name, ".json")
	if _, err := os.Stat(fname); err == nil {
		pwdFile, err := os.ReadFile(fname)
		if err != nil {
			return err
		}

		if err = json.Unmarshal(pwdFile, &pwd); err != nil {
			return err
		}

		Cfg.App.AuthSrvPassword = pwd["app.authSrvPassword"]
	}

	return nil
}
