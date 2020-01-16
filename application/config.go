package application

import (
	"reflect"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type weegoConfig struct {
	configType  reflect.Type
	configValue reflect.Value
}

func getConfigType(app interface{}) ([]reflect.Type, error) {
	var configTypes []reflect.Type

	appVal := reflect.ValueOf(app)
	appType := appVal.Type()

	for i := 0; i < appType.NumField(); i++ {
		field := appType.Field(i)

		if strings.HasSuffix(field.Name, "Config") {
			configTypes = append(configTypes, field.Type)
		}
	}

	return configTypes, nil
}

func newConfigInstance(configType reflect.Type) reflect.Value {
	configInstance := reflect.New(configType).Elem()
	return configInstance
}

func processConfig(app interface{}) ([]weegoConfig, error) {
	var configObjects []weegoConfig

	err := godotenv.Load("../.env.test", "../.env")
	if err != nil {
		return configObjects, err
	}

	configTypes, err := getConfigType(app)
	if err != nil {
		return configObjects, err
	}

	for _, configType := range configTypes {
		configInstance := newConfigInstance(configType)
		err = envconfig.Process("", configInstance.Addr().Interface())
		configObjects = append(configObjects, weegoConfig{configType, configInstance})
	}

	return configObjects, err
}
