package database

import (
	"fmt"
	"go-web-template/internal/util"
	"os"
	"path/filepath"

	"github.com/tenz-io/gokit/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"go-web-template/internal/config"
	"go-web-template/internal/constant"
	"go-web-template/internal/model"
)

const (
	initAdminName     = "admin"            // 默认用户名
	intiAdminPassword = "admin"            // 默认密码，生产环境建议修改
	initAdminSalt     = "1a2b3c4d5e6f7g8h" // 固定盐值，生产环境建议使用随机盐值
)

// DB GORM 数据库连接
type DB struct {
	conn *gorm.DB
}

// NewDB 创建 GORM 数据库连接
func NewDB(cfg *config.DBConfig) (*DB, error) {
	// 确保数据库目录存在
	dbDir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// 配置 GORM 日志
	gormLog := gormLogger.Default.LogMode(gormLogger.Silent)
	if cfg.Debug {
		gormLog = gormLogger.Default.LogMode(gormLogger.Info)
	}

	// 连接 SQLite 数据库
	conn, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{
		Logger: gormLog,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 获取底层 sql.DB 对象进行连接池配置
	_, err = conn.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	db := &DB{conn: conn}

	// 初始化数据库表
	if err := db.initTables(); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	// 初始化默认管理员账户
	if err := db.initDefaultAdmin(); err != nil {
		return nil, fmt.Errorf("failed to initialize default admin: %w", err)
	}

	logger.Info("GORM SQLite database initialized successfully")
	return db, nil
}

// GetConn 获取 GORM 连接
func (db *DB) GetConn() *gorm.DB {
	return db.conn
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// initTables 初始化数据库表
func (db *DB) initTables() error {
	// 自动迁移用户表
	return db.conn.AutoMigrate(&model.User{})
}

// initDefaultAdmin 初始化默认管理员账户
func (db *DB) initDefaultAdmin() error {
	// 检查是否已存在管理员账户
	var count int64
	err := db.conn.Model(&model.User{}).Where("role = ?", constant.RoleAdmin).Count(&count).Error
	if err != nil {
		return err
	}

	// 如果已存在管理员账户，跳过初始化
	if count > 0 {
		logger.Info("Admin user already exists, skipping initialization")
		return nil
	}

	initAdminPassHash := util.HashPasswordWithSalt(intiAdminPassword, initAdminSalt)

	// 创建默认管理员账户
	admin := &model.User{
		Username: initAdminName,
		Password: initAdminPassHash,
		Salt:     initAdminSalt,
		Role:     int32(constant.RoleAdmin),
		Email:    "admin@example.com",
		Profile:  "系统管理员",
	}

	if err := db.conn.Create(admin).Error; err != nil {
		return err
	}

	logger.Info("Default admin user created: admin/admin")
	return nil
}
