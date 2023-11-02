package models

import "gorm.io/gorm"

//signup validation purpose
type NewUser struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

//store to databas
type User struct {
	gorm.Model
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

// Define a new struct for login data
type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

///////////////////////////////////////////////////////////////////////////////////////////////////
//for validation purpose
type NewCompany struct {
	Name     string `json:"name" validate:"required"`
	Location string `json:"location" validate:"required"`
}

type Company struct {
	gorm.Model
	Name     string `gorm:"unique;not null" validate:"required,unique"`
	Location string `json:"location"`
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

type NewJob struct {
	Cid   uint   `json:"cid"`
	Title string `json:"title" validate:"required"`
	Desc  string `json:"desc" validate:"required"`
}

type Job struct {
	gorm.Model
	Company Company `json:"-" gorm:"ForeignKey:cid"`
	Cid     uint    `json:"cid"`
	Title   string  `json:"title"`
	Desc    string  `json:"desc"`
}
