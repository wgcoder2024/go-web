package models

type SQLQuery struct {
	SQL string `json:"sql" binding:"required"`
}

type TableData struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
	Total   int64           `json:"total"`
}

type TableBackup struct {
	Name      string `json:"name"`
	Structure string `json:"structure"`
	Data      string `json:"data"`
}
