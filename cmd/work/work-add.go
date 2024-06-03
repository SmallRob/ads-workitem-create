package work

import (
	"adswork/common"
	"adswork/global"
	"bufio"
	"context"
	"fmt"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/webapi"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/workitemtracking"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var WorkAddCmd = &cobra.Command{
	Use:     "add",
	Short:   "创建工作项",
	Long:    "读取默认文件/data/worktemplate.txt创建工作项",
	Aliases: []string{"a"},
	Example: `
  ads work add,a 
  ads work add ./data/worktemplate.txt`,
	Run: func(cmd *cobra.Command, args []string) {

		dirname := time.Now().Format("200601")
		var templatefilepath string
		if len(args) > 0 { // 从args中加载
			templatefilepath = args[0]
		} else {
			templatefilepath = filepath.Join(common.GetDirPath(dirname), FILEWORKTEMPLATE)
		}
		CreateWorkItemFile(templatefilepath)
	},
}

func init() {

}

var (
	FILECLOSEWORKID  = "closeworkid.txt"
	FILEWORKTEMPLATE = "worktemplate.txt"
)

// 解析文件创建工作项
func CreateWorkItemFile(path string) {

	if !common.CheckFileIsExist(path) {
		common.JiaLog.Info("找不到/data/worktemplate.txt文件，请使用：ads work t 命令生成模板文件")
		return
	}

	//https://studygolang.com/articles/26118?fr=sidebar
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	dirname := time.Now().Format("200601")
	closepath := filepath.Join(common.GetDirPath(dirname), FILECLOSEWORKID)

	var parent_url string
	var areapath string
	var iterationpath string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		rowstr := strings.Trim(scanner.Text(), "")
		if len(rowstr) == 0 {
			continue
		}
		workidx := strings.Index(rowstr, ":") + 1
		//切片默认是根据 byte 进行切割的，中文是 3 个 byte 组成，导致上面残留两个多余的 byte
		//https://blog.csdn.net/weixin_30443747/article/details/96759281
		lineText := []rune(rowstr) //解决中文切片取值乱码
		indextype := string(lineText[:workidx])

		switch indextype {
		case "work:":
			workarry := strings.Split(rowstr, ",")
			title := []rune(workarry[0])
			titlestr := string(title[workidx:])
			areapath = workarry[1]
			iterationpath = workarry[2]
			parent := CreateWorkItem(titlestr, areapath, iterationpath)
			parent_url = *parent.Url
			common.FileAppend(closepath, fmt.Sprintf("work:%v", *parent.Id))
		case "workid:":
			workid, _ := strconv.Atoi(rowstr[workidx:])
			workinfo := GetWorkItemId(workid)

			parent_url = *workinfo.Url
			fields := *workinfo.Fields
			for k, v := range fields {
				//fmt.Printf("key:%v,val:%v\r\n", k, v)
				if k == "System.AreaPath" {
					areapath = v.(string)
				}
				if k == "System.IterationPath" {
					iterationpath = v.(string)
				}
			}
			common.FileAppend(closepath, fmt.Sprintf("work:%v", workid))
		case "task:":
			//fmt.Printf("%v--%v--%v--%v\n", string(lineText[5:]), parent_url, areapath, iterationpath)
			//Url := "https://tfs.devopshub.cn/leansoft/fdb32917-2525-4167-8242-152f6f4d0687/_apis/wit/workItems/14513"
			title := string(lineText[workidx:])
			billabledate := string(lineText[workidx:15])
			task := CreateWorkItemTask(title, billabledate, areapath, iterationpath, parent_url)
			//关闭工作项
			common.FileAppend(closepath, fmt.Sprintf("task:%v", *task.Id))
		default:
			common.JiaLog.InfoF("已注释：%v", rowstr)
		}

	}
}

