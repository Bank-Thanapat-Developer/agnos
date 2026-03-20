package model

type RefGender struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Code string `gorm:"type:varchar(1);not null" json:"code"`
	Name string `gorm:"type:varchar(255);not null" json:"name"`
}

func (RefGender) TableName() string {
	return "ref_gender"
}
