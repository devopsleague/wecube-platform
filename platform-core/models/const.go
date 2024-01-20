package models

import "fmt"

const (
	GlobalProjectName = "platform"
	DateTimeFormat    = "2006-01-02 15:04:05"

	// header key
	AuthorizationHeader    = "Authorization"
	TransactionIdHeader    = "transactionId"
	RequestIdHeader        = "requestId"
	DefaultHttpErrorCode   = "ERROR"
	DefaultHttpSuccessCode = "OK"
	// context key
	ContextRequestBody  = "requestBody"
	ContextResponseBody = "responseBody"
	ContextOperator     = "operator"
	ContextRoles        = "roles"
	ContextAuth         = "auth"
	ContextAuthorities  = "authorities"
	ContextErrorCode    = "errorCode"
	ContextErrorKey     = "errorKey"
	ContextErrorMessage = "errorMessage"
	ContextUserId       = "userId"

	JwtSignKey = "authJwtSecretKey"
	AESPrefix  = "{AES}"

	// table name
	TableNameBatchExec                = "batch_execution"
	TableNameBatchExecJobs            = "batch_execution_jobs"
	TableNameBatchExecTemplate        = "batch_execution_template"
	TableNameBatchExecTemplateRole    = "batch_execution_template_role"
	TableNameBatchExecTemplateCollect = "batch_execution_template_collect"

	// batch execution
	BatchExecTemplateStatusAvailable = "available"
)

var (
	UrlPrefix = fmt.Sprintf("/%s", GlobalProjectName)
)
