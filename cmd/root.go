package cmd

import (
	"adswork/common"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var isDebug bool = false

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ads",
	Short: "ADS管理",
	Long:  "ADS工作项工时管理", // logo only show in init
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		// 初始化
		logLevel := ""
		if isDebug {
			logLevel = "debug"
		}
		common.JiaLog.InitLogger(logLevel)
	},
}

var Version AdsVersion

// Execute 将所有子命令添加到根命令并适当地设置标志。
// 这由 main.main() 调用。 它只需要对 rootCmd 发生一次。
func Execute() {
	//func Execute(smartVersion AdsVersion) {
	//Version = smartVersion
	common.JiaLog.Error(rootCmd.Execute())

}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.smartide-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Flags表示该类参数只能用于当前命令。
	rootCmd.Flags().BoolP("help", "h", false, "帮助")
	// 全局标示（Persistent Flags）会作用于其指定的命令与指定命令所有的子命令
	// Persistent Flags表示该类参数可以被用于当前命令及其子命令。
	rootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "是否开启Debug模式，在该模式下将显示更多的日志信息")

	// disable completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// override help command
	//rootCmd.SetHelpCommand(helpCmd)

	// usage template
	//usage_tempalte := strings.ReplaceAll(i18n.GetInstance().Main.Info_Usage_template, "\\n", "\n")
	//rootCmd.SetUsageTemplate(usage_tempalte)

	// custom command

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(workCmd)

	// 不允许命令直接按照名称排序
	cobra.EnableCommandSorting = false
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		//home, err := os.UserHomeDir()
		home, err := os.Getwd()
		cobra.CheckErr(err)

		// 在主目录中搜索配置名称“.smartide-cli”（无扩展）。
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		//viper.SetConfigName(".smartide-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
