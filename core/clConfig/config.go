package clConfig

/**
 *   INI文件格式加载实现包
 *   BY: Ciao Lan
 *   DATE: 2017-08-04
 *
 *   支持自定义配置节点和属性
 *   支持多种数据格式读取: string, int32, int64, float32, float64
 *   支持多种注释方式: #, //, 以及多行注释方式 /* ... * /
 *   使用"惰加载"方式避免重复加载, 提高程序性能
 *
 */

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
		"time"
	)

// 新建一个配置对象
func New(filename string, autoLoad time.Duration) (*Config) {
	var config  = Config {
		fileName: filename,
		config: make(map[string] sectionType),
		autoLoad: autoLoad,
		stringMap: make(map[string] autoLoadString),
		int64Map: make( map[/*section:key*/ string] autoLoadInt64),
		int32Map: make( map[/*section:key*/ string] autoLoadInt32),
		uint64Map: make( map[/*section:key*/ string] autoLoadUint64),
		uint32Map: make( map[/*section:key*/ string] autoLoadUint32),
		float32Map: make( map[/*section:key*/ string] autoLoadFloat32),
		float64Map: make( map[/*section:key*/ string] autoLoadFloat64),
		boolMap: make( map[/*section:key*/ string] autoLoadBool),
		sectionMap: make(map[/*section*/ string] autoLoadSection),
	}

	loadFile(&config)
	if autoLoad > 0 {
		go config.autoLoadConfig()
	}
	return &config
}


/**
	自动载入机制
 */
func (config *Config) autoLoadConfig() {
	for {
		<-time.After(config.autoLoad * time.Second )
		loadFile(config)
		config.lock.Lock()
		for _, val := range config.stringMap {
			config.GetStrConfig(val.section, val.key, val.def, val.value)
		}
		config.lock.Unlock()
	}
}

func loadFile(config *Config) {
	defer doWhenErr()

	dat, err := ioutil.ReadFile(config.fileName)
	if err != nil {
		return
	}

	// 将读入的文件进行换行切割
	confArr := strings.Split(string(dat),"\n")

	// 创建一个全局的配置小节
	section := "global"
	config.config[section] = make(sectionType)

	// 是否注释
	isHelp := false
	// 遍历每一行,提取里面有用的数据
	for _,v := range confArr {
		if len(v) < 2 {
			continue
		}
		v = strings.TrimSpace(v)
		v = strings.TrimPrefix(v, "\n")

		if strings.HasPrefix(v, "#") || strings.HasPrefix(v, "//") {
			continue
		}

		if isHelp {
			if strings.HasPrefix(v, "*/") {
				isHelp = false
			}
			continue
		}

		if strings.HasPrefix(v, "/*") {
			isHelp = true
			continue
		}

		if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
			// 小节
			section = strings.TrimPrefix(strings.TrimSuffix(v,"]"),"[")
			config.config[section] = make(sectionType)
		} else {
			// 配置项
			arr := strings.SplitN(v, "=", 2)
			if len(arr) != 2 {
				continue
			}
			config.config[section][strings.TrimSpace(arr[0])] = arr[1]
		}
	}
}

// 异常捕获机制
func doWhenErr() {
	// 发生错误
	if err:=recover(); err != nil {
		fmt.Printf("错误: %v\n",err)
	}
}

// 获取字符串型的配置
func (config *Config) GetStrConfig(section string, key string, def string, value *string) bool {

	var configVal = def
	var isExists = false
	if len(config.config[section]) > 0 {
		if len(config.config[section][key]) > 0 {
			configVal = config.config[section][key]
			isExists = true
		}
	}

	if !isExists {
		// 配置不存在，尝试从环境变量中获取
		envConfig := os.Getenv(strings.ToUpper(section + "_" + key))
		if envConfig != "" {
			configVal = envConfig
			isExists = true
		}
	}

	if config.autoLoad > 0 && isExists {
		config.stringMap[section+":"+key] = autoLoadString{
			section: section,
			key: key,
			def: def,
			value: value,
		}
	}
	*value = configVal
	return isExists
}

// 获取32位浮点型的配置
func (config *Config) GetFloat32Config(section string, key string, def float32, value *float32) bool {
	var configVal = def
	var isExists = false
	if len(config.config[section]) > 0 {
		if len(config.config[section][key]) > 0 {
			fval, err := strconv.ParseFloat(config.config[section][key], 32)
			if err == nil {
				configVal = float32(fval)
			}
			isExists = true
		}
	}

	if !isExists {
		// 配置不存在，尝试从环境变量中获取
		envConfig := os.Getenv(strings.ToUpper(section + "_" + key))
		if envConfig != "" {
			fval, err := strconv.ParseFloat(envConfig, 32)
			if err == nil {
				configVal = float32(fval)
			}
			isExists = true
		}
	}

	if config.autoLoad > 0 && isExists {
		config.float32Map[section+":"+key] = autoLoadFloat32{
			section: section,
			key: key,
			def: def,
			value: value,
		}
	}
	*value = configVal
	return isExists
}

// 获取64位浮点数的配置
func (config *Config) GetFloat64Config(section string, key string, def float64, value *float64) bool {
	var configVal = def
	var isExists = false
	if len(config.config[section]) > 0 {
		if len(config.config[section][key]) > 0 {
			fval, err := strconv.ParseFloat(config.config[section][key], 64)
			if err == nil {
				configVal = fval
			}
			isExists = true
		}
	}

	if !isExists {
		// 配置不存在，尝试从环境变量中获取
		envConfig := os.Getenv(strings.ToUpper(section + "_" + key))
		if envConfig != "" {
			fval, err := strconv.ParseFloat(envConfig, 64)
			if err == nil {
				configVal = fval
			}
			isExists = true
		}
	}

	if config.autoLoad > 0 && isExists {
		config.float64Map[section+":"+key] = autoLoadFloat64{
			section: section,
			key: key,
			def: def,
			value: value,
		}
	}
	*value = configVal
	return isExists
}

