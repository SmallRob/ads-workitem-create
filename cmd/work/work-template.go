package work

import (
	"adswork/common"
	"adswork/global"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6"
	adswork "github.com/microsoft/azure-devops-go-api/azuredevops/v6/work"
	"github.com/spf13/cobra"
)

var WorkTemplateCmd = &cobra.Command{
	Use:     "template",
	Short:   "工时模板",
	Long:    "获取指定数量的迭代(倒序)默认6周，创建工时模板",
	Aliases: []string{"t"},
	Example: `
  ads work template,t
  ads work template 6 `,
	Run: func(cmd *cobra.Command, args []string) {
		var number int
		if len(args) > 0 { // 从args中加载
			number, _ = strconv.Atoi(args[0])
			//fmt.Print(number)
		} else {
			//common.JiaLog.Error("")
			number = global.App.Config.ADS.IterationsCount
		}

		GetTeamIterations(number)

		var Iterations string

		for Iterations == "" {
			fmt.Print("输入迭代索引生成模板：")
			fmt.Scanln(&Iterations)
			if Iterations == "" {
				fmt.Print("\r")
			} else {
				if !checkScanlnIterations(Iterations) {
					Iterations = ""
					common.JiaLog.Info("输入多个迭代用,逗号分割\n")
				}
			}
		}
		//fmt.Println(Iterations)
		CreateWorkTemplate(Iterations)
	},
}

func checkScanlnIterations(Iterations string) bool {
	if len(Iterations) > 1 {
		if strings.IndexAny(",", Iterations) < 0 {
			//fmt.Print("输入多个迭代用,逗号分割：")
			//fmt.Scanln(&Iterations)
			return false
		}
	}
	return true
}

func init() {
	//WorkTypeCmd.Flags().StringP("username", "u", "", "i18nInstance.Start.Info_help_flag_username")
	//WorkTypeCmd.Flags().StringP("password", "t", "", "i18nInstance.Start.Info_help_flag_password")
	//WorkTypeCmd.Flags().IntP("port", "p", 22, "i18nInstance.Start.Info_help_flag_port")
}

var GetIterations []adswork.TeamSettingsIteration

// 获取指定数量迭代(倒序)
func GetTeamIterations(number int) *[]adswork.TeamSettingsIteration {

	organizationUrl := global.App.Config.ADS.OrganizationUrl         //
	personalAccessToken := global.App.Config.ADS.PersonalaccessToken //
	connection := azuredevops.NewPatConnection(organizationUrl, personalAccessToken)
	ctx := context.Background()

	Client, err := adswork.NewClient(ctx, connection)
	if err != nil {
		common.JiaLog.Error(err)
	}
	//project := global.App.Config.ADS.ProjectName
	//time := "current" //只支持current ,没办法查询到未来迭代
	item := &adswork.GetTeamIterationsArgs{}
	item.Project = &global.App.Config.ADS.ProjectName
	//item.Timeframe = &time

	responseValue, err := Client.GetTeamIterations(ctx, *item)
	if err != nil {
		common.JiaLog.Error(err)
	}
	data, err := json.MarshalIndent(responseValue, "", "      ") //这里返回的data值，类型是[]byte
	common.LogFile("GetTeamIterations_all", string(data))

	teamiterationscount := len(*responseValue)
	//common.JiaLog.InfoF("迭代总数：%v", len(*responseValue))

	//查询当月和上月迭代
	var tempIterations []adswork.TeamSettingsIteration
	currentmonthstart, currentmonthend := common.GetMonthDay(time.Now())

	// 2022-10-21 22:17:56
	//当前月的前面一个月 月初日期 2022-09-01 00:00:00
	pastmonthstarttime := common.ConvertLocalTime(currentmonthstart).AddDate(0, -1, 0)
	//当前月的后面一个月 月末日期 2022-11-30 23:59:59
	nextmonthtime := common.ConvertLocalTime(currentmonthend).AddDate(0, 1, -1)
	//fmt.Print(pastmonthstarttime.Format("2006-01-02 15:04:05") + "\n")
	//fmt.Print(nextmonthtime.Format("2006-01-02 15:04:05") + "\n")

	//取最后50条，减少循环次数，
	tempidx := 0
	if teamiterationscount > 50 {
		tempidx = teamiterationscount - 50
	}
	for i := tempidx; i < teamiterationscount; i++ {
		//xxx1 := *responseValue[i] //报错：'responseValue[i]' (类型 'TeamSettingsIteration')的间接引用无效，       无效运算: 'responseValue[i]' (类型 '*[]TeamSettingsIteration' 不支持索引)
		item := (*responseValue)[i]

		//After() StartDate在xxx之后 返回true       Before()   StartDate在xxx之前 返回true
		if item.Attributes.StartDate.Time.After(pastmonthstarttime) && item.Attributes.StartDate.Time.Before(nextmonthtime) {
			tempIterations = append(tempIterations, item)
		}
	}

	startidx := 0
	if number <= len(tempIterations) {
		startidx = len(tempIterations) - number
	}
	//fmt.Printf("数量 %v\n", count)
	idx := 0
	for i := startidx; i < len(tempIterations); i++ {
		item := tempIterations[i]
		GetIterations = append(GetIterations, item)
		//fmt.Printf("%v ---- %v \n", *item.Name, *item.Path)
		fmt.Printf("%v----%v \n", idx, *item.Name)
		idx++
	}

	//data, err := json.MarshalIndent(tempIterations, "", "      ") //这里返回的data值，类型是[]byte
	//common.LogFile("GetTeamIterations_11", string(data))
	//common.JiaLog.Info("获取团队所有迭代成功")

	return responseValue
}

