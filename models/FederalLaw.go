package models

// FederalLaw ...
type FederalLaw struct {
	RLID         int       `gorm:"column:fl_id;primary_key"`
	FLNameLaw    string    `gorm:"column:fl_name_law"`
	FLComment    string    `gorm:"column:fl_comment"`
}
