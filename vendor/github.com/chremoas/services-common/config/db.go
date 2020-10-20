package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

func (c Configuration) NewConnectionString() (string, error) {
	if !c.initialized {
		return "", errors.New("Configuration not initialized, call Load() before calling this.")
	}

	return viper.GetString("database.driver") +
		"://" +
		viper.GetString("database.username") +
		":" +
		viper.GetString("database.password") +
		"@" +
		viper.GetString("database.host") +
		":" +
		fmt.Sprintf("%d", viper.GetInt("database.port")) +
		"/" +
		viper.GetString("database.database") +
		"?" +
		viper.GetString("database.options"), nil
}
