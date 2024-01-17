package database

import (
	"context"
	"strings"
	"time"

	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-platform/platform-core/common/db"
	"github.com/WeBankPartners/wecube-platform/platform-core/common/exterror"
	"github.com/WeBankPartners/wecube-platform/platform-core/models"
)

// AddProcessDefinition 添加编排
func AddProcessDefinition(ctx context.Context, user string, param models.ProcessDefinitionParam) (draftEntity *models.ProcDef, err error) {
	draftEntity = &models.ProcDef{}
	now := time.Now()
	draftEntity.Id = "pdef_" + guid.CreateGuid()
	draftEntity.Status = string(models.Draft)
	draftEntity.Key = "pdef_key_" + guid.CreateGuid()
	draftEntity.Name = param.Name
	draftEntity.Tags = param.Tags
	draftEntity.ForPlugin = strings.Join(param.AuthPlugins, ",")
	draftEntity.Scene = param.Scene
	draftEntity.ConflictCheck = param.ConflictCheck
	draftEntity.CreatedBy = user
	draftEntity.CreatedTime = now
	draftEntity.UpdatedBy = user
	draftEntity.UpdatedTime = now
	draftEntity.RootEntity = param.RootEntity
	err = insertProcDef(ctx, draftEntity)
	if err != nil {
		return
	}
	return
}

// GetProcessDefinition 获取编排定义
func GetProcessDefinition(ctx context.Context, id string) (result *models.ProcDef, err error) {
	if id == "" {
		return
	}
	var list []*models.ProcDef
	err = db.MysqlEngine.Context(ctx).SQL("select * from proc_def where id = ?", id).Find(&list)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	if len(list) > 0 {
		result = list[0]
	}
	return
}

func GetProcessDefinitionByCondition(ctx context.Context, condition models.ProcDefCondition) (list []*models.ProcDef, err error) {
	var param []interface{}
	sql := "select * from proc_def where 1= 1"
	if condition.Name != "" {
		sql = sql + " and name = ?"
		param = append(param, condition.Name)
	}
	if condition.Key != "" {
		sql = sql + " and `key` = ?"
		param = append(param, condition.Key)
	}
	err = db.MysqlEngine.Context(ctx).SQL(sql, param...).Find(&list)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	return
}

func UpdateProcDef(ctx context.Context, procDef *models.ProcDef) (err error) {
	var actions []*db.ExecAction
	actions = append(actions, &db.ExecAction{Sql: "update proc_def set name=?,`version`=?,root_entity=?,tags=?,for_plugin=?,scene=?," +
		"conflict_check=?,updated_by=?,updated_time=? where id=?", Param: []interface{}{procDef.Name, procDef.Version, procDef.RootEntity,
		procDef.Tags, procDef.ForPlugin, procDef.Scene, procDef.ConflictCheck, procDef.UpdatedBy, procDef.UpdatedTime, procDef.Id}})
	err = db.Transaction(actions, ctx)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseExecuteError, err)
	}
	return
}

func UpdateProcDefStatus(ctx context.Context, procDef *models.ProcDef) (err error) {
	var actions []*db.ExecAction
	actions = append(actions, &db.ExecAction{Sql: "update proc_def set status=?,updated_by=?,updated_time=? where id = ?", Param: []interface{}{
		procDef.Status, procDef.UpdatedBy, procDef.UpdatedTime, procDef.Id}})
	err = db.Transaction(actions, ctx)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseExecuteError, err)
	}
	return
}
func UpdateProcDefStatusAndVersion(ctx context.Context, procDef *models.ProcDef) (err error) {
	var actions []*db.ExecAction
	actions = append(actions, &db.ExecAction{Sql: "update proc_def set status=?,version=?,updated_by=?,updated_time=? where id = ?", Param: []interface{}{
		procDef.Status, procDef.Version, procDef.UpdatedBy, procDef.UpdatedTime, procDef.Id}})
	err = db.Transaction(actions, ctx)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseExecuteError, err)
	}
	return
}

