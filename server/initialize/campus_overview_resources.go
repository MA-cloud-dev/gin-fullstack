package initialize

import (
	"errors"
	"fmt"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var campusOverviewMenuSeed = sysModel.SysBaseMenu{
	MenuLevel: 1,
	Hidden:    false,
	Path:      "campusOverview",
	Name:      "campusOverview",
	Component: "view/campus/overview/index.vue",
	Sort:      0,
	Meta: sysModel.Meta{
		Title: "数据总揽",
		Icon:  "data-line",
	},
}

var campusOverviewAPIs = []sysModel.SysApi{
	{ApiGroup: "校园数据总揽", Method: "GET", Path: "/campusOverview/getCampusOverview", Description: "获取校园数据总揽"},
}

var campusOverviewDefaultAuthorities = []uint{888, 9528}

func ensureCampusOverviewResources() {
	if global.GVA_DB == nil {
		return
	}

	if err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		campusMenu, err := ensureCampusRootMenu(tx)
		if err != nil {
			return err
		}

		childMenu, err := ensureCampusOverviewMenu(tx, campusMenu.ID)
		if err != nil {
			return err
		}

		if err := ensureCampusOverviewAPIs(tx); err != nil {
			return err
		}

		authorityIDs, err := collectCampusAuthorityIDs(tx, campusMenu.ID)
		if err != nil {
			return err
		}
		authorityIDs = appendAuthorityIDs(authorityIDs, campusOverviewDefaultAuthorities...)

		if err := ensureCampusOverviewMenuAuthorities(tx, authorityIDs, childMenu.ID); err != nil {
			return err
		}
		if err := ensureCampusOverviewCasbinRules(tx, authorityIDs); err != nil {
			return err
		}
		return nil
	}); err != nil {
		global.GVA_LOG.Error("补录校园数据总揽资源失败", zap.Error(err))
	} else {
		global.GVA_LOG.Info("校园数据总揽资源补录完成")
	}
}

func ensureCampusOverviewMenu(tx *gorm.DB, parentID uint) (sysModel.SysBaseMenu, error) {
	var menu sysModel.SysBaseMenu
	err := tx.Where("name = ?", campusOverviewMenuSeed.Name).First(&menu).Error
	if err == nil {
		updates := map[string]any{
			"menu_level": campusOverviewMenuSeed.MenuLevel,
			"parent_id":  parentID,
			"path":       campusOverviewMenuSeed.Path,
			"hidden":     campusOverviewMenuSeed.Hidden,
			"component":  campusOverviewMenuSeed.Component,
			"sort":       campusOverviewMenuSeed.Sort,
			"title":      campusOverviewMenuSeed.Meta.Title,
			"icon":       campusOverviewMenuSeed.Meta.Icon,
		}
		if err = tx.Model(&menu).Updates(updates).Error; err != nil {
			return menu, err
		}
		menu.ParentId = parentID
		menu.Path = campusOverviewMenuSeed.Path
		menu.Hidden = campusOverviewMenuSeed.Hidden
		menu.Component = campusOverviewMenuSeed.Component
		menu.Sort = campusOverviewMenuSeed.Sort
		menu.Meta = campusOverviewMenuSeed.Meta
		return menu, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return menu, err
	}

	menu = campusOverviewMenuSeed
	menu.ParentId = parentID
	return menu, tx.Create(&menu).Error
}

func ensureCampusOverviewAPIs(tx *gorm.DB) error {
	for _, api := range campusOverviewAPIs {
		var existing sysModel.SysApi
		err := tx.Where("path = ? AND method = ?", api.Path, api.Method).First(&existing).Error
		if err == nil {
			updates := map[string]any{
				"api_group":   api.ApiGroup,
				"description": api.Description,
			}
			if err = tx.Model(&existing).Updates(updates).Error; err != nil {
				return err
			}
			continue
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err = tx.Create(&api).Error; err != nil {
			return err
		}
	}
	return nil
}

func ensureCampusOverviewMenuAuthorities(tx *gorm.DB, authorityIDs []uint, menuID uint) error {
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
			MenuId:      fmt.Sprintf("%d", menuID),
			AuthorityId: fmt.Sprintf("%d", authorityID),
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
	}
	return nil
}

func ensureCampusOverviewCasbinRules(tx *gorm.DB, authorityIDs []uint) error {
	for _, authorityID := range authorityIDs {
		if authorityID == 0 {
			continue
		}
		for _, api := range campusOverviewAPIs {
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