// 通过输入迭代索引，生成工作项模板
func CreateWorkTemplate(IterationsIndex string) {

	indexlist := strings.Split(IterationsIndex, ",")

	dirname := time.Now().Format("200601")
	txtpath := filepath.Join(common.GetDirPath(dirname), FILEWORKTEMPLATE)
	isFile := common.CheckFileIsExist(txtpath)

	var redcmd string
	if isFile {
		fmt.Printf("模板文件:%v 已存在\n是否替换(y/n=替换/新建)？默认:y ", txtpath)
		fmt.Scanln(&redcmd)
		if strings.ToLower(redcmd) == "n" {
			txtpath = filepath.Join(common.GetDirPath(dirname), fmt.Sprintf("worktemplate%v.txt", dirname))
		} else {
			os.Remove(txtpath)
		}
	}
	f, _ := os.Create(txtpath) //创建文件

	for _, s := range indexlist {

		idx, _ := strconv.Atoi(strings.Trim(s, " "))
		fmt.Printf("输入：%v----%v \n", s, *GetIterations[idx].Name)
		Item := GetIterations[idx]

		days := common.GetBetweenDates(Item.Attributes.StartDate.Time, Item.Attributes.FinishDate.Time)

		fmt.Fprintln(f, fmt.Sprintf("#workid:xxxxx"))
		fmt.Fprintln(f, fmt.Sprintf("work:%v%s-%s,%v,%v", global.App.Config.ProductItem.Title, days[0].Format("2006/01/02"), days[len(days)-1].Format("2006/01/02"), global.App.Config.ProductItem.AreaPath, *Item.Path))
		for _, day := range days {
			fmt.Fprintln(f, day.Format("task:2006/01/02")+" ")
		}
	}
	f.Close()
	fmt.Println(txtpath)

}

// 获取指定月份工作日并生成工作项模板
func GetMonthWorkDays(datetime string) []WeekList {
	var weeklist []WeekList
	_, _, weekarry := common.GetMonthWeekCount(datetime)

	filename := time.Now().Format("20060102")
	dirname := strings.Replace(datetime, "-", "", -1)
	txtpath := filepath.Join(common.GetDirPath(dirname), fmt.Sprintf("工作日%s.txt", datetime))
	isFile := common.CheckFileIsExist(txtpath)
	var redcmd string
	if isFile {
		fmt.Printf("模板文件:%v 已存在\n是否替换(y/n=替换/新建)？默认:y ", txtpath)
		fmt.Scanln(&redcmd)
		if strings.ToLower(redcmd) == "n" {
			txtpath = filepath.Join(common.GetDirPath(dirname), fmt.Sprintf("工作日%v.txt", filename))
		} else {
			os.Remove(txtpath)
		}
	}

	f, _ := os.Create(txtpath) //创建文件

	for _, i := range weekarry {
		year, _ := strconv.Atoi(strings.Split(i, "_")[0])
		week, _ := strconv.Atoi(strings.Split(i, "_")[1])
		weekstart, weekend := common.WeekRange(year, week)

		// common.JiaLog.InfoF("WeekRange:%v-------%v\n", weekstart, weekend)
		days := common.GetBetweenDates(weekstart, weekend)
		weeklist = append(weeklist, WeekList{WeekId: fmt.Sprintf("%v_%v", year, week), Monday: common.ConvertJsonDate(days[0]), Sunday: common.ConvertJsonDate(days[6]), Days: common.ConvertJsonDateArray(days)})

		fmt.Fprintln(f, fmt.Sprintf("work:%v%s-%s,%v,%v", global.App.Config.ProductItem.Title, days[0].Format("2006/01/02"), days[4].Format("2006/01/02"), global.App.Config.ProductItem.AreaPath, "ProjectManagement\\FY 2022\\Week"))
		for i, day := range days {
			if i < 5 {
				fmt.Fprintln(f, day.Format("task:2006/01/02")+" ")
			}
		}
		//fmt.Fprintln(f, "\n")

	}
	f.Close()
	return weeklist
}

type WeekList struct {
	WeekId string            `json:"weekid"`
	Days   []common.JsonDate `json:"days"`
	Monday common.JsonDate   `json:"monday"`
	Sunday common.JsonDate   `json:"sunday"`
}
