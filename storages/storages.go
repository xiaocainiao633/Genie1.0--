package storages

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

var database *sql.DB

// 提供了保存简单数据、用户配置等的支持。保存的数据除非应用被卸载或者被主动删除，否则会一直保留

// 初始化
func init() {
	path := filepath.Dir(os.Args[0])
	var err error
	database, err = sql.Open("sqlite3", path+"/storages.db")
	if err != nil {
		panic("failed to open database: " + err.Error())
	}
}

// Get 从指定表中获取键值
func Get(table, key string) string {
	if table == "" {
		table = "storages"
	}
	var value string
	query := `SELECT value FROM ` + table + ` WHERE key = ?`
	_ = database.QueryRow(query, key).Scan(&value)
	return value
}

// Put 写入键值对，表不存在则自动创建
func Put(table, key, value string) {
	if table == "" {
		table = "storages"
	}
	sql := `CREATE TABLE IF NOT EXISTS ` + table + ` (
		key TEXT NOT NULL PRIMARY KEY,
		value TEXT
	);`
	_, _ = database.Exec(sql)
	query := `INSERT OR REPLACE INTO ` + table + ` (key, value) VALUES (?, ?)`
	_, _ = database.Exec(query, key, value)
}

// Remove 删除指定键
func Remove(table, key string) {
	if table == "" {
		table = "storages"
	}
	query := `DELETE FROM ` + table + ` WHERE key = ?`
	_, _ = database.Exec(query, key)
}

// Contains 判断键是否存在
func Contains(table, key string) bool {
	if table == "" {
		table = "storages"
	}
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM ` + table + ` WHERE key = ?)`
	_ = database.QueryRow(query, key).Scan(&exists)
	return exists
}

// GetAll 获取所有键值对
func GetAll(table string) map[string]string {
	if table == "" {
		table = "storages"
	}
	result := make(map[string]string)
	rows, err := database.Query(`SELECT key, value FROM ` + table)
	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err == nil {
			result[key] = value
		}
	}
	return result
}

// Clear 清空指定表数据
func Clear(table string) {
	if table == "" {
		table = "storages"
	}
	_, _ = database.Exec(`DELETE FROM ` + table)
}
