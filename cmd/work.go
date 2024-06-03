package cmd

import (
	"adswork/cmd/work"
	"adswork/common"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// initCmd represents the init command
var workCmd = &cobra.Command{
	Use:     "work",
	Short:   "工时操作",
	Long:    "每周工时批量操作",
	Aliases: []string{"w"},
	Example: `
  ads work template,t
  ads work --mwt,-m 2022-10
  ads work add,a
  ads work close,c`,
	Run: func(cmd *cobra.Command, args []string) {
		//args 取值flags参数后面的值，例：work -a 啊啊 sss       "sss" 在args里面获取

		/*		fmt.Printf(`
				Aliases:
				  work, w
				Examples:
				  ads work template,t
				  ads work add,a
				  ads work close,c
				`)
		*/
		fflags := cmd.Flags()

		if fflags.Changed("mwt") {
			wd := getFlagValue(fflags, "mwt")
			work.GetMonthWorkDays(wd)
		}

		// work --file ./logs/工作日2022-03.txt
		if fflags.Changed("file") {
			file := getFlagValue(fflags, "file")
			work.CreateWorkItemFile(file)
		}
		// work --close ./logs/workid.log
		if fflags.Changed("close") {
			close := getFlagValue(fflags, "close")
			work.CloseWorkItemFile(close)
		}

		if fflags.Changed("gti") {
			work.GetTeamIterations(6)

			var Iterations string
			for Iterations == "" {
				fmt.Print("输入迭代索引生成模板：")
				fmt.Scanln(&Iterations)
				if Iterations == "" {
					fmt.Print("\r")
				}
			}
			fmt.Println(Iterations)
			work.CreateWorkTemplate(Iterations)
		}

		if fflags.Changed("gwi") {
			gwi, _ := strconv.Atoi(getFlagValue(fflags, "gwi"))
			work.GetWorkItemId(gwi)
		}

		//aa := getFlagValue(fflags, "aa")
		//fmt.Println(aa)
		//
		//wt := getFlagValue(fflags, "wt")
		//fmt.Println(wt)

		cmdstr := ""
		if len(args) > 0 { // 从args中加载
			cmdstr = args[0]
			//	common.JiaLog.Info(cmdstr)
			switch cmdstr {
			case "wt":
				//	wt := getFlagValue(fflags, "wt")
				//	fmt.Println(wt)
				//workitemtypes()
				//GetWorkItemType("")
			}
		} else {
			//common.JiaLog.Error("无参数")
		}
	},
}

func init() {

	// Flags表示该类参数只能用于当前命令。
	// work worktype 任务
	workCmd.AddCommand(work.WorkTemplateCmd)
	workCmd.AddCommand(work.WorkAddCmd)
	workCmd.AddCommand(work.WorkCloseCmd)

	//workCmd.AddCommand(work.WorkTypeCmd)

}

// 检查参数是否填写
func checkFlagRequired(fflags *pflag.FlagSet, flagName string) error {
	if !fflags.Changed(flagName) {
		return fmt.Errorf("i18nInstance.Main.Err_flag_value_required", flagName)
	}
	return nil
}

// 获取Flag值
func getFlagValue(fflags *pflag.FlagSet, flag string) string {
	value, err := fflags.GetString(flag)
	if err != nil {
		if strings.Contains(err.Error(), "flag accessed but not defined:") {
			common.JiaLog.Debug(err.Error())
		} else {
			common.JiaLog.Error(err)
		}
	}
	return value
}