// GetProcDefNode 获取编排节点
func GetProcDefNode(ctx context.Context, procDefId, nodeId string) (result *models.ProcDefNode, err error) {
	var list []*models.ProcDefNode
	err = db.MysqlEngine.Context(ctx).SQL("select * from proc_def_node where proc_def_id = ? and node_id = ?", procDefId, nodeId).Find(&list)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	if len(list) > 0 {
		result = list[0]
	}
	return
}

// GetProcDefNodeLinkByLinkId  获取编排线
func GetProcDefNodeLinkByLinkId(ctx context.Context, id string) (result *models.ProcDefNodeLink, err error) {
	var list []*models.ProcDefNodeLink
	err = db.MysqlEngine.Context(ctx).SQL("select * from proc_def_node_link where link_id = ?", id).Find(&list)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	if len(list) > 0 {
		result = list[0]
	}
	return
}

func GetProcDefNodeByProcDefId(ctx context.Context, procDefId string) (result []*models.ProcDefNodeDto, err error) {
	var list []*models.ProcDefNode
	if procDefId == "" {
		return
	}
	err = db.MysqlEngine.Context(ctx).SQL("select * from proc_def_node where proc_def_id = ?", procDefId).Find(&list)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	for _, procDefNode := range list {
		var nodeParamList []*models.ProcDefNodeParam
		err = db.MysqlEngine.Context(ctx).SQL("select * from proc_def_node_param where proc_def_node_id = ?", procDefNode.Id).Find(&nodeParamList)
		if err != nil {
			err = exterror.Catch(exterror.New().DatabaseQueryError, err)
			return
		}
		result = append(result, models.ConvertProcDefNode2Dto(procDefNode, nodeParamList))
	}
	return
}

func GetProcDefNodeLinkBySource(ctx context.Context, source string) (list []*models.ProcDefNodeLinkDto, err error) {
	var nodeLinkList []*models.ProcDefNodeLink
	list = make([]*models.ProcDefNodeLinkDto, 0)
	err = db.MysqlEngine.Context(ctx).SQL("select * from proc_def_node_link where source = ?", source).Find(&nodeLinkList)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	if len(nodeLinkList) > 0 {
		for _, nodeLink := range nodeLinkList {
			list = append(list, models.ConvertProcDefNodeLink2Dto(nodeLink))
		}
	}
	return
}

// InsertProcDefNodeLink 添加编排节点线
func InsertProcDefNodeLink(ctx context.Context, nodeLink *models.ProcDefNodeLink) (err error) {
	var actions []*db.ExecAction
	actions = append(actions, &db.ExecAction{Sql: "insert into  proc_def_node_link(id,source,target,name,ui_style) values(?,?,?,?,?)",
		Param: []interface{}{nodeLink.Id, nodeLink.Source, nodeLink.Target, nodeLink.Name, nodeLink.UiStyle}})
	err = db.Transaction(actions, ctx)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseExecuteError, err)
	}
	return
}

// UpdateProcDefNodeLink 更新编排节点线
func UpdateProcDefNodeLink(ctx context.Context, procDefNodeLink *models.ProcDefNodeLink) (err error) {
	var actions []*db.ExecAction
	sql, params := transProcDefNodeLinkUpdateConditionToSQL(procDefNodeLink)
	actions = append(actions, &db.ExecAction{Sql: sql, Param: params})
	err = db.Transaction(actions, ctx)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseExecuteError, err)
	}
	return
}

// InsertProcDefNode 添加编排节点
func InsertProcDefNode(ctx context.Context, node *models.ProcDefNode) (err error) {
	var actions []*db.ExecAction
	actions = append(actions, &db.ExecAction{Sql: "insert into  proc_def_node(id,node_id,proc_def_id,name,description,status,node_type,service_name," +
		"dynamic_bind,bind_node_id,risk_check,routine_expression,context_param_nodes,timeout,time_config,ordered_no,ui_style,created_by,created_time," +
		"updated_by,updated_time) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", Param: []interface{}{node.Id, node.NodeId, node.ProcDefId, node.Name, node.Description,
		node.Status, node.NodeType, node.ServiceName, node.DynamicBind, node.BindNodeId, node.RiskCheck, node.RoutineExpression, node.ContextParamNodes,
		node.Timeout, node.TimeConfig, node.OrderedNo, node.UiStyle, node.CreatedBy, node.CreatedTime.Format(models.DateTimeFormat), node.UpdatedBy, node.UpdatedTime.Format(models.DateTimeFormat)}})
	err = db.Transaction(actions, ctx)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseExecuteError, err)
	}
	return
}

