package models

import "time"

type Author struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"not null"`
    Email     string    `json:"email" gorm:"unique"`
    Bio       string    `json:"bio"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    Books []Book `json:"books" gorm:"foreignKey:AuthorID"`
}
