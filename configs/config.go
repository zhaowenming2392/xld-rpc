package configs

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

//Config 配置
type Config map[string]map[string]string

var configFile string

//NewConfig 配置，configFile 配置文件
func NewConfig(configFileName string) *Config {
	//默认配置
	c := Config{}

	if configFileName == "" {
		configFile = "my.conf"
		fmt.Println("使用默认 my.conf 配置文件")
	} else {
		configFile = configFileName
		fmt.Println("正在使用 " + configFileName + " 配置文件")
	}

	f, err := os.OpenFile(configFile, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		fmt.Println("配置文件打开失败")
		return &c
	}

	defer f.Close()

	r := bufio.NewReader(f)

	//基本配置组
	group := "Base"
	c[group] = map[string]string{}

	for {
		//fmt.Println("==========================================")

		str, err := r.ReadString('\n') //读到一个换行就结束

		//读到末尾又没有字符串了
		if str == "" && err == io.EOF { //io.EOF 表示文件的末尾
			fmt.Println("**************配置文件读取结束*******************")
			break
		}

		//fmt.Printf("字符串是：%+v 长度是：%d \n",str,len(str))
		//fmt.Println("-------------------------------------")

		str = strings.TrimSpace(str)
		if str == "" {
			fmt.Println("空退出")
			continue
		}

		//注释
		if string(str[0]) == "#" {
			fmt.Println("注释退出")
			continue
		}

		//配置组
		if string(str[0]) == "[" {
			fmt.Println("配置组")

			str = strings.TrimLeft(str, "[")
			str = strings.TrimRight(str, "]")
			str = strings.TrimSpace(str)

			if str == "" {
				fmt.Println("空配置组退出")
				continue
			}

			//组
			if str != group {
				str = strings.Title(strings.ToLower(str))
				_, ok := c[str]
				if !ok {
					fmt.Println(str, "组开始")

					//新组，先小写再转成首字母大写
					group = str
					c[group] = map[string]string{}
				}
			}

			continue
		}

		//普通的配置项

		//分隔
		strs := strings.Split(str, "=")
		if len(strs) != 2 {
			//错误的配置
			//fmt.Println("错误的配置，忽略")
			continue
		}

		name := strings.TrimSpace(strs[0])
		value := strings.TrimSpace(strs[1])

		if name == "" {
			//错误的配置
			//fmt.Println("错误的配置，忽略")
			continue
		}

		//加入配置
		c[group][name] = value
	}

	//fmt.Println()
	//fmt.Println()
	return &c
}

//Get 获取某个配置
func (c *Config) Get(name, group string) (string, error) {
	group = strings.Title(strings.ToLower(strings.TrimSpace(group)))

	if group == "" {
		group = "Base"
	}

	_, ok := (*c)[group]

	if !ok {
		return "", errors.New("配置组不存在")
	}

	v, ok := (*c)[group][name]

	if ok {
		return v, nil
	}

	return "", errors.New("配置项不存在")
}

//Set 设置某个配置，会判断是否有值
func (c *Config) Set(name, value, group string) error {
	group = strings.Title(strings.ToLower(strings.TrimSpace(group)))

	if group == "" {
		group = "Base"
	}

	_, ok := (*c)[group]

	if !ok {
		return errors.New("配置组不存在")
	}

	(*c)[group][name] = value

	return nil
}

//Del 删除配置项
func (c *Config) Del(name, group string) {
	group = strings.Title(strings.ToLower(strings.TrimSpace(group)))

	if group == "" {
		group = "Base"
	}

	_, ok := (*c)[group]

	if ok {
		delete((*c)[group], name)
	}

	c.Save()
}

//Save 保存配置
func (c *Config) Save() {
	str := ""
	for k, v := range *c {
		str += "# " + k + " 配置组开始\n"
		str += "[ " + k + " ]\n"
		for i, m := range v {
			str += i + " = " + m + "\n"
		}

		str += "# " + k + " 配置组结束\n\n"
	}

	fmt.Println("整个配置\n", str)

	//-----------写文件
	f, err := os.OpenFile(configFile, os.O_WRONLY|os.O_CREATE, 0775)
	if err != nil {
		//fmt.Println("wu.config 打开失败")
	}

	defer f.Close()

	if err == nil {
		str += "#最后保存时间" + time.Now().Format("2006-01-02 15:04:05")
		_, err := f.WriteString(str)
		//n,err := f.WriteString(str)
		if err != nil {
			//fmt.Println(err)
		} else {
			//fmt.Println("保存成功，共",n,"字节")
		}
	}
}
