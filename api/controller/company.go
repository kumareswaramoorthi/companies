package controller

import (
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/kumareswaramoorthi/companies/api/constants"
	"github.com/kumareswaramoorthi/companies/api/errors"
	"github.com/kumareswaramoorthi/companies/api/logging"
	"github.com/kumareswaramoorthi/companies/api/models"
	service "github.com/kumareswaramoorthi/companies/api/service"
	"github.com/kumareswaramoorthi/companies/api/utils"
)

type Controller interface {
	CreateCompany(c *gin.Context)
	GetCompany(c *gin.Context)
	DeleteCompany(c *gin.Context)
	UpdateCompany(c *gin.Context)
}

type controller struct {
	svc service.Company
}

func NewController(svc service.Company) Controller {
	return &controller{svc: svc}
}

// Company godoc
// @Tags Company
// @Summary create company
// @Description creation of new company
// @Accept json
// @Produce  json
// @Success 201 {object} models.Company
// @Failure 400 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Param CreateCompany body models.Company true "request body"
// @param authorization header string true "string" default(authorization)
// @Router /api/v1/company [POST]
func (ctrl controller) CreateCompany(c *gin.Context) {
	logger := logging.GetLogger(c).
		WithField(constants.ReqID, requestid.Get(c)).
		WithField(constants.Interface, "Controller").
		WithField(constants.Method, "CreateCompany")

	companyReq := models.Company{}

	if err := c.ShouldBindJSON(&companyReq); err != nil {
		c.AbortWithStatusJSON(errors.ErrBadRequest.HttpStatusCode, errors.ErrBadRequest)
		return
	}
	_, validationerr := govalidator.ValidateStruct(companyReq)
	if validationerr != nil {
		logger.Errorf("CreateCompany - %s", validationerr.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Validation Failed "+validationerr.Error())
		return
	}
	company, err := ctrl.svc.CreateCompany(c, companyReq)
	if err != nil {
		logger.Errorf("CreateCompany - %s", err.Error())
		c.AbortWithStatusJSON(err.HttpStatusCode, err)
		return
	}

	c.JSON(http.StatusCreated, company)
}

// Company godoc
// @Tags Company
// @Summary get company
// @Description get company info by ID
// @Accept json
// @Produce  json
// @Success 200 {object} models.Company
// @Failure 400 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/v1/company/:id [GET]
func (ctrl controller) GetCompany(c *gin.Context) {
	logger := logging.GetLogger(c).
		WithField(constants.ReqID, requestid.Get(c)).
		WithField(constants.Interface, "Controller").
		WithField(constants.Method, "GetCompany")

	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(errors.ErrBadRequest.HttpStatusCode, errors.ErrBadRequest)
		return
	}

	company, err := ctrl.svc.GetCompany(c, id)
	if err != nil {
		logger.Errorf("GetCompany - %s", err.Error())
		c.AbortWithStatusJSON(err.HttpStatusCode, err)
		return
	}

	c.JSON(http.StatusOK, company)
}

// Company godoc
// @Tags Company
// @Summary delete a company
// @Description delete company by ID
// @Accept json
// @Produce  json
// @Success 200 {string} successfully deleted company
// @Failure 400 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @param authorization header string true "string" default(authorization)
// @Router /api/v1/company/:id [DELETE]
func (ctrl controller) DeleteCompany(c *gin.Context) {
	logger := logging.GetLogger(c).
		WithField(constants.ReqID, requestid.Get(c)).
		WithField(constants.Interface, "Controller").
		WithField(constants.Method, "DeleteCompany")

	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(errors.ErrBadRequest.HttpStatusCode, errors.ErrBadRequest)
		return
	}

	err := ctrl.svc.DeleteCompany(c, id)
	if err != nil {
		logger.Errorf("DeleteCompany - %s", err.Error())
		c.AbortWithStatusJSON(err.HttpStatusCode, err)
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("successfully deleted company with id: %s", id))
}

// Company godoc
// @Tags Company
// @Summary update a company
// @Description update company by ID
// @Accept json
// @Produce  json
// @Success 200 {object} models.Company
// @Failure 400 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Param updateReq body models.Company true "request body"
// @param authorization header string true "string" default(authorization)
// @Router /api/v1/company/:id [PATCH]
func (ctrl controller) UpdateCompany(c *gin.Context) {
	logger := logging.GetLogger(c).
		WithField(constants.ReqID, requestid.Get(c)).
		WithField(constants.Interface, "Controller").
		WithField(constants.Method, "UpdateCompany")

	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(errors.ErrBadRequest.HttpStatusCode, errors.ErrBadRequest)
		return
	}

	var updateReq map[string]interface{}
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.AbortWithStatusJSON(errors.ErrBadRequest.HttpStatusCode, errors.ErrBadRequest)
		return
	}
	_, validationerr := govalidator.ValidateMap(updateReq, utils.GetMapValidations())
	if validationerr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Validation Failed "+validationerr.Error())
		return
	}

	company, err := ctrl.svc.UpdateCompany(c, id, updateReq)
	if err != nil {
		logger.Errorf("UpdateCompany - %s", err.Error())
		c.AbortWithStatusJSON(err.HttpStatusCode, err)
		return
	}

	c.JSON(http.StatusOK, company)
}
