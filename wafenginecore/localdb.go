package wafenginecore

import (
	"SamWaf/global"
	"SamWaf/innerbean"
	"SamWaf/model"
	"SamWaf/utils"
	"SamWaf/utils/zlog"
	"fmt"
	"net/url"
	//"github.com/kangarooxin/gorm-plugin-crypto"
	//"github.com/kangarooxin/gorm-plugin-crypto/strategy"
	"github.com/pengge/sqlitedriver"
	"gorm.io/gorm"
)

func InitCoreDb(currentDir string) {
	if currentDir == "" {
		currentDir = utils.GetCurrentDir()
	}
	if global.GWAF_LOCAL_DB == nil {
		path := currentDir + "/data/local.db"
		key := url.QueryEscape(global.GWAF_PWD_COREDB)
		dns := fmt.Sprintf("%s?_db_key=%s", path, key)
		db, err := gorm.Open(sqlite.Open(dns), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		// 启用 WAL 模式
		_ = db.Exec("PRAGMA journal_mode=WAL;")
		global.GWAF_LOCAL_DB = db
		db.DB()
		//db.Use(crypto.NewCryptoPlugin())
		// 注册默认的AES加解密策略
		//crypto.RegisterCryptoStrategy(strategy.NewAesCryptoStrategy("3Y)(27EtO^tK8Bj~"))
		// Migrate the schema
		db.AutoMigrate(&model.Hosts{})
		db.AutoMigrate(&model.Rules{})

		//隐私处理
		db.AutoMigrate(&model.LDPUrl{})

		//白名单处理
		db.AutoMigrate(&model.IPWhiteList{})
		db.AutoMigrate(&model.URLWhiteList{})

		//限制处理
		db.AutoMigrate(&model.IPBlockList{})
		db.AutoMigrate(&model.URLBlockList{})

		//抵抗CC
		db.AutoMigrate(&model.AntiCC{})

		//waf自身账号
		db.AutoMigrate(&model.TokenInfo{})
		db.AutoMigrate(&model.Account{})

		//系统参数
		db.AutoMigrate(&model.SystemConfig{})

		//延迟信息
		db.AutoMigrate(&model.DelayMsg{})

		global.GWAF_LOCAL_DB.Callback().Query().Before("gorm:query").Register("tenant_plugin:before_query", before_query)
		global.GWAF_LOCAL_DB.Callback().Query().Before("gorm:update").Register("tenant_plugin:before_update", before_update)

		//重启需要删除无效规则
		db.Where("user_code = ? and rule_status = 999", global.GWAF_USER_CODE).Delete(model.Rules{})

	}
}

func InitLogDb(currentDir string) {
	if currentDir == "" {
		currentDir = utils.GetCurrentDir()
	}
	if global.GWAF_LOCAL_LOG_DB == nil {
		path := currentDir + "/data/local_log.db"
		key := url.QueryEscape(global.GWAF_PWD_LOGDB)
		dns := fmt.Sprintf("%s?_db_key=%s", path, key)
		db, err := gorm.Open(sqlite.Open(dns), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		// 启用 WAL 模式
		_ = db.Exec("PRAGMA journal_mode=WAL;")
		global.GWAF_LOCAL_LOG_DB = db
		//logDB.Use(crypto.NewCryptoPlugin())
		// 注册默认的AES加解密策略
		//crypto.RegisterCryptoStrategy(strategy.NewAesCryptoStrategy("3Y)(27EtO^tK8Bj~"))
		// Migrate the schema
		//统计处理
		db.AutoMigrate(&innerbean.WebLog{})
		db.AutoMigrate(&model.AccountLog{})
		db.AutoMigrate(&model.WafSysLog{})
		global.GWAF_LOCAL_LOG_DB.Callback().Query().Before("gorm:query").Register("tenant_plugin:before_query", before_query)
		global.GWAF_LOCAL_LOG_DB.Callback().Query().Before("gorm:update").Register("tenant_plugin:before_update", before_update)

	}
}

func InitStatsDb(currentDir string) {
	if currentDir == "" {
		currentDir = utils.GetCurrentDir()
	}
	if global.GWAF_LOCAL_STATS_DB == nil {
		path := currentDir + "/data/local_stats.db"
		key := url.QueryEscape(global.GWAF_PWD_STATDB)
		dns := fmt.Sprintf("%s?_db_key=%s", path, key)
		db, err := gorm.Open(sqlite.Open(dns), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		// 启用 WAL 模式
		_ = db.Exec("PRAGMA journal_mode=WAL;")

		global.GWAF_LOCAL_STATS_DB = db
		//db.Use(crypto.NewCryptoPlugin())
		// 注册默认的AES加解密策略
		//crypto.RegisterCryptoStrategy(strategy.NewAesCryptoStrategy("3Y)(27EtO^tK8Bj~"))
		// Migrate the schema
		//统计处理
		db.AutoMigrate(&model.StatsTotal{})
		db.AutoMigrate(&model.StatsDay{})
		db.AutoMigrate(&model.StatsIPDay{})
		db.AutoMigrate(&model.StatsIPCityDay{})
		global.GWAF_LOCAL_STATS_DB.Callback().Query().Before("gorm:query").Register("tenant_plugin:before_query", before_query)
		global.GWAF_LOCAL_STATS_DB.Callback().Query().Before("gorm:update").Register("tenant_plugin:before_update", before_update)

	}
}

func before_query(db *gorm.DB) {
	if global.GWAF_RELEASE == "false" {
		db.Debug()
	}
	db.Where("tenant_id = ? and user_code=? ", global.GWAF_TENANT_ID, global.GWAF_USER_CODE)
	zlog.Debug("before_query")
}
func before_update(db *gorm.DB) {
	if global.GWAF_RELEASE == "false" {
		db.Debug()
	}
}
