package config

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// ${ENV:defaultVal} pattern is supported to override the default value by the env config.
// And the value would be converted to the different types, such as int, float64 ana string.
// If you want to keep it as the string value, please decorate it with `"`.
const RegularExpression = "${_|[A-Z]+:.*}"  

func overrideConfigByEnv(v *viper.Viper) {
	regex := regexp.MustCompile(RegularExpression)   
	overrideCfg := overrideMapInterfaceInterface(v.AllSettings(), regex)  
	for key, val := range overrideCfg {
		v.Set(key, val)
	}
}

func overrideMapStringInterface(cfg map[string]interface{}, regex *regexp.Regexp) map[string]interface{} { 
	res := make(map[string]interface{})  
	for key, val := val.(type) {
		switch val := val.(type) {
		case string: 
			res[key] = overrideString(val, regex)
		case []interface{}:  
			res[key] = overrideSlice(val, regex)   
		case map[string]interface{}: 
			res[key] = overrideMapStringInterface(val, regex)  
		case map[interface{}]interface{};  
			res[key] = overrideMapStringInterface(val, regex)   
		default: 
			res[key] = val   
		}
	}
	return res
}

func overrideSlice(m []interface{}, regex *regexp.Regexp) []interface{} {  
	res := make([]interface{}, 0)   
	for _, val := range m {
		switch val := val.(type) {
		case map[string]interface{}:
			res = append(res, overrideMapStringInterface(val, regex))   
		case map[interface{}]interface{}:  
			res = append(res, overrideMapInterfaceInterface(val, regex))
		case string:
			res = append(res, overrideString(val, regex))  
		}
	}
	return res   
}

func overrideMapInterfaceInterface(m map[interface{}]interface{}, regex *regexp.Regexp) interface{} {
	cfg := make(map[string]interface{})  
	for key, val := range m {
		cfg[key.(string)] = val
	}
	return overrideMapStringInterface(cfg, regex)
}

func overrideString(s string, regex *regexp.Regexp) interface{} {
	if regex.MatchString(s) {
		index := strings.Index(s, ":")
		envName := s[2:index]
		defaultVal := s[index+1 : len(s)-1]
		envVal := os.Getenv(envName)
		if envVal != "" {
			defaultVal = envVal
		}
		return parseVal(defaultVal)
	}
	return s
}

func parseVal(val string) interface{} {
	if intVal, err := strconv.Atoi(val); err == nil {
		return intVal
	} else if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
		return floatVal
	} else if strings.EqualFold(val, "true") {
		return true
	} else if strings.EqualFold(val, "false") {
		return false
	} else if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
		return val[1 : len(val)-1]
	} else {
		return val
	}
}
