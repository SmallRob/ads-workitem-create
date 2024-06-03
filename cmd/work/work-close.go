package work

import (
	"adswork/common"
	"adswork/global"
	"bufio"
	"context"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/webapi"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/workitemtracking"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var WorkCloseCmd = &cobra.Command{
	Use:     "close",
	Short:   "关闭工作项",
	Long:    "读取默认文件/data/closeworkid.txt关闭工作项",
	Aliases: []string{"c"},
	Example: `
  ads work close,c 
  ads work close ./data/closeworkid.txt`,
	Run: func(cmd *cobra.Command, args []string) {

		dirname := time.Now().Format("200601")
		var closefilepath string
		if len(args) > 0 { // 从args中加载
			closefilepath = args[0]
		} else {
			closefilepath = filepath.Join(common.GetDirPath(dirname), FILECLOSEWORKID)
		}
		CloseWorkItemFile(closefilepath)
	},
}

func init() {

}

// 关闭工作项
func CloseWorkItem(ID, State string) *workitemtracking.WorkItem {

	organizationUrl := global.App.Config.ADS.OrganizationUrl         //
	personalAccessToken := global.App.Config.ADS.PersonalaccessToken //

	connection := azuredevops.NewPatConnection(organizationUrl, personalAccessToken)
	ctx := context.Background()

	Client, err := workitemtracking.NewClient(ctx, connection)
	if err != nil {
		common.JiaLog.Error(err)
	}

	project := global.App.Config.ADS.ProjectName

	var workjson []webapi.JsonPatchOperation
	workPathState := "/fields/System.State"

	workjson = append(workjson, webapi.JsonPatchOperation{Op: &webapi.OperationValues.Add, Path: &workPathState, Value: State})

	//wk, err := json.MarshalIndent(&workjson, "", "      ")
	//common.LogFile("log-产品积压工作项json", string(wk))

	wkid, _ := strconv.Atoi(ID)
	item := &workitemtracking.UpdateWorkItemArgs{}
	item.Document = &workjson
	item.Project = &project
	item.Id = &wkid

	responseValue, err := Client.UpdateWorkItem(ctx, *item)
	if err != nil {
		common.JiaLog.Error(err)
	}

	//data, err := json.MarshalIndent(responseValue, "", "      ") //这里返回的data值，类型是[]byte
	//common.LogFile("add", string(data))
	common.JiaLog.InfoF("工作项:%v----关闭成功", ID)

	return responseValue
}

// 关闭工作项
func CloseWorkItemFile(path string) {

	if !common.CheckFileIsExist(path) {
		common.JiaLog.Info("找不到/data/closeworkid.txt文件，请使用：ads work a 命令新增工作项")
		return
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		id := string(scanner.Text()[5:])
		CloseWorkItem(id, global.App.Config.ProductItem.State)
	}
}
