package data_trans

import (
	"context"
	"fmt"
	"github.com/WeBankPartners/wecube-platform/platform-core/services/remote/monitor"
	"io/fs"
	"os"
	"time"

	"github.com/WeBankPartners/wecube-platform/platform-core/api/v1/process"
	"github.com/WeBankPartners/wecube-platform/platform-core/common/log"
	"github.com/WeBankPartners/wecube-platform/platform-core/models"
	"github.com/WeBankPartners/wecube-platform/platform-core/services/database"
	"github.com/WeBankPartners/wecube-platform/platform-core/services/remote"
)

var importFuncList []func(context.Context, *models.TransImportJobParam) (string, error)

func init() {
	importFuncList = append(importFuncList, importRole)
	importFuncList = append(importFuncList, importCmdbConfig)
	importFuncList = append(importFuncList, importPluginConfig)
	importFuncList = append(importFuncList, importWorkflow)
	importFuncList = append(importFuncList, importBatchExecution)
	importFuncList = append(importFuncList, importArtifactPackage)
	importFuncList = append(importFuncList, importMonitorBaseConfig)
	importFuncList = append(importFuncList, importTaskManComponentLibrary)
	importFuncList = append(importFuncList, importTaskManTemplate)
	importFuncList = append(importFuncList, execWorkflow)
	importFuncList = append(importFuncList, importMonitorServiceConfig)
}

// StartTransImport 执行导入
func StartTransImport(ctx context.Context, param models.ExecImportParam) (err error) {
	var transImport *models.TransImportTable
	var localPath string
	var transImportAction *models.TransImportActionTable
	if transImport, err = database.GetTransImport(ctx, param.TransImportId); err != nil {
		log.Logger.Error("GetTransImport err", log.Error(err))
		return
	}
	// 文件解压
	if localPath, err = database.DecompressExportZip(ctx, param.ExportNexusUrl, param.TransImportId); err != nil {
		log.Logger.Error("DecompressExportZip err", log.Error(err))
		return
	}
	if transImport == nil || transImport.Id == "" {
		// 初始化导入
		if err = database.InitTransImport(ctx, param.TransImportId, param.ExportNexusUrl, localPath, param.Operator); err != nil {
			log.Logger.Error("initTransImport err", log.Error(err))
			return
		}
	}
	if transImportAction, err = database.GetLatestTransImportAction(ctx, param.TransImportId); err != nil {
		log.Logger.Error("GetLatestTransImportAction err", log.Error(err))
		return
	}
	actionParam := &models.CallTransImportActionParam{
		TransImportId: param.TransImportId,
		Action:        string(models.TransImportStatusStart),
		Operator:      param.Operator,
		ActionId:      transImportAction.Id,
		DirPath:       localPath,
		Token:         param.Token,
		Language:      param.Language,
	}
	go doImportAction(ctx, actionParam)
	return
}

// StartTransImport
// 开始导入
// 1、导入角色
// 2、导入cmdb插件服务、导入cmdb数据、同步cmdb数据模型、
// 3、导入其它插件服务
// 4、导入编排
// 5、导入批量执行
// 6、导入物料包
// 7、导入监控基础类型、对象组、基础类型指标、对象组指标、对象组阈值配置、业务配置模版
// 8、导入taskman组件库
// 9. 导入taskman请求模版

// 开始执行
// 10、开始执行编排(创建资源、初始化资源、应用部署)
// 继续导入
// 11、导入监控业务配置、层级对象指标、层级对象阈值配置、自定义看板
func doImportAction(ctx context.Context, callParam *models.CallTransImportActionParam) (err error) {
	transImportJobParam, getConfigErr := database.GetTransImportWithDetail(ctx, callParam.TransImportId, false)
	if getConfigErr != nil {
		err = getConfigErr
		log.Logger.Error("GetTransImportWithDetail err", log.Error(err))
		return
	}
	if err = database.RecordTransImportAction(ctx, callParam); err != nil {
		err = fmt.Errorf("record trans import action table fail,%s ", err.Error())
		log.Logger.Error("RecordTransImportAction err", log.Error(err))
		return
	}
	transImportJobParam.DirPath = callParam.DirPath
	transImportJobParam.Token = callParam.Token
	transImportJobParam.Language = callParam.Language
	transImportJobParam.Operator = callParam.Operator
	if callParam.Action == string(models.TransImportActionStart) {
		var currentStep int
		for _, detailRow := range transImportJobParam.Details {
			if detailRow.Status == string(models.TransImportStatusNotStart) {
				currentStep = detailRow.Step
				break
			}
		}
		if currentStep == int(models.TransImportStepInitWorkflow) {
			transImportJobParam.CurrentDetail = transImportJobParam.Details[currentStep-1]
			if err = callImportFunc(ctx, transImportJobParam, execWorkflow); err != nil {
				return
			}
		} else if currentStep == int(models.TransImportStepMonitorBusiness) {
			transImportJobParam.CurrentDetail = transImportJobParam.Details[currentStep-1]
			if err = callImportFunc(ctx, transImportJobParam, importMonitorServiceConfig); err != nil {
				return
			}
		} else {
			for currentStep <= int(models.TransImportStepRequestTemplate) {
				transImportJobParam.CurrentDetail = transImportJobParam.Details[currentStep-1]
				funcObj := importFuncList[currentStep-1]
				if err = callImportFunc(ctx, transImportJobParam, funcObj); err != nil {
					break
				}
				currentStep = currentStep + 1
			}
			if err != nil {
				return
			}
		}
	}
	if err != nil {
		callParam.ErrorMsg = err.Error()
		log.Logger.Error("doImportAction fail", log.JsonObj("callParam", callParam), log.Error(err))
		database.RecordTransImportAction(ctx, callParam)
	}
	return
}

