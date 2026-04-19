package initialize

import (
	"errors"
	"fmt"
	"strconv"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var campusAgentReviewTestMenuSeed = sysModel.SysBaseMenu{
	MenuLevel: 1,
	Hidden:    false,
	Path:      "campusAgentReviewTest",
	Name:      "campusAgentReviewTest",
	Component: "view/campus/agentReviewTest/index.vue",
	Sort:      2,
	Meta: sysModel.Meta{
		Title: "Agent审核测试",
		Icon:  "cpu",
	},
}

var campusAgentReviewTestDefaultAuthorities = []uint{888, 9528}

var campusAgentReviewTestAPIs = []sysModel.SysApi{
	{ApiGroup: "Agent审核测试", Method: "POST", Path: "/campusAgentReviewTest/submit", Description: "B端模拟提交校园身份审核申请"},
}

const campusAgentReviewTestLegacySubmitPath = "/campusAuth/submitCampusAuth"

func ensureCampusAgentReviewTestResources() {
	if global.GVA_DB == nil {
		return
	}

	if err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		campusMenu, err := ensureCampusRootMenu(tx)
		if err != nil {
			return err
		}

		childMenu, err := ensureCampusAgentReviewTestMenu(tx, campusMenu.ID)
		if err != nil {
			return err
		}

		if err := ensureCampusAgentReviewTestAPIs(tx); err != nil {
			return err
		}

		authorityIDs, err := collectCampusAuthorityIDs(tx, campusMenu.ID)
		if err != nil {
			return err
		}
		authorityIDs = appendAuthorityIDs(authorityIDs, campusAgentReviewTestDefaultAuthorities...)

		if err := ensureCampusAgentReviewTestMenuAuthorities(tx, authorityIDs, childMenu.ID); err != nil {
			return err
		}
		if err := ensureCampusAgentReviewTestCasbinRules(tx, authorityIDs); err != nil {
			return err
		}
		return nil
	}); err != nil {
		global.GVA_LOG.Error("补录校园 Agent 审核测试资源失败", zap.Error(err))
	} else {
		global.GVA_LOG.Info("校园 Agent 审核测试资源补录完成")
	}
}

func ensureCampusAgentReviewTestMenu(tx *gorm.DB, parentID uint) (sysModel.SysBaseMenu, error) {
	var menu sysModel.SysBaseMenu
	err := tx.Where("name = ?", campusAgentReviewTestMenuSeed.Name).First(&menu).Error
	if err == nil {
		updates := map[string]any{
			"menu_level": campusAgentReviewTestMenuSeed.MenuLevel,
			"parent_id":  parentID,
			"path":       campusAgentReviewTestMenuSeed.Path,
			"hidden":     campusAgentReviewTestMenuSeed.Hidden,
			"component":  campusAgentReviewTestMenuSeed.Component,
			"sort":       campusAgentReviewTestMenuSeed.Sort,
			"title":      campusAgentReviewTestMenuSeed.Meta.Title,
			"icon":       campusAgentReviewTestMenuSeed.Meta.Icon,
		}
		if err = tx.Model(&menu).Updates(updates).Error; err != nil {
			return menu, err
		}
		menu.ParentId = parentID
		menu.Path = campusAgentReviewTestMenuSeed.Path
		menu.Hidden = campusAgentReviewTestMenuSeed.Hidden
		menu.Component = campusAgentReviewTestMenuSeed.Component
		menu.Sort = campusAgentReviewTestMenuSeed.Sort
		menu.Meta = campusAgentReviewTestMenuSeed.Meta
		return menu, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return menu, err
	}

	menu = campusAgentReviewTestMenuSeed
	menu.ParentId = parentID
	return menu, tx.Create(&menu).Error
}

func ensureCampusAgentReviewTestAPIs(tx *gorm.DB) error {
	for _, api := range campusAgentReviewTestAPIs {
		var existing sysModel.SysApi
		err := tx.Where("path = ? AND method = ?", api.Path, api.Method).First(&existing).Error

		var legacy sysModel.SysApi
		legacyErr := tx.Where("path = ? AND method = ?", campusAgentReviewTestLegacySubmitPath, api.Method).First(&legacy).Error
		if legacyErr != nil && !errors.Is(legacyErr, gorm.ErrRecordNotFound) {
			return legacyErr
		}

		if err == nil {
			updates := map[string]any{
				"api_group":   api.ApiGroup,
				"description": api.Description,
			}
			if err = tx.Model(&existing).Updates(updates).Error; err != nil {
				return err
			}
			if legacy.ID != 0 {
				if err = tx.Delete(&legacy).Error; err != nil {
					return err
				}
			}
			continue
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if legacy.ID != 0 {
			updates := map[string]any{
				"path":        api.Path,
				"api_group":   api.ApiGroup,
				"description": api.Description,
			}
			if err = tx.Model(&legacy).Updates(updates).Error; err != nil {
				return err
			}
			continue
		}

		if err = tx.Create(&api).Error; err != nil {
			return err
		}
	}
	return nil
}

func ensureCampusAgentReviewTestMenuAuthorities(tx *gorm.DB, authorityIDs []uint, menuID uint) error {
	for _, authorityID := range authorityIDs {
		if authorityID == 0 {
			continue
		}
		var count int64
		if err := tx.Model(&sysModel.SysAuthorityMenu{}).
			Where("sys_authority_authority_id = ? AND sys_base_menu_id = ?", authorityID, menuID).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		record := sysModel.SysAuthorityMenu{
			MenuId:      strconv.FormatUint(uint64(menuID), 10),
			AuthorityId: strconv.FormatUint(uint64(authorityID), 10),
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
	}
	return nil
}

func ensureCampusAgentReviewTestCasbinRules(tx *gorm.DB, authorityIDs []uint) error {
	for _, authorityID := range authorityIDs {
		if authorityID == 0 {
			continue
		}
		for _, api := range campusAgentReviewTestAPIs {
			var count int64
			if err := tx.Model(&adapter.CasbinRule{}).
				Where("ptype = ? AND v0 = ? AND v1 = ? AND v2 = ?", "p", fmt.Sprintf("%d", authorityID), api.Path, api.Method).
				Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				continue
			}
			rule := adapter.CasbinRule{
				Ptype: "p",
				V0:    fmt.Sprintf("%d", authorityID),
				V1:    api.Path,
				V2:    api.Method,
			}
			if err := tx.Create(&rule).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
