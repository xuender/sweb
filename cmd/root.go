package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/kataras/iris"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	_address = "address" // 端口地址
)

var _cfgFile string // 配置文件

var rootCmd = &cobra.Command{
	Use:   "sweb",
	Short: "Static Web Server",
	Long:  `一个简单的静态Web服务`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}
		app := iris.Default()
		address := getString(cmd, _address)
		if !strings.HasPrefix(address, ":") {
			address = ":" + address
		}
		app.StaticWeb("/", args[0])
		app.Run(iris.Addr(address))
	},
}

// Execute 运行命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// getString 读取配置String
func getString(cmd *cobra.Command, name string) string {
	f := cmd.Flag(name)
	// 命令行优先
	if f.Changed {
		return f.Value.String()
	}
	ret := viper.GetString(name)
	if ret == "" {
		return f.Value.String()
	}
	return ret
}

func init() {
	cobra.OnInitialize(func() {
		if _cfgFile != "" {
			viper.SetConfigFile(_cfgFile)
		} else {
			home, err := homedir.Dir()
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			viper.AddConfigPath(home)
			viper.AddConfigPath(".")
			viper.SetConfigName(".sweb.yml")
		}
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err == nil {
			_cfgFile = viper.ConfigFileUsed()
			log.Println("读取配置文件:", _cfgFile)
		}
	})

	pflags := rootCmd.PersistentFlags()
	pflags.StringVarP(&_cfgFile, "config", "c", "sweb.yml", "配置文件")

	flags := rootCmd.Flags()
	flags.StringP(_address, "a", "8080", "访问地址端口号")
}
