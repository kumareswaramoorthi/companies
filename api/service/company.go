package service

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/kumareswaramoorthi/companies/api/constants"
	errors "github.com/kumareswaramoorthi/companies/api/errors"
	"github.com/kumareswaramoorthi/companies/api/logging"
	"github.com/kumareswaramoorthi/companies/api/models"
	"github.com/kumareswaramoorthi/companies/api/repository"
)

type Company interface {
	CreateCompany(c *gin.Context, company models.Company) (models.Company, *errors.ErrorResponse)
	GetCompany(c *gin.Context, id string) (models.Company, *errors.ErrorResponse)
	DeleteCompany(c *gin.Context, id string) *errors.ErrorResponse
	UpdateCompany(c *gin.Context, id string, updateReq map[string]interface{}) (models.Company, *errors.ErrorResponse)
}

type company struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Company {
	return &company{repo: repo}
}

func (s company) CreateCompany(c *gin.Context, companyReq models.Company) (models.Company, *errors.ErrorResponse) {
	logger := logging.GetLogger(c).
		WithField(constants.ReqID, requestid.Get(c)).
		WithField(constants.Interface, "Service").
		WithField(constants.Method, "CreateCompany")

	nameExists, err := s.repo.CheckCompanyExistsByName(c, companyReq.Name)
	if err != nil {
		logger.Errorf("service: CreateCompany name [%s] error: %s", companyReq.Name, err.Error())
		return models.Company{}, errors.ErrInternalServerError
	}

	if nameExists {
		return models.Company{}, errors.ErrRecordAlreadyExistsForGivenName
	}

	err = s.repo.CreateCompany(c, companyReq)
	if err != nil {
		logger.Errorf("service: CreateCompany name [%s] error: %s", companyReq.Name, err.Error())
		return models.Company{}, errors.ErrUnableToCreateCompany
	}

	company, err := s.repo.GetCompany(c, companyReq.ID)
	if err != nil {
		logger.Errorf("service: GetCompany ID [%s] error: %s", companyReq.ID, err.Error())
		return models.Company{}, errors.ErrInternalServerError
	}

	logger.Debugf("created company with ID: [%s]", company.ID)
	return company, nil
}

func (s company) GetCompany(c *gin.Context, id string) (models.Company, *errors.ErrorResponse) {
	logger := logging.GetLogger(c).
		WithField(constants.ReqID, requestid.Get(c)).
		WithField(constants.Interface, "Service").
		WithField(constants.Method, "GetCompany")

	company, err := s.repo.GetCompany(c, id)
	if err != nil {
		logger.Errorf("service: GetCompany ID [%s] error: %s", id, err.Error())
		return models.Company{}, errors.ErrUnableToFetchCompany
	}

	logger.Debugf("fetched company with ID: [%s]", id)
	return company, nil
}

func (s company) DeleteCompany(c *gin.Context, id string) *errors.ErrorResponse {
	logger := logging.GetLogger(c).
		WithField(constants.ReqID, requestid.Get(c)).
		WithField(constants.Interface, "Service").
		WithField(constants.Method, "DeleteCompany")

	exists, err := s.repo.CheckCompanyExistsByID(c, id)
	if err != nil {
		logger.Errorf("service: CreateCompany ID [%s] error: %s", id, err.Error())
		return errors.ErrInternalServerError
	}

	if !exists {
		return errors.ErrNoCompanyRecordsFoundByID
	}

	err = s.repo.DeleteCompany(c, id)
	if err != nil {
		logger.Errorf("service: DeleteCompany ID [%s] error: %s", id, err.Error())
		return errors.ErrUnableToDeleteCompany
	}

	logger.Debugf("deleted company with ID: %s", id)
	return nil
}

func (s company) UpdateCompany(c *gin.Context, id string, updateReq map[string]interface{}) (models.Company, *errors.ErrorResponse) {
	logger := logging.GetLogger(c).
		WithField(constants.ReqID, requestid.Get(c)).
		WithField(constants.Interface, "Service").
		WithField(constants.Method, "UpdateCompany")

	exists, err := s.repo.CheckCompanyExistsByID(c, id)
	if err != nil {
		logger.Errorf("service: UpdateCompany ID [%s] error: %s", id, err.Error())
		return models.Company{}, errors.ErrInternalServerError
	}

	if !exists {
		return models.Company{}, errors.ErrNoCompanyRecordsFoundByID
	}

	err = s.repo.UpdateCompany(c, updateReq, id)
	if err != nil {
		logger.Errorf("service: UpdateCompany ID [%s] error: %s", id, err.Error())
		return models.Company{}, errors.ErrUnableToUpdateCompany
	}

	company, err := s.repo.GetCompany(c, id)
	if err != nil {
		logger.Errorf("service: GetCompany ID [%s] error: %s", id, err.Error())
		return models.Company{}, errors.ErrInternalServerError
	}

	logger.Debugf("updated company with ID: [%s]", id)
	return company, nil
}
