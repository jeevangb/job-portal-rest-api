package services

import (
	"context"
	"jobportalapi/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockgen -source serviceInterface.go -destination mockmodels/serviceInterface_mock.go -package mockmodels
type Service interface {
	CreateUser(ctx context.Context, nu models.NewUser) (models.User, error)
	Authenticate(ctx context.Context, email, password string) (jwt.RegisteredClaims, error)

	StoreCompany(ctx context.Context, nc models.NewCompany) (models.Company, error)
	GetCompanyData(uid uint64) (models.Company, error)
	GetCompanyAllData() ([]models.Company, error)

	StoreJob(ctx context.Context, nc models.NewJob, cid uint64) (models.Job, error)
	GetJobData(uid uint64) (models.Job, error)
	GetAllJobData() ([]models.Job, error)
	GetJobByCompany(uid uint64) ([]models.Job, error)
}

type Store struct {
	Service
}

func NewStore(s Service) Store {
	return Store{Service: s}
}