// 获取32位整数型的配置
func (config *Config) GetInt32Config(section string, key string, def int32, value *int32) bool {
	var configVal = def
	var isExists = false
	if len(config.config[section]) > 0 {
		if len(config.config[section][key]) > 0 {
			ival, err := strconv.ParseInt(config.config[section][key], 0, 32)
			if err == nil {
				configVal = int32(ival)
			}
			isExists = true
		}
	}

	if !isExists {
		// 配置不存在，尝试从环境变量中获取
		envConfig := os.Getenv(strings.ToUpper(section + "_" + key))
		if envConfig != "" {
			ival, err := strconv.ParseInt(envConfig, 0, 32)
			if err == nil {
				configVal = int32(ival)
			}
			isExists = true
		}
	}

	if config.autoLoad > 0 && isExists {
		config.int32Map[section+":"+key] = autoLoadInt32{
			section: section,
			key: key,
			def: def,
			value: value,
		}
	}
	*value = configVal
	return isExists
}

// 获取32位无符号整数型的配置
// 如果配置不存在返回false, 存在返回true
func (config *Config) GetUint32Config(section string, key string, def uint32, value *uint32) bool {
	var configVal = def
	var isExists = false
	if len(config.config[section]) > 0 {
		if len(config.config[section][key]) > 0 {
			ival, err := strconv.ParseUint(config.config[section][key], 0, 32)
			if err == nil {
				configVal = uint32(ival)
			}
			isExists = true
		}
	}

	if !isExists {
		// 配置不存在，尝试从环境变量中获取
		envConfig := os.Getenv(strings.ToUpper(section + "_" + key))
		if envConfig != "" {
			ival, err := strconv.ParseInt(envConfig, 0, 32)
			if err == nil {
				configVal = uint32(ival)
			}
			isExists = true
		}
	}

	if config.autoLoad > 0 && isExists {
		config.uint32Map[section+":"+key] = autoLoadUint32{
			section: section,
			key: key,
			def: def,
			value: value,
		}
	}
	*value = configVal
	return isExists
}

// 获取64位整数型的配置
// 如果配置不存在返回false, 存在返回true
func (config *Config) GetInt64Config(section string, key string, def int64, value *int64) bool {
	var configVal = def
	var isExists = false
	if len(config.config[section]) > 0 {
		if len(config.config[section][key]) > 0 {
			ival, err := strconv.ParseInt(config.config[section][key], 0, 64)
			if err == nil {
				configVal = ival
			}
			isExists = true
		}
	}

	if !isExists {
		// 配置不存在，尝试从环境变量中获取
		envConfig := os.Getenv(strings.ToUpper(section + "_" + key))
		if envConfig != "" {
			ival, err := strconv.ParseInt(envConfig, 0, 64)
			if err == nil {
				configVal = int64(ival)
			}
			isExists = true
		}
	}


	if config.autoLoad > 0 && isExists {
		config.int64Map[section+":"+key] = autoLoadInt64{
			section: section,
			key: key,
			def: def,
			value: value,
		}
	}
	*value = configVal
	return isExists
}

// 获取64位无符号整数型的配置
func (config *Config) GetUint64Config(section string, key string, def uint64, value *uint64) bool {
	var configVal = def
	var isExists = false
	if len(config.config[section]) > 0 {
		if len(config.config[section][key]) > 0 {
			ival, err := strconv.ParseUint(config.config[section][key], 0, 64)
			if err == nil {
				configVal = ival
			}
			isExists = true
		}
	}

	if !isExists {
		// 配置不存在，尝试从环境变量中获取
		envConfig := os.Getenv(strings.ToUpper(section + "_" + key))
		if envConfig != "" {
			ival, err := strconv.ParseInt(envConfig, 0, 64)
			if err == nil {
				configVal = uint64(ival)
			}
			isExists = true
		}
	}

	if config.autoLoad > 0 && isExists {
		config.uint64Map[section+":"+key] = autoLoadUint64{
			section: section,
			key: key,
			def: def,
			value: value,
		}
	}
	*value = configVal
	return isExists
}

// 获取布尔型的配置
func (config *Config)GetBoolConfig(section string, key string, def bool, value *bool) bool {
	var configVal = def
	var isExists = false
	if len(config.config[section]) > 0 {
		if len(config.config[section][key]) > 0 {
			bval, err := strconv.ParseBool(config.config[section][key])
			if err == nil {
				configVal = bval
			}
			isExists = true
		}
	}

	if !isExists {
		// 配置不存在，尝试从环境变量中获取
		envConfig := os.Getenv(strings.ToUpper(section + "_" + key))
		if envConfig != "" {
			bval, err := strconv.ParseBool(envConfig)
			if err == nil {
				configVal = bval
			}
			isExists = true
		}
	}

	if config.autoLoad > 0 && isExists {
		config.boolMap[section+":"+key] = autoLoadBool{
			section: section,
			key: key,
			def: def,
			value: value,
		}
	}
	*value = configVal
	return isExists
}

// 获取整个SECTION结构
func (config *Config)GetFullSection(section string, value *map[string]string) bool {
	if len(config.config[section]) > 0 {
		if config.autoLoad > 0 {
			config.sectionMap[section] = autoLoadSection{
				section: section,
				value: value,
			}
		}
		*value = config.config[section]
		return true
	}
	return false
}