func callImportFunc(ctx context.Context, transImportJobParam *models.TransImportJobParam, funcObj func(context.Context, *models.TransImportJobParam) (string, error)) (err error) {
	if err = database.UpdateTransImportDetailStatus(ctx, transImportJobParam.TransImport.Id, transImportJobParam.CurrentDetail.Id, "doing", "", ""); err != nil {
		return
	}
	var output string
	output, err = funcObj(ctx, transImportJobParam)
	if err != nil {
		database.UpdateTransImportDetailStatus(ctx, transImportJobParam.TransImport.Id, transImportJobParam.CurrentDetail.Id, "fail", output, err.Error())
	} else {
		database.UpdateTransImportDetailStatus(ctx, transImportJobParam.TransImport.Id, transImportJobParam.CurrentDetail.Id, "success", output, "")
	}
	return
}

// 1、导入角色
func importRole(ctx context.Context, transImportParam *models.TransImportJobParam) (output string, err error) {
	// 解析role.json,导入角色
	log.Logger.Info("1. importRole start!!!")
	var roleList []*models.SimpleLocalRoleDto
	var response models.QuerySingleRolesResponse
	if err = database.ParseJsonData(fmt.Sprintf("%s/role.json", transImportParam.DirPath), &roleList); err != nil {
		return
	}
	for _, role := range roleList {
		if response, err = remote.RegisterLocalRole(role, transImportParam.Token, transImportParam.Language); err != nil {
			log.Logger.Error("RegisterLocalRole err", log.Error(err))
			return
		}
		if response.Status != "OK" || response.Data.ID == "" {
			err = fmt.Errorf("RegisterLocalRole fail,msg:%s", response.Message)
			log.Logger.Error("RegisterLocalRole fail", log.String("roleName", role.Name), log.String("msg", response.Message))
			return
		}
	}
	log.Logger.Info("1. importRole success end!!!")
	return
}

// 2.导入cmdb数据
func importCmdbConfig(ctx context.Context, transImportParam *models.TransImportJobParam) (output string, err error) {

	return
}

// 3、导入cmdb插件服务
func importPluginConfig(ctx context.Context, transImportParam *models.TransImportJobParam) (output string, err error) {

	return
}

// 4、导入编排
func importWorkflow(ctx context.Context, transImportParam *models.TransImportJobParam) (output string, err error) {
	// 解析workflow.json,导入编排
	log.Logger.Info("4. importWorkflow start!!!")
	var procDefList []*models.ProcessDefinitionDto
	if err = database.ParseJsonData(fmt.Sprintf("%s/workflow.json", transImportParam.DirPath), &procDefList); err != nil {
		return
	}
	if len(procDefList) > 0 {
		if _, err = process.ProcDefImport(ctx, procDefList, transImportParam.Operator, transImportParam.Token, transImportParam.Language); err != nil {
			return
		}
	}
	log.Logger.Info("4. importWorkflow success end!!!")
	return
}

// 5、导入批量执行
func importBatchExecution(ctx context.Context, transImportParam *models.TransImportJobParam) (output string, err error) {
	log.Logger.Info("5. importBatchExecution start!!!")
	var batchExecutionTemplateList []*models.BatchExecutionTemplate
	if err = database.ParseJsonData(fmt.Sprintf("%s/batch_execution.json", transImportParam.DirPath), &batchExecutionTemplateList); err != nil {
		return
	}
	if len(batchExecutionTemplateList) > 0 {
		if err = database.ImportTemplate(ctx, transImportParam.Operator, batchExecutionTemplateList); err != nil {
			log.Logger.Error("importBatchExecution ImportTemplate fail", log.Error(err))
			return
		}
	}
	log.Logger.Info("5. importWorkflow success end!!!")
	return
}

