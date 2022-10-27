package configs

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

//ViperConfig Viper的封装
//
//采用Viper来进行配置管理
type ViperConfig struct {
	Viper *viper.Viper
}

//配置名称，一组初始值
//
//请参数GetXxxViper来提供Viper
func NewViperConfig(name string, viper *viper.Viper) *ViperConfig {
	return &ViperConfig{
		Viper: viper,
	}
}

//获取一个配置项的值
func (dc *ViperConfig) Get(name string) (interface{}, error) {
	if dc.Viper.IsSet(name) {
		return dc.Viper.Get(name), nil
	}

	return 0, fmt.Errorf("%s不存在此参数", name)
}

//设置一个配置项的值
func (dc *ViperConfig) Set(name string, value interface{}, mustExist ...bool) error {
	if mustExist != nil {
		if dc.Viper.IsSet(name) {
			dc.Viper.Set(name, value)
			return nil
		}

		return fmt.Errorf("%s不存在此参数", name)
	} else {
		dc.Viper.Set(name, value)
		return nil
	}
}

//设置一组配置项的值
func (dc *ViperConfig) SetAll(defParams map[string]float64, mustExist ...bool) {
	for k, v := range defParams {
		dc.Set(k, v, mustExist...)
	}
}

//获取一组配置项的值
func (dc *ViperConfig) All() map[string]interface{} {
	return dc.Viper.AllSettings()
}

//从Yaml字符串创建一个新的Viper
func GetYamlStrViper(name string, params string) *viper.Viper {
	//自定义
	myViper := viper.New()
	myViper.SetConfigName(name) //配置的名称，不带扩展
	//key: value
	myViper.SetConfigType("yaml") // 如果配置文件的名称中没有扩展名，则为必需
	// myViper.AddConfigPath(".")    // 添加配置文件路径，可以添加多个
	// Find and read the config file
	// if err := myViper.ReadInConfig(); err != nil {
	// 	panic(fmt.Errorf("fatal error config file: %w", err))
	// }

	//读取配置源
	if err := viper.ReadConfig(bytes.NewBuffer([]byte(params))); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	//只需要调用viper.WatchConfig，viper 会自动监听配置修改。如果有修改，重新加载的配置。
	myViper.WatchConfig()

	return myViper
}

//从文件源创建
func GetFileViper(name string, params string) *viper.Viper {
	return nil
}

//从etcd源创建
func GetEtcdViper(name string, params string) *viper.Viper {
	return nil
}

