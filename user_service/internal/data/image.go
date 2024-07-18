package data

import (
	"time"

	"gorm.io/gorm"
)

type Image struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name     string
	Size     int64
	Location string
	UserID   uint
	User     User `gorm:"constraint:OnDelete:CASCADE;"`
}
type ImageModel struct {
	DB *gorm.DB
}

func (i ImageModel) GetProfilePicture(userid uint) (*Image, error) {
	var image Image
	t := i.DB.Where("user_id = ?", userid).Find(&image)
	return &image, t.Error
}

func (i ImageModel) UpdateProfilePicture(image *Image) error {
	t := i.DB.Save(image)
	return t.Error
}
