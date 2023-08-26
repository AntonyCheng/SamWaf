package waf_service

import (
	"SamWaf/global"
	"SamWaf/model"
	"SamWaf/model/request"
	"errors"
	uuid "github.com/satori/go.uuid"
	"time"
)

type WafSystemConfigService struct{}

var WafSystemConfigServiceApp = new(WafSystemConfigService)

func (receiver *WafSystemConfigService) AddApi(wafSystemConfigAddReq request.WafSystemConfigAddReq) error {
	var bean = &model.SystemConfig{
		UserCode:       global.GWAF_USER_CODE,
		TenantId:       global.GWAF_TENANT_ID,
		Id:             uuid.NewV4().String(),
		Item:           wafSystemConfigAddReq.Item,
		Value:          wafSystemConfigAddReq.Value,
		IsSystem:       "0",
		Remarks:        wafSystemConfigAddReq.Remarks,
		HashInfo:       "",
		CreateTime:     time.Now(),
		LastUpdateTime: time.Now(),
	}
	global.GWAF_LOCAL_DB.Create(bean)
	return nil
}

func (receiver *WafSystemConfigService) CheckIsExistApi(wafSystemConfigAddReq request.WafSystemConfigAddReq) error {
	return global.GWAF_LOCAL_DB.First(&model.SystemConfig{}, "item = ? ", wafSystemConfigAddReq.Item).Error
}
func (receiver *WafSystemConfigService) ModifyApi(req request.WafSystemConfigEditReq) error {
	var sysConfig model.SystemConfig
	global.GWAF_LOCAL_DB.Where("id = ?", req.Id).Find(&sysConfig)
	if req.Id != "" && req.Item != req.Item {
		return errors.New("当前配置已经存在")
	}
	editMap := map[string]interface{}{
		"Item":             req.Item,
		"Value":            req.Value,
		"Remarks":          req.Remarks,
		"last_update_time": time.Now(),
	}
	err := global.GWAF_LOCAL_DB.Model(model.SystemConfig{}).Where("id = ?", req.Id).Updates(editMap).Error

	return err
}
func (receiver *WafSystemConfigService) GetDetailApi(req request.WafSystemConfigDetailReq) model.SystemConfig {
	var bean model.SystemConfig
	global.GWAF_LOCAL_DB.Where("id=?", req.Id).Find(&bean)
	return bean
}
func (receiver *WafSystemConfigService) GetDetailByIdApi(id string) model.SystemConfig {
	var bean model.SystemConfig
	global.GWAF_LOCAL_DB.Where("id=?", id).Find(&bean)
	return bean
}
func (receiver *WafSystemConfigService) GetDetailByItem(item string) model.SystemConfig {
	var bean model.SystemConfig
	global.GWAF_LOCAL_DB.Debug().Where("Item=?", item).Find(&bean)
	return bean
}
func (receiver *WafSystemConfigService) GetListApi(req request.WafSystemConfigSearchReq) ([]model.SystemConfig, int64, error) {
	var beans []model.SystemConfig
	var total int64 = 0
	global.GWAF_LOCAL_DB.Limit(req.PageSize).Offset(req.PageSize * (req.PageIndex - 1)).Find(&beans)
	global.GWAF_LOCAL_DB.Model(&model.SystemConfig{}).Count(&total)
	return beans, total, nil
}
func (receiver *WafSystemConfigService) DelApi(req request.WafSystemConfigDelReq) error {
	var bean model.SystemConfig
	err := global.GWAF_LOCAL_DB.Where("id = ? and is_system=0", req.Id).First(&bean).Error
	if err != nil {
		return err
	}
	err = global.GWAF_LOCAL_DB.Where("id = ? and is_system=0", req.Id).Delete(model.SystemConfig{}).Error
	return err
}