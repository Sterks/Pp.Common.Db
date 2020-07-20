package models

// Seq Получение последнего id
type Seq struct {
	Last_value int `gorm:"column:last_value"`
	Log_cnt int `gorm:"column:log_cnt"`
	Is_called bool `gorm:"column:is_called"`
}