func DeleteProcDef(ctx context.Context, procDefId string) (err error) {
	var actions []*db.ExecAction
	actions = append(actions, &db.ExecAction{Sql: "delete  from proc_def where id=? and status = ?", Param: []interface{}{procDefId, string(models.Draft)}})
	err = db.Transaction(actions, ctx)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseExecuteError, err)
	}
	return
}

func DeleteProcDefNode(ctx context.Context, procDefId, nodeId string) (err error) {
	var actions []*db.ExecAction
	actions = append(actions, &db.ExecAction{Sql: "delete  from proc_def_node where proc_def_id=? and node_id=?", Param: []interface{}{procDefId, nodeId}})
	err = db.Transaction(actions, ctx)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseExecuteError, err)
	}
	return
}

func DeleteProcDefNodeParamByNodeId(ctx context.Context, nodeId string) (err error) {
	var actions []*db.ExecAction
	actions = append(actions, &db.ExecAction{Sql: "delete  from proc_def_node_param where proc_def_node_id=?", Param: []interface{}{nodeId}})
	err = db.Transaction(actions, ctx)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseExecuteError, err)
	}
	return
}

func DeleteProcDefNodeLinkByNode(ctx context.Context, nodeId string) (err error) {
	var actions []*db.ExecAction
	actions = append(actions, &db.ExecAction{Sql: "delete  from proc_def_node_link where source=? or target=?", Param: []interface{}{nodeId, nodeId}})
	err = db.Transaction(actions, ctx)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseExecuteError, err)
	}
	return
}

func DeleteProcDefNodeLinkByLinkId(ctx context.Context, id string) (err error) {
	var actions []*db.ExecAction
	actions = append(actions, &db.ExecAction{Sql: "delete  from proc_def_node_link where link_id= ?", Param: []interface{}{id}})
	err = db.Transaction(actions, ctx)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseExecuteError, err)
	}
	return
}

// UpdateProcDefNode 更新编排节点
func UpdateProcDefNode(ctx context.Context, procDefNode *models.ProcDefNode) (err error) {
	var actions []*db.ExecAction
	sql, params := transProcDefNodeUpdateConditionToSQL(procDefNode)
	actions = append(actions, &db.ExecAction{Sql: sql, Param: params})
	err = db.Transaction(actions, ctx)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseExecuteError, err)
	}
	return
}

// UpdateProcDefNodeStatusByProcDefId 根据编排id更新编排节点状态
func UpdateProcDefNodeStatusByProcDefId(ctx context.Context, procDefId, status string) (err error) {
	var actions []*db.ExecAction
	actions = append(actions, &db.ExecAction{Sql: "update proc_def_node set status = ? where proc_def_id= ?", Param: []interface{}{status, procDefId}})
	err = db.Transaction(actions, ctx)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseExecuteError, err)
	}
	return
}

// GetProcDefPermissionByCondition 根据条件 获取编排权限
func GetProcDefPermissionByCondition(ctx context.Context, permission models.ProcDefPermission) (list []*models.ProcDefPermission, err error) {
	var params []interface{}
	sql := "select * from proc_def_permission where 1=1"
	if permission.Id != "" {
		sql = sql + " and id = ?"
		params = append(params, permission.Id)

	}
	if permission.ProcDefId != "" {
		sql = sql + " and proc_def_id = ?"
		params = append(params, permission.ProcDefId)
	}
	if permission.RoleName != "" {
		sql = sql + " and role_name = ?"
		params = append(params, permission.RoleName)
	}
	if permission.Permission != "" {
		sql = sql + " and permission = ?"
		params = append(params, permission.Permission)
	}
	err = db.MysqlEngine.Context(ctx).SQL(sql, params...).Find(&list)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	return
}

