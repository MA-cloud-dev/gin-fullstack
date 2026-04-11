package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize/internal"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GormMysql 初始化Mysql数据库
// Author [piexlmax](https://github.com/piexlmax)
// Author [SliverHorn](https://github.com/SliverHorn)
// Author [ByteZhou-2018](https://github.com/ByteZhou-2018)
func GormMysql() *gorm.DB {
	m := global.GVA_CONFIG.Mysql
	db, err := openMysqlDatabase(m)
	if err != nil {
		panic(err)
	}
	return db
}

// GormMysqlByConfig 通过传入配置初始化Mysql数据库
func GormMysqlByConfig(m config.Mysql) *gorm.DB {
	db, err := openMysqlDatabase(m)
	if err != nil {
		panic(err)
	}
	return db
}

// openMysqlDatabase 初始化Mysql数据库的辅助函数
func openMysqlDatabase(m config.Mysql) (*gorm.DB, error) {
	if m.Dbname == "" {
		return nil, nil
	}

	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	// 数据库配置
	general := m.GeneralDB
	if db, err := gorm.Open(mysql.New(mysqlConfig), internal.Gorm.Config(general)); err != nil {
		return nil, err
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db, nil
	}
}
