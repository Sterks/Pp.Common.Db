package models

import "time"

// Task - Задачи для запуска
type Task struct {
	TSID        int       `gorm:"column:ts_id;primary_key"`
	TSName      string    `gorm:"column:ts_name"`
	TSDataStart time.Time `gorm:"column:ts_data_start"`
	TSRunTimes  int       `gorm:"column:ts_run_times"`
	TSComment   string    `gorm:"column:ts_comment"`
}
