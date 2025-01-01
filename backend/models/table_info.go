package models

type TableInfo struct {
	Name        string       `json:"name"`
	ColumnCount int          `json:"columns"`
	Rows        int          `json:"rows"`
	CreateTime  string       `json:"createTime"`
	Columns     []ColumnInfo `json:"columnDetails,omitempty"`
}

type ColumnInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
	Default  string `json:"default"`
	Key      string `json:"key"`
}
