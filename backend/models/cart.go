package models

import "time"

type CartItem struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `gorm:"uniqueIndex:idx_user_product" json:"user_id"`
    ProductID uint      `gorm:"uniqueIndex:idx_user_product" json:"product_id"`
    Quantity  int       `gorm:"not null;default:1" json:"quantity"`
    CreatedAt time.Time `json:"created_at"`

    Product Product `gorm:"foreignKey:ProductID" json:"product"`
}