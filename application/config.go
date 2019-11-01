package application

import (
	"fmt"
	"reflect"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type weegoConfig struct {
	configType  reflect.Type
	configValue reflect.Value
}

func getConfigType(app interface{}) (reflect.Type, error) {
	appVal := reflect.ValueOf(app)
	appType := appVal.Type()

	configField, ok := appType.FieldByName("Config")
	if !ok {
		return nil, fmt.Errorf("")
	}

	return configField.Type, nil
}

func newConfigInstance(configType reflect.Type) reflect.Value {
	configInstance := reflect.New(configType).Elem()
	return configInstance
}

func processConfig(app interface{}) (weegoConfig, error) {
	godotenv.Load("../.env.test", "../.env")

	configType, err := getConfigType(app)
	if err != nil {
		return weegoConfig{nil, reflect.Value{}}, err
	}

	configInstance := newConfigInstance(configType)

	err = envconfig.Process("", configInstance.Addr().Interface())

	return weegoConfig{configType, configInstance}, err
}