// 6、导入物料包
func importArtifactPackage(ctx context.Context, transImportParam *models.TransImportJobParam) (output string, err error) {

	return
}

// 7、导入监控基础类型、对象组、基础类型指标、对象组指标、对象组阈值配置、业务配置模版
func importMonitorBaseConfig(ctx context.Context, transImportParam *models.TransImportJobParam) (output string, err error) {
	log.Logger.Info("6. importMonitorBaseConfig start!!!")
	var monitorTypeList []string
	var response monitor.BatchAddTypeConfigResp
	// 导入监控基础类型
	if err = database.ParseJsonData(fmt.Sprintf("%s/monitor/monitor_type.json", transImportParam.DirPath), &monitorTypeList); err != nil {
		return
	}
	if response, err = monitor.ImportMonitorType(monitorTypeList, transImportParam.Token); err != nil {
		log.Logger.Error("ImportMonitorType fail", log.Error(err))
		return
	}
	if response.Status != "OK" {
		err = fmt.Errorf("ImportMonitorType %s", response.Message)
		log.Logger.Error("ImportMonitorType fail", log.String("msg", response.Message))
		return
	}
	log.Logger.Info("6-1. import monitorType success!")

	// 导入对象组
	err = monitor.ImportEndpointGroup(fmt.Sprintf("%s/monitor/endpoint_group.json", transImportParam.DirPath), transImportParam.Token, transImportParam.Language)
	if err != nil {
		log.Logger.Error("ImportEndpointGroup fail", log.Error(err))
		return
	}
	log.Logger.Info("6-2. import endpointGroup success!")
	// 导入基础类型指标、对象组指标
	metricPath := fmt.Sprintf("%s/monitor/metric", transImportParam.DirPath)
	var files []fs.DirEntry
	if files, err = os.ReadDir(metricPath); err != nil {
		log.Logger.Error("ReadDir fail", log.String("metricPath", metricPath), log.Error(err))
		return
	}
	// 遍历文件和子目录
	for _, file := range files {
		if file.IsDir() {
			continue
		}
	}
	log.Logger.Info("6. importMonitorBaseConfig success end!!!")
	return
}

// 8.importTaskManComponentLibrary 导入组件库
func importTaskManComponentLibrary(ctx context.Context, transImportParam *models.TransImportJobParam) (output string, err error) {
	// 判断是否要导入组件库
	log.Logger.Info("8. importTaskManComponentLibrary start!!!")
	if transImportParam.CurrentDetail == nil {
		err = fmt.Errorf("importTaskManTemplate CurrentDetail is empty")
		log.Logger.Error("err:", log.Error(err))
		return
	}
	if transImportParam.CurrentDetail.Input == "true" {
		// 导入组件库
		err = remote.ImportComponentLibrary(fmt.Sprintf("%s/component_library.json", transImportParam.DirPath), transImportParam.Token, transImportParam.Language)
		if err != nil {
			log.Logger.Error("ImportComponentLibrary err", log.Error(err))
			return
		}
	}
	log.Logger.Info("8. importTaskManComponentLibrary success end!!!")
	return
}

// 9、导入taskman模版
func importTaskManTemplate(ctx context.Context, transImportParam *models.TransImportJobParam) (output string, err error) {
	log.Logger.Info("9. importTaskManTemplate start!!!")
	// 导入模版
	err = remote.ImportRequestTemplate(fmt.Sprintf("%s/request_template.json", transImportParam.DirPath), transImportParam.Token, transImportParam.Language)
	if err != nil {
		log.Logger.Error("ImportRequestTemplate fail", log.Error(err))
		return
	}
	log.Logger.Info("8. importTaskManTemplate success end!!!")
	return
}

// 10、开始执行编排(创建资源、初始化资源、应用部署)
func execWorkflow(ctx context.Context, transImportParam *models.TransImportJobParam) (output string, err error) {

	return
}

// 11、导入监控业务配置、层级对象指标、层级对象阈值配置、自定义看板
func importMonitorServiceConfig(ctx context.Context, transImportParam *models.TransImportJobParam) (output string, err error) {

	return
}

func StartExecWorkflowCron() {
	t := time.NewTicker(5 * time.Second).C
	for {
		<-t
		doExecWorkflowDaemonJob()
	}
}

func doExecWorkflowDaemonJob() {
	procExecList, err := database.GetTransImportProcExecList()
	if err != nil {
		log.Logger.Error("doExecWorkflowDaemonJob fail with get proc exec list", log.Error(err))
		return
	}
	log.Logger.Debug("doExecWorkflowDaemonJob", log.JsonObj("procExecList", procExecList))
}