// 获取指定工作项ID相关属性
func GetWorkItemId(workid int) *workitemtracking.WorkItem {

	organizationUrl := global.App.Config.ADS.OrganizationUrl
	personalAccessToken := global.App.Config.ADS.PersonalaccessToken
	connection := azuredevops.NewPatConnection(organizationUrl, personalAccessToken)

	ctx := context.Background()
	// Create a client to interact with the Core area
	Client, err := workitemtracking.NewClient(ctx, connection)
	if err != nil {
		log.Fatal(err)
	}

	project := global.App.Config.ADS.ProjectName
	responseValue, err := Client.GetWorkItem(ctx, workitemtracking.GetWorkItemArgs{Project: &project, Id: &workid})
	if err != nil {
		log.Fatal(err)
	}

	//data, err := json.MarshalIndent(responseValue, "", "      ") //这里返回的data值，类型是[]byte
	//common.LogFile("getworkid", string(data))
	return responseValue
}

// 创建产品积压工作项
func CreateWorkItem(Title, AreaPath, IterationPath string) *workitemtracking.WorkItem {

	organizationUrl := global.App.Config.ADS.OrganizationUrl         //
	personalAccessToken := global.App.Config.ADS.PersonalaccessToken //

	connection := azuredevops.NewPatConnection(organizationUrl, personalAccessToken)
	ctx := context.Background()

	Client, err := workitemtracking.NewClient(ctx, connection)
	if err != nil {
		common.JiaLog.Error(err)
	}

	project := global.App.Config.ADS.ProjectName
	workype := "产品积压工作项"

	var workjson []webapi.JsonPatchOperation
	workPathTitle := "/fields/System.Title"
	workPathAreaPath := "/fields/System.AreaPath"
	workPathIterationPath := "/fields/System.IterationPath"
	workPathAssignedTo := "/fields/System.AssignedTo"

	workjson = append(workjson, webapi.JsonPatchOperation{Op: &webapi.OperationValues.Add, Path: &workPathTitle, Value: Title})
	workjson = append(workjson, webapi.JsonPatchOperation{Op: &webapi.OperationValues.Add, Path: &workPathAreaPath, Value: AreaPath})
	workjson = append(workjson, webapi.JsonPatchOperation{Op: &webapi.OperationValues.Add, Path: &workPathIterationPath, Value: IterationPath})
	workjson = append(workjson, webapi.JsonPatchOperation{Op: &webapi.OperationValues.Add, Path: &workPathAssignedTo, Value: global.App.Config.Task.AssignedTo})

	//wk, err := json.MarshalIndent(&workjson, "", "      ")
	//common.LogFile("log-产品积压工作项json", string(wk))

	item := &workitemtracking.CreateWorkItemArgs{}
	item.Document = &workjson
	item.Project = &project
	item.Type = &workype

	responseValue, err := Client.CreateWorkItem(ctx, *item)
	if err != nil {
		common.JiaLog.Error(err)
	}

	//data, err := json.MarshalIndent(responseValue, "", "      ") //这里返回的data值，类型是[]byte
	//common.LogFile("add", string(data))
	common.JiaLog.InfoF("产品积压工作项:%v----创建成功", Title)

	return responseValue
}

