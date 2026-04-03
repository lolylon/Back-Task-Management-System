package models

import "time"

type Category struct {
    ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    Name        string    `json:"name" gorm:"not null;unique"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
    
    Books []Book `json:"books" gorm:"foreignKey:CategoryID"`
}