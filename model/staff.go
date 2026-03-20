package model

import "time"

type Staff struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Username   string    `gorm:"type:varchar(100);uniqueIndex:idx_username_hospital;not null" json:"username"`
	Password   string    `gorm:"type:varchar(255);not null" json:"-"`
	HospitalID uint      `gorm:"uniqueIndex:idx_username_hospital;not null" json:"hospital_id"`
	Hospital   Hospital  `gorm:"foreignKey:HospitalID" json:"hospital"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Staff) TableName() string {
	return "staff"
}