// 创建任务，需要父链接
func CreateWorkItemTask(Title, BillableDate, AreaPath, IterationPath, Parent_Url string) *workitemtracking.WorkItem {

	organizationUrl := global.App.Config.ADS.OrganizationUrl         //
	personalAccessToken := global.App.Config.ADS.PersonalaccessToken //

	connection := azuredevops.NewPatConnection(organizationUrl, personalAccessToken)
	ctx := context.Background()

	Client, err := workitemtracking.NewClient(ctx, connection)
	if err != nil {
		common.JiaLog.Error(err)
	}

	project := global.App.Config.ADS.ProjectName
	workype := "任务"

	var workjson []webapi.JsonPatchOperation
	workPathTitle := "/fields/System.Title"
	workPathAreaPath := "/fields/System.AreaPath"
	workPathIterationPath := "/fields/System.IterationPath"
	workPathAssignedTo := "/fields/System.AssignedTo"

	workPathOriginalEstimate := "/fields/Microsoft.VSTS.Scheduling.OriginalEstimate" //初始估计
	workPathCompletedWork := "/fields/Microsoft.VSTS.Scheduling.CompletedWork"       //已完成工作
	workPathBillableDate := "/fields/Leansoft.ProjectManagement.BillableDate"

	workPathWorkType := "/fields/Leansoft.ProjectManagement.WorkType"
	workPathBillableStatus := "/fields/Leansoft.ProjectManagement.BillableStatus"
	workPathTravel := "/fields/Leansoft.ProjectManagement.Travel"
	workPathTravelLocation := "/fields/Leansoft.ProjectManagement.TravelLocation"

	workjson = append(workjson, webapi.JsonPatchOperation{Op: &webapi.OperationValues.Add, Path: &workPathTitle, Value: Title})
	workjson = append(workjson, webapi.JsonPatchOperation{Op: &webapi.OperationValues.Add, Path: &workPathBillableDate, Value: BillableDate})

	workjson = append(workjson, webapi.JsonPatchOperation{Op: &webapi.OperationValues.Add, Path: &workPathAreaPath, Value: AreaPath})
	workjson = append(workjson, webapi.JsonPatchOperation{Op: &webapi.OperationValues.Add, Path: &workPathIterationPath, Value: IterationPath})
	workjson = append(workjson, webapi.JsonPatchOperation{Op: &webapi.OperationValues.Add, Path: &workPathAssignedTo, Value: global.App.Config.Task.AssignedTo})

	workjson = append(workjson, webapi.JsonPatchOperation{Op: &webapi.OperationValues.Add, Path: &workPathOriginalEstimate, Value: global.App.Config.Task.OriginalEstimate})
	workjson = append(workjson, webapi.JsonPatchOperation{Op: &webapi.OperationValues.Add, Path: &workPathCompletedWork, Value: global.App.Config.Task.CompletedWork})
	workjson = append(workjson, webapi.JsonPatchOperation{Op: &webapi.OperationValues.Add, Path: &workPathWorkType, Value: global.App.Config.Task.WorkType})
	workjson = append(workjson, webapi.JsonPatchOperation{Op: &webapi.OperationValues.Add, Path: &workPathBillableStatus, Value: global.App.Config.Task.BillableStatus})
	workjson = append(workjson, webapi.JsonPatchOperation{Op: &webapi.OperationValues.Add, Path: &workPathTravel, Value: global.App.Config.Task.Travel})
	if global.App.Config.Task.Travel == "是" {
		workjson = append(workjson, webapi.JsonPatchOperation{Op: &webapi.OperationValues.Add, Path: &workPathTravelLocation, Value: global.App.Config.Task.TravelLocation})
	}

	/*
			链接父工作项
			{
				"op": "add",
				"path": "/relations/-",
				"value": {
				"rel": "System.LinkTypes.Hierarchy-Reverse",
		        "url": "https://marinaliu.visualstudio.com/f7855e29-6f8d-429d-8c9b-41fd4d7e70a4/_apis/wit/workItems/184"
			      }
			 }
	*/

	workPath_relations := "/relations/-"
	parent := &ParentUrl{}
	parent.Rel = "System.LinkTypes.Hierarchy-Reverse"
	//	parent.Url = "https://tfs.devopshub.cn/leansoft/fdb32917-2525-4167-8242-152f6f4d0687/_apis/wit/workItems/14506"
	parent.Url = Parent_Url

	//parent := "{\"rel\": \"System.LinkTypes.Hierarchy-Reverse\",\"url\": \"https://tfs.devopshub.cn/leansoft/fdb32917-2525-4167-8242-152f6f4d0687/_apis/wit/workItems/14506\"}"
	workjson = append(workjson, webapi.JsonPatchOperation{Op: &webapi.OperationValues.Add, Path: &workPath_relations, Value: parent})

	//wk, err := json.MarshalIndent(&workjson, "", "      ")
	//common.LogFile("", string(wk))

	item := &workitemtracking.CreateWorkItemArgs{}
	item.Document = &workjson
	item.Project = &project
	item.Type = &workype

	responseValue, err := Client.CreateWorkItem(ctx, *item)
	if err != nil {
		//log.Fatal(err)
		common.JiaLog.Error(err)
	}

	//data, err := json.MarshalIndent(responseValue, "", "      ") //这里返回的data值，类型是[]byte
	//common.LogFile("addtask", string(data))

	common.JiaLog.InfoF("任务:%v----创建成功", Title)
	return responseValue
}

type ParentUrl struct {
	//注意：需要公共访问的首字母大写开头，私有小写开头
	Rel string `json:"rel"`
	Url string `json:"url"`
}
