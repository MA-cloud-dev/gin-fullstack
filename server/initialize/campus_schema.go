package initialize

import (
	"errors"
	"strings"

	campusSchema "github.com/flipped-aurora/gin-vue-admin/server/model/campus/schema"
	"gorm.io/gorm"
)

type campusCategorySeed struct {
	Name      string
	SortOrder int
}

var defaultCampusCategorySeeds = []campusCategorySeed{
	{Name: "教材书籍", SortOrder: 10},
	{Name: "数码电子", SortOrder: 20},
	{Name: "宿舍日用", SortOrder: 30},
	{Name: "服饰鞋包", SortOrder: 40},
	{Name: "运动户外", SortOrder: 50},
	{Name: "其他闲置", SortOrder: 60},
}

func migrateCampusTables(db *gorm.DB) error {
	for _, table := range campusSchemaTables() {
		if err := db.AutoMigrate(table); err != nil {
			return err
		}
	}
	return nil
}

func hasCampusTables(db *gorm.DB) bool {
	for _, table := range campusSchemaTables() {
		if !db.Migrator().HasTable(table) {
			return false
		}
	}
	return true
}

func campusSchemaTables() []interface{} {
	return append([]interface{}{}, campusSchema.Tables()...)
}

func ensureCampusBaseSeed(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, seed := range defaultCampusCategorySeeds {
			name := strings.TrimSpace(seed.Name)
			if name == "" {
				continue
			}

			var category campusSchema.CampusCategory
			err := tx.Table(category.TableName()).Where("name = ? AND parent_id IS NULL", name).First(&category).Error
			switch {
			case err == nil:
				updates := map[string]interface{}{
					"status":     0,
					"sort_order": seed.SortOrder,
				}
				if updateErr := tx.Table(category.TableName()).Where("id = ?", category.ID).Updates(updates).Error; updateErr != nil {
					return updateErr
				}
			case errors.Is(err, gorm.ErrRecordNotFound):
				category = campusSchema.CampusCategory{
					Name:      name,
					SortOrder: seed.SortOrder,
					Status:    0,
				}
				if createErr := tx.Table(category.TableName()).Create(&category).Error; createErr != nil {
					return createErr
				}
			default:
				return err
			}
		}
		return nil
	})
}
