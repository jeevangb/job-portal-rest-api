package services

import (
	"context"
	"errors"
	"fmt"
	"jobportalapi/internal/models"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Conn struct {

	// db is an instance of the SQLite database.
	Db *gorm.DB
}

func NewService(db *gorm.DB) (*Conn, error) {

	// We check if the database instance is nil, which would indicate an issue.
	if db == nil {
		return nil, errors.New("please provide a valid connection")
	}

	// We initialize our service with the passed database instance.
	s := &Conn{Db: db}
	return s, nil
}

func (s *Conn) CreateUser(ctx context.Context, nu models.NewUser) (models.User, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, fmt.Errorf("generating password hash: %w", err)
	}

	// We prepare the User record.
	u := models.User{
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: string(hashedPass),
	}
	err = s.Db.Create(&u).Error
	if err != nil {
		return models.User{}, err
	}

	// Successfully created the record, return the user.
	return u, nil

}

// Authenticate is a method that checks a user's provided email and password against the database.
func (s *Conn) Authenticate(ctx context.Context, email, password string) (jwt.RegisteredClaims, error) {

	// We attempt to find the User record where the email
	// matches the provided email.
	var u models.User
	tx := s.Db.Where("email = ?", email).First(&u)
	if tx.Error != nil {
		return jwt.RegisteredClaims{}, tx.Error
	}

	// We check if the provided password matches the hashed password in the database.
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return jwt.RegisteredClaims{}, err
	}

	// Successful authentication! Generate JWT claims.
	c := jwt.RegisteredClaims{
		Issuer:    "service project",
		Subject:   strconv.FormatUint(uint64(u.ID), 10),
		Audience:  jwt.ClaimStrings{"students"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	return c, nil

}
