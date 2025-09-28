package config

import (
	"log/slog"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

func Load(config any, cfgFile string) error {
	if cfgFile != "" {

		if _, err := os.Stat(cfgFile); err != nil {
			return err
		}
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".config")
	}

	/*
	PROCESO:
		Tomas tu key: "database.host"
		Applica replacer: "database.host" → "database_host"
		Mayúsculas: "DATABASE_HOST"
		Ejecuta: os.Getenv("DATABASE_HOST")
		Te devuelve el resultado
	*/
	//viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	//viper.AutomaticEnv()

	// there is a bug, for this reason we are bypassing the env variables
	for _, key := range getAllKeys(config) {
		keyEnv := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))

		if val, ok := os.LookupEnv(keyEnv); ok {
			viper.BindEnv(key, keyEnv)
			viper.Set(key, val)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	slog.Info("Using config file", "file", viper.ConfigFileUsed())
	return viper.Unmarshal(config)
}

func getAllKeys(iface interface{}, parts ...string) []string {
	var keys []string

	ifv := reflect.ValueOf(iface)
	if ifv.Kind() == reflect.Ptr {
		ifv = ifv.Elem()
	}

	for i := 0; i < ifv.NumField(); i++ {
		v := ifv.Field(i)
		t := ifv.Type().Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}

		switch v.Kind() {
		case reflect.Struct:
			keys = append(keys, getAllKeys(v.Interface(), append(parts, tv)...)...)
		case reflect.Ptr:
			if v.IsNil() && v.CanSet() {
				v.Set(reflect.New(v.Type().Elem()))
			}
			if v.Elem().Kind() == reflect.Struct {
				keys = append(keys, getAllKeys(v.Interface(), append(parts, tv)...)...)
			}
			keys = append(keys, strings.Join(append(parts, tv), "."))
		default:
			keys = append(keys, strings.Join(append(parts, tv), "."))
		}
	}

	return keys
}