//测试 viper配置扩展包
func GetTextConfig() {

	/*
		读取配置文件

		Viper 需要最少的配置，因此它知道在哪里查找配置文件。
		Viper 支持 JSON、TOML、YAML、HCL、INI、envfile 和 Java 属性文件。
		Viper 可以搜索多个路径，但目前单个 Viper 实例仅支持单个配置文件。
		Viper 不默认任何配置搜索路径，将默认决定留给应用程序。

		下面是一个如何使用 Viper 搜索和读取配置文件的示例。不需要任何特定路径，但应至少提供一个路径，其中需要配置文件。
	*/
	viper.SetConfigName("config")         // 配置的名称，不带扩展
	viper.SetConfigType("yaml")           // 如果配置文件的名称中没有扩展名，则为必需
	viper.AddConfigPath("/etc/appname/")  // 添加搜索路径
	viper.AddConfigPath("$HOME/.appname") // 添加配置文件路径，可以添加多个搜索路径
	viper.AddConfigPath(".")              // 添加搜索当前路径

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
		}

		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	/*
		编写配置文件
		从配置文件中读取很有用，但有时您希望存储在运行时所做的所有修改。为此，可以使用一堆命令，每个命令都有自己的用途：

		WriteConfig - 将当前 viper 配置写入预定义路径（如果存在）。如果没有预定义的路径，则会出错。如果存在，将覆盖当前配置文件。
		SafeWriteConfig - 将当前 viper 配置写入预定义路径。如果没有预定义的路径，则会出错。如果存在，则不会覆盖当前配置文件。
		WriteConfigAs - 将当前的 viper 配置写入给定的文件路径。将覆盖给定文件（如果存在）。
		SafeWriteConfigAs - 将当前的 viper 配置写入给定的文件路径。不会覆盖给定文件（如果存在）。
		根据经验，标有安全的所有文件都不会覆盖任何文件，如果不存在则创建，而默认行为是创建或截断。
	*/
	viper.WriteConfig() // writes current config to predefined path set by 'viper.AddConfigPath()' and 'viper.SetConfigName'
	viper.SafeWriteConfig()
	viper.WriteConfigAs("/path/to/my/.config")
	viper.SafeWriteConfigAs("/path/to/my/.config") // will error since it has already been written
	viper.SafeWriteConfigAs("/path/to/my/.other_config")

	/*
			查看和重新读取配置文件
		Viper 支持让您的应用程序在运行时实时读取配置文件的能力。

		需要重新启动服务器才能使配置生效的日子已经一去不复返了，viper 驱动的应用程序可以在运行时读取配置文件的更新而不会错过任何一个节拍。

		只需告诉 viper 实例 watchConfig。或者，您可以为 Viper 提供一个函数，以便在每次发生更改时运行。

		确保在调用之前添加所有 configPathsWatchConfig()
	*/
	//配置修改回调
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	//只需要调用viper.WatchConfig，viper 会自动监听配置修改。如果有修改，重新加载的配置。
	viper.WatchConfig()

	/*
			远程键/值存储支持
		要在 Viper 中启用远程支持，请对包进行空白导入viper/remote ：

		import _ "github.com/spf13/viper/remote"

		Viper 将读取从 Key/Value 存储（如 etcd 或 Consul）中的路径检索到的配置字符串（如 JSON、TOML、YAML、HCL 或 envfile）。这些值优先于默认值，但会被从磁盘、标志或环境变量中检索到的配置值覆盖。

		Viper 使用crypt从 K/V 存储中检索配置，这意味着您可以加密存储配置值，如果您拥有正确的 gpg 密钥环，则可以自动解密它们。加密是可选的。

		您可以将远程配置与本地配置结合使用，也可以独立使用。
	*/

	// alternatively, you can create a new viper instance.
	var runtime_viper = viper.New()

	runtime_viper.AddRemoteProvider("etcd3", "http://127.0.0.1:4001", "/config/hugo.yml")
	runtime_viper.SetConfigType("yaml") // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"

	// read from remote config the first time.
	err := runtime_viper.ReadRemoteConfig()
	if err != nil {
		log.Printf("unable to read remote config: %v", err)
	}

	// unmarshal config
	runtime_conf := make(map[string]string)
	runtime_viper.Unmarshal(&runtime_conf)

	// open a goroutine to watch remote changes forever
	go func() {
		for {
			time.Sleep(time.Second * 5) // delay after each request

			// currently, only tested with etcd support
			err := runtime_viper.WatchRemoteConfig()
			if err != nil {
				log.Printf("unable to read remote config: %v", err)
				continue
			}

			// unmarshal new config into our runtime config struct. you can also use channel
			// to implement a signal to notify the system of the changes
			runtime_viper.Unmarshal(&runtime_conf)
		}
	}()

	/*
		在 Viper 中，有几种方法可以根据值的类型获取值。存在以下功能和方法：

		Get(key string) : interface{}
		GetBool(key string) : bool
		GetFloat64(key string) : float64
		GetInt(key string) : int
		GetIntSlice(key string) : []int
		GetString(key string) : string
		GetStringMap(key string) : map[string]interface{}
		GetStringMapString(key string) : map[string]string
		GetStringSlice(key string) : []string
		GetTime(key string) : time.Time
		GetDuration(key string) : time.Duration
		IsSet(key string) : bool
		AllSettings() : map[string]interface{}

		要认识到的一件重要的事情是，如果没有找到每个 Get 函数，它将返回一个零值。为了检查给定的密钥是否存在，提供了该IsSet()方法。


		Viper 可以通过传递.键的分隔路径来访问嵌套字段：

		GetString("datastore.metric.host") // (returns "127.0.0.1")
		viper.Set("redis.port", 5381)
	*/
}
