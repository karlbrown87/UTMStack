package main

import (
	"path"

	"github.com/AtlasInsideCorp/UTMStackInstaller/templates"
	"github.com/AtlasInsideCorp/UTMStackInstaller/utils"
)

func InstallNginx() error {
	env := []string{"DEBIAN_FRONTEND=noninteractive"}

	if err := utils.RunEnvCmd(env, "apt", "update"); err != nil {
		return err
	}

	if err := utils.RunEnvCmd(env, "apt", "install", "-y", "nginx"); err != nil {
		return err
	}

	return nil
}

type NginxConfig struct {
	SharedKey string
}

func ConfigureNginx(conf *Config, stack *StackConfig) error {
	c := NginxConfig{
		SharedKey: conf.InternalKey,
	}

	err := utils.GenerateConfig(c, templates.FrontEnd, path.Join(stack.FrontEndNginx, "00_nginx_panel.conf"))
	if err != nil {
		return err
	}

	err = utils.GenerateConfig(c, templates.Proxy, path.Join("/", "etc", "nginx", "sites-available", "default"))
	if err != nil {
		return err
	}

	return nil
}