// InsertProcDefNodeParam 添加编排节点参数
func InsertProcDefNodeParam(ctx context.Context, node *models.ProcDefNodeParam) (err error) {
	var actions []*db.ExecAction
	actions = append(actions, &db.ExecAction{Sql: "insert into  proc_def_node_param(id,proc_def_node_id,param_id,name,bind_type," +
		"value,ctx_bind_node,ctx_bind_type,ctx_bind_name) values (?,?,?,?,?,?,?,?,?)", Param: []interface{}{node.Id, node.ProcDefNodeId, node.ParamId,
		node.Name, node.BindType, node.Value, node.CtxBindNode, node.CtxBindType, node.CtxBindName}})
	err = db.Transaction(actions, ctx)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseExecuteError, err)
	}
	return
}

func DeleteProcDefNodeParam(ctx context.Context, procDefNodeId, paramId string) (err error) {
	var actions []*db.ExecAction
	actions = append(actions, &db.ExecAction{Sql: "delete from proc_def_node_param where proc_def_node_id=?,param_id=?", Param: []interface{}{procDefNodeId, paramId}})
	err = db.Transaction(actions, ctx)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseExecuteError, err)
	}
	return
}

// GetProcDefNodeParam 获取编排节点参数
func GetProcDefNodeParam(ctx context.Context, id string) (result *models.ProcDefNodeParam, err error) {
	if id == "" {
		return
	}
	var list []*models.ProcDefNodeParam
	err = db.MysqlEngine.Context(ctx).SQL("select * from proc_def_node where id = ?", id).Find(&list)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	if len(list) > 0 {
		result = list[0]
	}
	return
}

// GetProcDefNodeParamByNodeId 根据节点获取编排节点参数
func GetProcDefNodeParamByNodeId(ctx context.Context, nodeId string) (list []*models.ProcDefNodeParam, err error) {
	err = db.MysqlEngine.Context(ctx).SQL("select * from proc_def_node_param where proc_def_nnode_id = ?", nodeId).Find(&list)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	return
}

func BatchAddProcDefPermission(ctx context.Context, procDefId string, permission models.PermissionToRole) (err error) {
	var actions []*db.ExecAction
	actions = append(actions, &db.ExecAction{Sql: "delete from proc_def_permission where proc_def_id= ? ",
		Param: []interface{}{procDefId}})
	if len(permission.USE) > 0 {
		for _, roleName := range permission.USE {
			actions = append(actions, &db.ExecAction{Sql: "insert into proc_def_permission(id,proc_def_id,role_id,role_name,permission)values(" +
				"?,?,?,?,?)", Param: []interface{}{guid.CreateGuid(), procDefId, roleName, roleName, string(models.USE)}})
		}
	}
	if len(permission.MGMT) > 0 {
		for _, roleName := range permission.MGMT {
			actions = append(actions, &db.ExecAction{Sql: "insert into proc_def_permission(id,proc_def_id,role_id,role_name,permission)values(" +
				"?,?,?,?,?)", Param: []interface{}{guid.CreateGuid(), procDefId, roleName, roleName, string(models.MGMT)}})
		}
	}
	err = db.Transaction(actions, ctx)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseExecuteError, err)
	}
	return
}

func insertProcDef(ctx context.Context, procDef *models.ProcDef) (err error) {
	var actions []*db.ExecAction
	actions = append(actions, &db.ExecAction{Sql: "insert into  proc_def (id,`key`,name,`version`,root_entity,status,tags,for_plugin,scene," +
		"conflict_check,created_by,created_time,updated_by,updated_time) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)", Param: []interface{}{procDef.Id,
		procDef.Key, procDef.Name, procDef.Version, procDef.RootEntity, procDef.Status, procDef.Tags, procDef.ForPlugin, procDef.Scene,
		procDef.ConflictCheck, procDef.CreatedBy, procDef.CreatedTime.Format(models.DateTimeFormat), procDef.UpdatedBy, procDef.UpdatedTime.Format(models.DateTimeFormat)}})
	err = db.Transaction(actions, ctx)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseExecuteError, err)
	}
	return
}

