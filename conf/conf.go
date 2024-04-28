package conf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rwx-yxu/weather/app"
	"github.com/rwxrob/fs/file"
	"github.com/spf13/viper"
)

/*
type C struct {
	apiKey   string `mapstructure:"apiKey"`
	location string `mapstructure:"locationID"`
	region   string `mapstructure: "regionID"`
}
*/

var ErrorCannotResolvePath = errors.New("cannot resolve path")
var ErrorConfigFileExists = errors.New("config file already exists")

func Load() error {
	d := dirPath()
	if d == "" {
		return errors.New("invalid dir path")
	}
	viper.AddConfigPath(d)
	viper.SetConfigName(app.Name)
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Errorf("error reading config file, %s", err)
	}
	return nil
}

func exist(path string) bool { return file.Exists(path) }

func dirPath() string {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	return filepath.Join(userCacheDir, app.Name)
}

func filePath() string {
	name := fmt.Sprintf("%s.env", app.Name)
	return filepath.Join(dirPath(), name)
}

func Init() error {
	d := dirPath()
	if d == "" {
		return fmt.Errorf("%w for %s", ErrorCannotResolvePath, d)
	}

	f, err := os.Create(filePath())
	if err != nil {
		return fmt.Errorf("unable to initialise create file: %w", err)
	}

	defer f.Close()

	return nil
}

func SoftInit() error {
	if _, err := os.Stat(filePath()); err == nil {
		return fmt.Errorf("%w at: %s", ErrorConfigFileExists, filePath)
	}

	err := Init()
	if err != nil {
		return err
	}
	return nil
}

func keys() ([]string, error) {
	err := Load()
	if err != nil {
		return nil, err
	}

	return viper.AllKeys(), nil
}
