package api

import (
	"bytes"
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-platform/platform-core/api/middleware"
	"github.com/WeBankPartners/wecube-platform/platform-core/api/v1/plugin"
	"github.com/WeBankPartners/wecube-platform/platform-core/api/v1/system"
	"github.com/WeBankPartners/wecube-platform/platform-core/common/db"
	"github.com/WeBankPartners/wecube-platform/platform-core/common/log"
	"github.com/WeBankPartners/wecube-platform/platform-core/models"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

type handlerFuncObj struct {
	HandlerFunc  func(c *gin.Context)
	Method       string
	Url          string
	LogOperation bool
	PreHandle    func(c *gin.Context)
	ApiCode      string
}

var (
	httpHandlerFuncList []*handlerFuncObj
	apiCodeMap          = make(map[string]string)
)

func init() {
	httpHandlerFuncList = append(httpHandlerFuncList,
		// health check
		&handlerFuncObj{Url: "/health-check", Method: "GET", HandlerFunc: healthCheck, ApiCode: "health"},
		// system-variable
		&handlerFuncObj{Url: "/system-variables/retrieve", Method: "POST", HandlerFunc: system.QuerySystemVariables, ApiCode: "query-system-variables"},
		&handlerFuncObj{Url: "/system-variables/create", Method: "POST", HandlerFunc: system.CreateSystemVariable, ApiCode: "create-system-variables"},
		&handlerFuncObj{Url: "/system-variables/update", Method: "POST", HandlerFunc: system.UpdateSystemVariable, ApiCode: "update-system-variables"},
		&handlerFuncObj{Url: "/system-variables/delete", Method: "POST", HandlerFunc: system.DeleteSystemVariable, ApiCode: "delete-system-variables"},
		// resource
		&handlerFuncObj{Url: "/resource/constants/resource-server-status", Method: "GET", HandlerFunc: system.GetResourceServerStatus, ApiCode: "get-resource-server-status"},
		&handlerFuncObj{Url: "/resource/constants/resource-server-types", Method: "GET", HandlerFunc: system.GetResourceServerTypes, ApiCode: "get-resource-server-types"},
		&handlerFuncObj{Url: "/resource/constants/resource-item-status", Method: "GET", HandlerFunc: system.GetResourceItemStatus, ApiCode: "get-resource-item-status"},
		&handlerFuncObj{Url: "/resource/constants/resource-item-types", Method: "GET", HandlerFunc: system.GetResourceItemTypes, ApiCode: "get-resource-item-types"},
		&handlerFuncObj{Url: "/resource/servers/retrieve", Method: "POST", HandlerFunc: system.QueryResourceServer, ApiCode: "query-resource-server"},
		&handlerFuncObj{Url: "/resource/items/retrieve", Method: "POST", HandlerFunc: system.QueryResourceItem, ApiCode: "query-resource-item"},
		&handlerFuncObj{Url: "/resource/servers/create", Method: "POST", HandlerFunc: system.CreateResourceServer, ApiCode: "create-resource-server"},
		&handlerFuncObj{Url: "/resource/servers/update", Method: "POST", HandlerFunc: system.UpdateResourceServer, ApiCode: "update-resource-server"},
		&handlerFuncObj{Url: "/resource/servers/delete", Method: "POST", HandlerFunc: system.DeleteResourceServer, ApiCode: "delete-resource-server"},
		&handlerFuncObj{Url: "/resource/servers/delete", Method: "POST", HandlerFunc: system.DeleteResourceServer, ApiCode: "delete-resource-server"},
		&handlerFuncObj{Url: "/resource/servers/:resourceServerId/product-serial", Method: "GET", HandlerFunc: system.GetResourceServerSerialNum, ApiCode: "get-serial-num"},
		// plugin
		&handlerFuncObj{Url: "/packages", Method: "GET", HandlerFunc: plugin.GetPackages, ApiCode: "get-packages"},
		&handlerFuncObj{Url: "/packages", Method: "POST", HandlerFunc: plugin.UploadPackage, ApiCode: "upload-packages"},
		&handlerFuncObj{Url: "/packages/:pluginPackage/dependencies", Method: "GET", HandlerFunc: plugin.GetPluginDependencies, ApiCode: "get-plugin-dependencies"},
		&handlerFuncObj{Url: "/packages/:pluginPackage/menus", Method: "GET", HandlerFunc: plugin.GetPluginMenus, ApiCode: "get-plugin-menus"},
		&handlerFuncObj{Url: "/packages/:pluginPackage/models", Method: "GET", HandlerFunc: plugin.GetPluginModels, ApiCode: "get-plugin-models"},
		&handlerFuncObj{Url: "/packages/:pluginPackage/system-parameters", Method: "GET", HandlerFunc: plugin.GetPluginSystemParameters, ApiCode: "get-plugin-system-parameters"},
		&handlerFuncObj{Url: "/packages/:pluginPackage/authorities", Method: "GET", HandlerFunc: plugin.GetPluginAuthorities, ApiCode: "get-plugin-authorities"},
		&handlerFuncObj{Url: "/packages/:pluginPackage/runtime-resource", Method: "GET", HandlerFunc: plugin.GetPluginRuntimeResources, ApiCode: "get-runtime-resource"},
		&handlerFuncObj{Url: "/packages/register/:pluginPackage", Method: "POST", HandlerFunc: plugin.RegisterPackage, ApiCode: "register-package"},
		&handlerFuncObj{Url: "/available-container-hosts", Method: "GET", HandlerFunc: plugin.GetAvailableContainerHost, ApiCode: "get-available-host"},
		&handlerFuncObj{Url: "/hosts/:hostIp/next-available-port", Method: "GET", HandlerFunc: plugin.GetHostAvailablePort, ApiCode: "get-available-port"},
		&handlerFuncObj{Url: "/packages/:pluginPackage/hosts/:hostIp/ports/:port/instance/launch", Method: "POST", HandlerFunc: plugin.LaunchPlugin, ApiCode: "launch-plugin"},
		&handlerFuncObj{Url: "/packages/instances/:pluginPackage/remove", Method: "DELETE", HandlerFunc: plugin.RemovePlugin, ApiCode: "remove-plugin"},
	)
}

func InitHttpServer() {
	middleware.InitHttpError()
	r := gin.New()
	// access log
	r.Use(httpLogHandle())
	// recover
	r.Use(gin.CustomRecovery(recoverHandle))
	// register handler func with auth
	authRouter := r.Group(models.UrlPrefix, middleware.AuthToken)
	for _, funcObj := range httpHandlerFuncList {
		apiCodeMap[fmt.Sprintf("%s_%s", funcObj.Method, funcObj.Url)] = funcObj.ApiCode
		handleFuncList := []gin.HandlerFunc{funcObj.HandlerFunc}
		if funcObj.PreHandle != nil {
			handleFuncList = append([]gin.HandlerFunc{funcObj.PreHandle}, funcObj.HandlerFunc)
		}
		switch funcObj.Method {
		case "GET":
			authRouter.GET(funcObj.Url, handleFuncList...)
			break
		case "POST":
			authRouter.POST(funcObj.Url, handleFuncList...)
			break
		case "PUT":
			authRouter.PUT(funcObj.Url, handleFuncList...)
			break
		case "DELETE":
			authRouter.DELETE(funcObj.Url, handleFuncList...)
			break
		}
	}
	r.Run(":" + models.Config.HttpServer.Port)
}

func httpLogHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestId := c.GetHeader(models.RequestIdHeader)
		transactionId := c.GetHeader(models.TransactionIdHeader)
		if requestId == "" {
			requestId = "req_" + guid.CreateGuid()
		}
		if transactionId == "" {
			transactionId = "trans_" + guid.CreateGuid()
		}
		c.Set(models.RequestIdHeader, requestId)
		c.Set(models.TransactionIdHeader, transactionId)
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body.Close()
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		c.Set(models.ContextRequestBody, string(bodyBytes))
		log.AccessLogger.Info(fmt.Sprintf("[%s] [%s] ->", requestId, transactionId), log.String("uri", c.Request.RequestURI), log.String("method", c.Request.Method), log.String("sourceIp", getRemoteIp(c)), log.String(models.ContextOperator, c.GetString(models.ContextOperator)), log.String(models.ContextRequestBody, c.GetString(models.ContextRequestBody)))
		c.Next()
		costTime := time.Now().Sub(start).Seconds() * 1000
		userId := c.GetString(models.ContextUserId)
		if log.DebugEnable {
			log.AccessLogger.Info(fmt.Sprintf("[%s] [%s] [%s] <-", requestId, transactionId, userId), log.String("uri", c.Request.RequestURI), log.String("method", c.Request.Method), log.Int("httpCode", c.Writer.Status()), log.Int(models.ContextErrorCode, c.GetInt(models.ContextErrorCode)), log.String(models.ContextErrorMessage, c.GetString(models.ContextErrorMessage)), log.Float64("costTime", costTime), log.String(models.ContextResponseBody, c.GetString(models.ContextResponseBody)))
		} else {
			log.AccessLogger.Info(fmt.Sprintf("[%s] [%s] [%s] <-", requestId, transactionId, userId), log.String("uri", c.Request.RequestURI), log.String("method", c.Request.Method), log.Int("httpCode", c.Writer.Status()), log.Int(models.ContextErrorCode, c.GetInt(models.ContextErrorCode)), log.String(models.ContextErrorMessage, c.GetString(models.ContextErrorMessage)), log.Float64("costTime", costTime))
		}
	}
}

func getRemoteIp(c *gin.Context) string {
	return c.ClientIP()
}

func recoverHandle(c *gin.Context, err interface{}) {
	var errorMessage string
	if err != nil {
		errorMessage = err.(error).Error()
	}
	log.Logger.Error("Handle recover error", log.Int("code", -2), log.String("message", errorMessage))
	c.JSON(http.StatusInternalServerError, models.HttpResponseMeta{Code: -2, Status: models.DefaultHttpErrorCode})
}

// @Summary 健康检查
// @description 健康检查
// @Tags 健康检查接口
// @Produce  application/json
// @Success 200 {object} models.ResponseJson
// @Router /health-check [get]
func healthCheck(c *gin.Context) {
	if err := db.CheckDbConnection(); err != nil {
		c.JSON(http.StatusInternalServerError, models.HttpResponseMeta{Status: models.DefaultHttpErrorCode, Message: err.Error()})
	} else {
		middleware.ReturnSuccess(c)
	}
}