package services

import (
	"context"
	"errors"
	"jobportalapi/internal/models"

	"github.com/rs/zerolog/log"
)

func (s *Conn) StoreCompany(ctx context.Context, nc models.NewCompany) (models.Company, error) {

	// We prepare the Company record.
	c := models.Company{
		Name:     nc.Name,
		Location: nc.Location,
	}
	err := s.Db.Create(&c).Error
	if err != nil {
		return models.Company{}, err
	}

	// Successfully created the record, return the user.
	return c, nil
}

func (s *Conn) GetCompanyData(cid uint64) (models.Company, error) {
	var companyData models.Company

	result := s.Db.Where("id=?", cid).First(&companyData)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.Company{}, errors.New("company not found")
	}
	return companyData, nil
}

func (s *Conn) GetCompanyAllData() ([]models.Company, error) {
	var companyDetails []models.Company

	result := s.Db.Find(&companyDetails)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return nil, errors.New("company table not present")
	}
	return companyDetails, nil
}
