package cmd

import (
	"adswork/config"
	_ "embed"
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

//国际化
//https://github.com/gohouse/i18n

// initCmd represents the init command
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "版本",
	Long:    "当前版本",
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {

		//root里面定义的Version
		//common.JiaLog.Console(Version.ConvertToJson())
		PrintVersion()
	},
}

func PrintVersion() {
	fmt.Println("ADS Tool:", config.Version)
	fmt.Println("Go Version:", runtime.Version())
	fmt.Println("OS/Arch:", runtime.GOOS+"/"+runtime.GOARCH)
	fmt.Println("Git Commit:", config.Commit)
	fmt.Println("BuildTime:", config.BuildTime)
	fmt.Println("Author:", config.Author)
	fmt.Println("")
}

type AdsVersion struct {
	VersionNumber        string `json:"version_number"`
	TagName              string `json:"tag_name"`
	BuildNumber          string `json:"build_number"`
	TargetCommitish      string `json:"target_commitish"`
	TargetCommitishShort string
	BuildQuqueTime       string `json:"build_ququeTime"`
	Company              string `json:"company"`
	// 编译时间
	BuildTime time.Time
}

func (smartVersion *AdsVersion) ConvertToJson() string {
	json := smartVersion.VersionNumber
	return json
}

func init() {

}
