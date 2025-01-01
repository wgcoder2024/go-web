package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wgcoder2024/go-web/backend/config"
	"github.com/wgcoder2024/go-web/backend/models"
)

// GetTables 获取所有表信息
func GetTables(c *gin.Context) {
	var tables []models.TableInfo
	db := config.DB

	// 获取所有表名
	rows, err := db.Raw(`
		SELECT 
			table_name,
			(
				SELECT COUNT(*) 
				FROM information_schema.columns 
				WHERE table_schema = DATABASE() 
				AND table_name = t.table_name
			) as columns,
			COALESCE(table_rows, 0) as rows,
			create_time
		FROM information_schema.tables t
		WHERE table_schema = DATABASE()
	`).Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var table models.TableInfo
		rows.Scan(&table.Name, &table.ColumnCount, &table.Rows, &table.CreateTime)
		tables = append(tables, table)
	}

	c.JSON(http.StatusOK, tables)
}

// GetTableDetails 获取表详细信息
func GetTableDetails(c *gin.Context) {
	tableName := c.Param("name")
	var tableInfo models.TableInfo
	tableInfo.Name = tableName

	// 获取列信息
	rows, err := config.DB.Raw(`
		SELECT 
			column_name,
			column_type,
			is_nullable = 'YES' as nullable,
			COALESCE(column_default, '') as column_default,
			column_key
		FROM information_schema.columns
		WHERE table_schema = DATABASE()
		AND table_name = ?
		ORDER BY ordinal_position
	`, tableName).Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var columns []models.ColumnInfo
	for rows.Next() {
		var col models.ColumnInfo
		rows.Scan(&col.Name, &col.Type, &col.Nullable, &col.Default, &col.Key)
		columns = append(columns, col)
	}

	tableInfo.Columns = columns
	c.JSON(http.StatusOK, tableInfo)
}

// DeleteTable 删除表
func DeleteTable(c *gin.Context) {
	tableName := c.Param("name")

	if err := config.DB.Exec("DROP TABLE IF EXISTS " + tableName).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "表删除成功"})
}

// CreateTable 创建新表
func CreateTable(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Columns []struct {
			Name     string `json:"name" binding:"required"`
			Type     string `json:"type" binding:"required"`
			Nullable bool   `json:"nullable"`
			Default  string `json:"default"`
		} `json:"columns" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构建 CREATE TABLE 语句
	sql := "CREATE TABLE " + req.Name + " (\n"
	for i, col := range req.Columns {
		sql += col.Name + " " + col.Type
		if !col.Nullable {
			sql += " NOT NULL"
		}
		if col.Default != "" {
			sql += " DEFAULT " + col.Default
		}
		if i < len(req.Columns)-1 {
			sql += ",\n"
		}
	}
	sql += "\n)"

	if err := config.DB.Exec(sql).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "表创建成功"})
}

// AlterTable 修改表结构
func AlterTable(c *gin.Context) {
	tableName := c.Param("name")
	var req struct {
		Action string `json:"action" binding:"required,oneof=add modify drop"`
		Column struct {
			Name     string `json:"name" binding:"required"`
			NewName  string `json:"newName"`
			Type     string `json:"type"`
			Nullable bool   `json:"nullable"`
			Default  string `json:"default"`
		} `json:"column"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var sql string
	switch req.Action {
	case "add":
		sql = fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s",
			tableName, req.Column.Name, req.Column.Type)
	case "modify":
		sql = fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s %s",
			tableName, req.Column.Name, req.Column.Type)
	case "drop":
		sql = fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s",
			tableName, req.Column.Name)
	}

	if err := config.DB.Exec(sql).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "表结构修改成功"})
}

// GetTableData 获取表数据
func GetTableData(c *gin.Context) {
	tableName := c.Param("name")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortField := c.DefaultQuery("sortField", "id")
	sortOrder := c.DefaultQuery("sortOrder", "asc")

	// 获取总记录数
	var total int64
	if err := config.DB.Table(tableName).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取列信息
	columns, err := config.DB.Raw(`
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_schema = DATABASE() 
		AND table_name = ?
		ORDER BY ordinal_position
	`, tableName).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer columns.Close()

	var columnNames []string
	for columns.Next() {
		var name string
		columns.Scan(&name)
		columnNames = append(columnNames, name)
	}

	// 构建查询
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(
		"SELECT * FROM %s ORDER BY %s %s LIMIT %d OFFSET %d",
		tableName, sortField, sortOrder, pageSize, offset,
	)

	rows, err := config.DB.Raw(query).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var tableData [][]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columnNames))
		valuePtrs := make([]interface{}, len(columnNames))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		tableData = append(tableData, values)
	}

	c.JSON(http.StatusOK, models.TableData{
		Columns: columnNames,
		Rows:    tableData,
		Total:   total,
	})
}

// ExecuteSQL 执行自定义 SQL
func ExecuteSQL(c *gin.Context) {
	var query models.SQLQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 只允许 SELECT 语句
	if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(query.SQL)), "SELECT") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只允许执行 SELECT 语句"})
		return
	}

	rows, err := config.DB.Raw(query.SQL).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	var result []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)

		entry := make(map[string]interface{})
		for i, col := range columns {
			entry[col] = values[i]
		}
		result = append(result, entry)
	}

	c.JSON(http.StatusOK, result)
}

// ExportTable 导出表
func ExportTable(c *gin.Context) {
	tableName := c.Param("name")

	// 获取表结构
	var createSQL string
	err := config.DB.Raw(`SHOW CREATE TABLE `+tableName).Row().Scan(&tableName, &createSQL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取表数据
	rows, err := config.DB.Raw("SELECT * FROM " + tableName).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var data []string
	columns, _ := rows.Columns()
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)

		// 构建 INSERT 语句
		data = append(data, fmt.Sprintf("INSERT INTO %s VALUES (...);", tableName))
	}

	backup := models.TableBackup{
		Name:      tableName,
		Structure: createSQL,
		Data:      strings.Join(data, "\n"),
	}

	c.JSON(http.StatusOK, backup)
}

// ImportTable 导入表
func ImportTable(c *gin.Context) {
	var backup models.TableBackup
	if err := c.ShouldBindJSON(&backup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()

	// 创建表
	if err := tx.Exec(backup.Structure).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 导入数据
	for _, insert := range strings.Split(backup.Data, "\n") {
		if err := tx.Exec(insert).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "表导入成功"})
}
