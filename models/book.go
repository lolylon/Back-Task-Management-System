package models

import "time"

type Book struct {
    ID         uint      `json:"id" gorm:"primaryKey"`
    Title      string    `json:"title" gorm:"not null"`
    AuthorID   uint      `json:"author_id" gorm:"not null"`
    CategoryID uint      `json:"category_id" gorm:"not null"`
    Price      float64   `json:"price" gorm:"not null"`
    ISBN       string    `json:"isbn" gorm:"unique"`
    Stock      int       `json:"stock" gorm:"default:0"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    
    Author   Author   `json:"author" gorm:"foreignKey:AuthorID"`
    Category Category `json:"category" gorm:"foreignKey:CategoryID"`
}
