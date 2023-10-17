package model

import "time"

type Category struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	NIK       string    `json:"nik" gorm:"type:varchar(255);not null;unique"`
	Fullname  string    `json:"fullname" gorm:"type:varchar(255);"`
	Address   string    `json:"address" gorm:"type:varchar(255);"`
	Email     string    `json:"email" gorm:"type:varchar(255);not null"`
	Password  string    `json:"-" gorm:"type:varchar(255);not null"`
	IDCard    string    `json:"id_card" gorm:"type:varchar(255);"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegister struct {
	NIK      string `json:"nik" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	IDCard   string `json:"id_card" binding:"required"`
}

type Task struct {
	ID           int    `gorm:"primaryKey" json:"id"`
	Title        string `json:"title"`
	Deadline     string `json:"deadline"`
	Priority     int    `json:"priority"`
	Status       string `json:"status"`
	CategoryID   int    `json:"category_id"`
	UserID       int    `json:"user_id"`
	DocumentPath string `json:"document_path"`
}

// type Document struct {
// 	ID       int    `gorm:"primaryKey" json:"id"`
// 	UserID   int    `json:"user_id"`
// 	FilePath string `json:"file_path"`
// }

type Session struct {
	ID     int       `gorm:"primaryKey" json:"id"`
	Token  string    `json:"token"`
	Email  string    `json:"email"`
	Expiry time.Time `json:"expiry"`
}

type TaskCategory struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
}

type UserTaskCategory struct {
	ID           int    `json:"id"`
	Fullname     string `json:"fullname"`
	Email        string `json:"email"`
	TaskID       int    `json:"task_id"`
	Task         string `json:"task"`
	Deadline     string `json:"deadline"`
	Priority     int    `json:"priority"`
	Status       string `json:"status"`
	Category     string `json:"category"`
	DocumentPath string `json:"document_path"`
}

type Credential struct {
	Host         string
	Username     string
	Password     string
	DatabaseName string
	Port         int
	Schema       string
}

type UserProfile struct {
	ID       int    `json:"id"`
	NIK      string `json:"nik"`
	Fullname string `json:"fullname"`
	Address  string `json:"address"`
	Email    string `json:"email"`
	IDCard   string `json:"id_card"`
}