func transProcDefNodeUpdateConditionToSQL(procDefNode *models.ProcDefNode) (sql string, params []interface{}) {
	sql = "update proc_def_node set id = ?"
	params = append(params, procDefNode.Id)
	if procDefNode.ProcDefId != "" {
		sql = sql + ",proc_def_id=?"
		params = append(params, procDefNode.ProcDefId)
	}
	if procDefNode.NodeId != "" {
		sql = sql + ",node_id=?"
		params = append(params, procDefNode.NodeId)
	}
	if procDefNode.Name != "" {
		sql = sql + ",name=?"
		params = append(params, procDefNode.Name)
	}
	if procDefNode.Description != "" {
		sql = sql + ",description=?"
		params = append(params, procDefNode.Description)
	}
	if procDefNode.Status != "" {
		sql = sql + ",status=?"
		params = append(params, procDefNode.Status)
	}
	if procDefNode.NodeType != "" {
		sql = sql + ",node_type=?"
		params = append(params, procDefNode.NodeType)
	}
	if procDefNode.ServiceName != "" {
		sql = sql + ",service_name=?"
		params = append(params, procDefNode.ServiceName)
	}
	sql = sql + ",dynamic_bind=?"
	params = append(params, procDefNode.DynamicBind)
	if procDefNode.BindNodeId != "" {
		sql = sql + ",bind_node_id=?"
		params = append(params, procDefNode.BindNodeId)
	}
	sql = sql + ",risk_check=?"
	params = append(params, procDefNode.RiskCheck)
	if procDefNode.RoutineExpression != "" {
		sql = sql + ",routine_expression=?"
		params = append(params, procDefNode.RoutineExpression)
	}
	if procDefNode.ContextParamNodes != "" {
		sql = sql + ",context_param_nodes=?"
		params = append(params, procDefNode.ContextParamNodes)
	}
	if procDefNode.Timeout != 0 {
		sql = sql + ",timeout=?"
		params = append(params, procDefNode.Timeout)
	}
	if procDefNode.TimeConfig != "" {
		sql = sql + ",time_config=?"
		params = append(params, procDefNode.TimeConfig)
	}
	if procDefNode.OrderedNo != 0 {
		sql = sql + ",ordered_no=?"
		params = append(params, procDefNode.OrderedNo)
	}
	if procDefNode.UiStyle != "" {
		sql = sql + ",ui_style=?"
		params = append(params, procDefNode.UiStyle)
	}
	if procDefNode.UpdatedBy != "" {
		sql = sql + ",updated_by=?"
		params = append(params, procDefNode.UpdatedBy)
	}
	sql = sql + ",updated_time=?"
	params = append(params, procDefNode.UpdatedTime.Format(models.DateTimeFormat))
	sql = sql + " where id= ?"
	params = append(params, procDefNode.Id)
	return
}

func transProcDefNodeLinkUpdateConditionToSQL(procDefNodeLink *models.ProcDefNodeLink) (sql string, params []interface{}) {
	sql = "update proc_def_node_link set id=?"
	params = append(params, procDefNodeLink.Id)
	if procDefNodeLink.LinkId != "" {
		sql = sql + ",link_id=?"
		params = append(params, procDefNodeLink.LinkId)
	}
	if procDefNodeLink.Source != "" {
		sql = sql + ",source=?"
		params = append(params, procDefNodeLink.Source)
	}
	if procDefNodeLink.Target != "" {
		sql = sql + ",target=?"
		params = append(params, procDefNodeLink.Target)
	}
	if procDefNodeLink.Name != "" {
		sql = sql + ",name=?"
		params = append(params, procDefNodeLink.Name)
	}
	if procDefNodeLink.UiStyle != "" {
		sql = sql + ",ui_style=?"
		params = append(params, procDefNodeLink.UiStyle)
	}
	sql = sql + " where id= ?"
	params = append(params, procDefNodeLink.Id)
	return
}
