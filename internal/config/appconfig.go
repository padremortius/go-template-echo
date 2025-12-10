package config

type (
	App struct {
		AuthSrvPassword string `yaml:"-" json:"-" validate:"required"`
	}
)

func (c *Config) ReadPwd() error {
	pwd, err := fillPwdMap(c.SecPath)
	if err != nil {
		return err
	}

	c.AuthSrvPassword = pwd["authSrvPassword"]

	return nil
}
