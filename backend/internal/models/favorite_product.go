package models

type FavoriteProduct struct {
	UserID    uint `gorm:"primaryKey"`
	User      User
	ProductID uint `gorm:"primaryKey"`
	Product   Product
}
