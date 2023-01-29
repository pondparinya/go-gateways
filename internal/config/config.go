package config

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/spf13/viper"
)

type ConfigEntry struct {
	v *viper.Viper
}

func LoadConfigs(prefix, path, name string, configs interface{}) error {
	viper, err := initialize(prefix, path, name)
	if err != nil {
		return err
	}
	if err := viper.checkAndBindEnvironment(configs); err != nil {
		return err
	}
	if err := viper.v.Unmarshal(configs); err != nil {
		return err
	}

	return viper.checkMissing()
}

func initialize(prefix, path, name string) (*ConfigEntry, error) {
	v := viper.New()
	v.SetConfigName(name)
	v.AddConfigPath(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Fail to read config error : %v", err)
	}
	v.SetEnvPrefix(prefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return &ConfigEntry{v}, nil
}

func (e *ConfigEntry) checkAndBindEnvironment(configs interface{}) error {
	// Check type is pointer
	if pt := reflect.TypeOf(configs).Kind(); pt != reflect.Ptr {
		return fmt.Errorf("invalid type, should be pointer instead of %v", pt)
	}

	// Check type is struct
	t := reflect.Indirect(reflect.ValueOf(configs)).Type()
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("invalid type, should be struct instead of %v", t.Kind())
	}
	e.bindEnvironments(t)

	return nil
}

func (e *ConfigEntry) bindEnvironments(ptype reflect.Type, parts ...string) {
	for i := 0; i < ptype.NumField(); i++ {
		field := ptype.Field(i)
		newParts := make([]string, len(parts))

		tv, ok := field.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}

		copy(newParts, parts)
		if tv != ",squash" {
			newParts = append(newParts, tv)
		}
		switch field.Type.Kind() {
		case reflect.Struct:
			e.bindEnvironments(field.Type, newParts...)
		default:
			_ = e.v.BindEnv(strings.Join(newParts, "."))
		}
	}
}
func (e *ConfigEntry) checkMissing() error {
	var missingKeys []string
	keys := e.v.AllKeys()
	for _, v := range keys {
		if e.v.Get(v) == nil {
			missingKeys = append(missingKeys, strings.Replace(v, ".", "_", -1))
		}
	}

	if len(missingKeys) > 0 {
		sort.Strings(missingKeys)
		return fmt.Errorf("missing env: %v", strings.Join(missingKeys, ","))
	}

	return nil
}
