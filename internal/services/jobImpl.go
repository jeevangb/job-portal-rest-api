package services

import (
	"context"
	"errors"
	"jobportalapi/internal/models"

	"github.com/rs/zerolog/log"
)

func (s *Conn) StoreJob(ctx context.Context, nj models.NewJob, cid uint64) (models.Job, error) {

	// We prepare the Company record.
	c := models.Job{
		Title: nj.Title,
		Desc:  nj.Desc,
		Cid:   uint(cid),
	}
	err := s.Db.
		Create(&c).Error
	if err != nil {
		return models.Job{}, err
	}

	// Successfully created the record, return the user.
	return c, nil
}

// ///////////////////////////////////////////////////////////////////////////////
func (c *Conn) GetJobData(jid uint64) (models.Job, error) {
	var jobData models.Job

	result := c.Db.Where("id=?", jid).First(&jobData)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.Job{}, errors.New("company not found")
	}
	return jobData, nil

}

// ////////////////////////////////////////////////////////////////////////////////
func (s *Conn) GetAllJobData() ([]models.Job, error) {
	var jobDetails []models.Job

	result := s.Db.Find(&jobDetails)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return nil, errors.New("job table not present")
	}
	return jobDetails, nil
}

// ///////////////////////////////////////////////////////////////////////////////////////
func (s *Conn) GetJobByCompany(cjid uint64) ([]models.Job, error) {
	var jobDetails []models.Job

	result := s.Db.Where("cid=?", cjid).Find(&jobDetails)

	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return nil, errors.New("job not found")
	}
	return jobDetails, nil

}
