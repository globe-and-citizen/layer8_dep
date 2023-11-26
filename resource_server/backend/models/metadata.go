package models

type UserMetadata struct {
	ID        uint   `gorm:"primaryKey; unique; autoIncrement; not null" json:"id"`
	UserID    uint   `gorm:"column:user_id; not null" json:"user_id"`
	Key       string `gorm:"column:key; not null" json:"key"`
	Value     string `gorm:"column:value; not null" json:"value"`
	CreatedAt string `gorm:"column:created_at; not null" json:"created_at"`
	UpdatedAt string `gorm:"column:updated_at; not null" json:"updated_at"`
}

func (UserMetadata) TableName() string {
	return "user_metadata"
}
