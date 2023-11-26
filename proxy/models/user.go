package models

type User struct {
	ID               uint   `gorm:"primaryKey; unique; autoIncrement; not null" json:"id"`
	Email            string `gorm:"column:email; unique; not null" json:"email"`
	Username         string `gorm:"column:username; unique; not null" json:"username"`
	Password         string `gorm:"column:password; not null" json:"password"`
	FirstName        string `gorm:"column:first_name; not null" json:"first_name"`
	LastName         string `gorm:"column:last_name; not null" json:"last_name"`
	DisplayName      string `gorm:"column:display_name; not null" json:"display_name"`
	ShareDisplayName bool   `gorm:"column:share_display_name; not null" json:"share_display_name"`
	// PhoneNumber         string `gorm:"column:phone_number; not null" json:"phone_number"`
	// Address             string `gorm:"column:address; not null" json:"address"`
	// EmailVerified       bool   `gorm:"column:email_verified; default:false" json:"email_verified"`
	// PhoneNumberVerified bool   `gorm:"column:phone_number_verified; default:false" json:"phone_number_verified"`
	// LocationVerified    bool   `gorm:"column:location_verified; default:false" json:"location_verified"`
	// NationalIdVerified  bool   `gorm:"column:national_id_verified; default:false" json:"national_id_verified"`
	Salt string `gorm:"column:salt; not null" json:"salt"`
}

func (User) TableName() string {
	return "users"
}
