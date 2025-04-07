package bootstrap

import (
	"fmt"
	"forum/g"
	"forum/pkg/zerolog2gorm"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化默认db
func InitDB() *gorm.DB {
	db := initMysql(
		g.Conf.GetString("Database.Default.Username"),
		g.Conf.GetString("Database.Default.Password"),
		g.Conf.GetString("Database.Default.Host"),
		g.Conf.GetInt("Database.Default.Port"),
		g.Conf.GetString("Database.Default.Dbname"),
	)

	// if g.Conf.GetBool("Debug") && strings.ToLower(g.Conf.GetString("Env")) == "development" {
	// 	// 迁移数据
	// 	err := db.AutoMigrate(
	// 	)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	return db
}

// initMysql 初始化mysql
func initMysql(username, password, host string, port int, dbname string) *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		dbname,
	)
	l := zerolog2gorm.New(&log.Logger)
	l.SetAsDefault()
	if g.Conf.GetBool("Debug") {
		// 调试模式,所有sql都打印
		l.LogMode(logger.Info)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   l,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(fmt.Errorf("数据库连接失败:%w", err))
	}
	if g.Conf.GetBool("Debug") {
		db = db.Debug()
	}
	return db
}
