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

var campusOperationLogMenuSeed = sysModel.SysBaseMenu{
	MenuLevel: 1,
	Hidden:    false,
	Path:      "campusOperationLog",
	Name:      "campusOperationLog",
	Component: "view/campus/operationLog/index.vue",
	Sort:      9,
	Meta: sysModel.Meta{
		Title: "审核操作记录",
		Icon:  "memo",
	},
}

var campusOperationLogAPIs = []sysModel.SysApi{
	{ApiGroup: "审核操作记录", Method: "GET", Path: "/campusOperationLog/getCampusOperationLogList", Description: "获取审核操作记录列表"},
	{ApiGroup: "审核操作记录", Method: "GET", Path: "/campusOperationLog/findCampusOperationLog", Description: "获取审核操作记录详情"},
}

var campusOperationLogDefaultAuthorities = []uint{888, 9528}

func ensureCampusOperationLogResources() {
	if global.GVA_DB == nil {
		return
	}

	if err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		campusMenu, err := ensureCampusRootMenu(tx)
		if err != nil {
			return err
		}

		childMenu, err := ensureCampusOperationLogMenu(tx, campusMenu.ID)
		if err != nil {
			return err
		}

		if err := ensureCampusOperationLogAPIs(tx); err != nil {
			return err
		}

		authorityIDs, err := collectCampusAuthorityIDs(tx, campusMenu.ID)
		if err != nil {
			return err
		}
		authorityIDs = appendAuthorityIDs(authorityIDs, campusOperationLogDefaultAuthorities...)

		if err := ensureCampusOperationLogMenuAuthorities(tx, authorityIDs, childMenu.ID); err != nil {
			return err
		}
		if err := ensureCampusOperationLogCasbinRules(tx, authorityIDs); err != nil {
			return err
		}
		return nil
	}); err != nil {
		global.GVA_LOG.Error("补录校园审核操作记录资源失败", zap.Error(err))
	} else {
		global.GVA_LOG.Info("校园审核操作记录资源补录完成")
	}
}

func ensureCampusRootMenu(tx *gorm.DB) (sysModel.SysBaseMenu, error) {
	var menu sysModel.SysBaseMenu
	err := tx.Where("name = ?", "campus").First(&menu).Error
	if err == nil {
		return menu, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return menu, err
	}

	menu = sysModel.SysBaseMenu{
		MenuLevel: 0,
		ParentId:  0,
		Hidden:    false,
		Path:      "campus",
		Name:      "campus",
		Component: "view/routerHolder.vue",
		Sort:      8,
		Meta: sysModel.Meta{
			Title: "校园管理",
			Icon:  "school",
		},
	}
	return menu, tx.Create(&menu).Error
}

func ensureCampusOperationLogMenu(tx *gorm.DB, parentID uint) (sysModel.SysBaseMenu, error) {
	var menu sysModel.SysBaseMenu
	err := tx.Where("name = ?", campusOperationLogMenuSeed.Name).First(&menu).Error
	if err == nil {
		updates := map[string]any{
			"menu_level": campusOperationLogMenuSeed.MenuLevel,
			"parent_id":  parentID,
			"path":       campusOperationLogMenuSeed.Path,
			"hidden":     campusOperationLogMenuSeed.Hidden,
			"component":  campusOperationLogMenuSeed.Component,
			"sort":       campusOperationLogMenuSeed.Sort,
			"title":      campusOperationLogMenuSeed.Meta.Title,
			"icon":       campusOperationLogMenuSeed.Meta.Icon,
		}
		if err = tx.Model(&menu).Updates(updates).Error; err != nil {
			return menu, err
		}
		menu.ParentId = parentID
		menu.Path = campusOperationLogMenuSeed.Path
		menu.Hidden = campusOperationLogMenuSeed.Hidden
		menu.Component = campusOperationLogMenuSeed.Component
		menu.Sort = campusOperationLogMenuSeed.Sort
		menu.Meta = campusOperationLogMenuSeed.Meta
		return menu, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return menu, err
	}

	menu = campusOperationLogMenuSeed
	menu.ParentId = parentID
	return menu, tx.Create(&menu).Error
}

func ensureCampusOperationLogAPIs(tx *gorm.DB) error {
	for _, api := range campusOperationLogAPIs {
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

func collectCampusAuthorityIDs(tx *gorm.DB, campusMenuID uint) ([]uint, error) {
	var authorityIDs []uint
	err := tx.
		Table("sys_authority_menus AS sam").
		Joins("JOIN sys_base_menus AS sbm ON sbm.id = sam.sys_base_menu_id").
		Where("sbm.id = ? OR sbm.parent_id = ?", campusMenuID, campusMenuID).
		Distinct().
		Pluck("sam.sys_authority_authority_id", &authorityIDs).Error
	if err != nil {
		return nil, err
	}

	if len(authorityIDs) == 0 {
		var existing []uint
		if err = tx.Model(&sysModel.SysAuthority{}).Where("authority_id IN ?", campusOperationLogDefaultAuthorities).Pluck("authority_id", &existing).Error; err != nil {
			return nil, err
		}
		authorityIDs = append(authorityIDs, existing...)
	}

	return uniqueAuthorityIDs(authorityIDs), nil
}

func ensureCampusOperationLogMenuAuthorities(tx *gorm.DB, authorityIDs []uint, menuID uint) error {
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

func ensureCampusOperationLogCasbinRules(tx *gorm.DB, authorityIDs []uint) error {
	for _, authorityID := range authorityIDs {
		if authorityID == 0 {
			continue
		}
		for _, api := range campusOperationLogAPIs {
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

func appendAuthorityIDs(base []uint, extras ...uint) []uint {
	return uniqueAuthorityIDs(append(base, extras...))
}

func uniqueAuthorityIDs(values []uint) []uint {
	result := make([]uint, 0, len(values))
	seen := make(map[uint]struct{}, len(values))
	for _, value := range values {
		if value == 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
