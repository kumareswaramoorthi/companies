package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/kumareswaramoorthi/companies/api/constants"
	"github.com/kumareswaramoorthi/companies/api/logging"
	"github.com/kumareswaramoorthi/companies/api/models"
)

type Repository interface {
	CreateCompany(c *gin.Context, company models.Company) error
	GetCompany(c *gin.Context, id string) (models.Company, error)
	DeleteCompany(c *gin.Context, id string) error
	CheckCompanyExistsByName(c *gin.Context, name string) (bool, error)
	CheckCompanyExistsByID(c *gin.Context, id string) (bool, error)
	UpdateCompany(c *gin.Context, updateFields map[string]interface{}, id string) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return repository{db: db}
}

const (
	insertCompany            = `INSERT INTO companies (id,name,description,amount_of_employees,registered,type) VALUES ($1,$2,$3,$4,$5,$6)`
	getCompany               = `SELECT * FROM companies WHERE id  = $1`
	checkCompanyExistsByName = `SELECT EXISTS(SELECT 1 FROM companies where name = $1)`
	checkCompanyExistsByID   = `SELECT EXISTS(SELECT 1 FROM companies where id = $1)`
	deleteCompany            = `DELETE  FROM companies WHERE id  = $1`
)

func (r repository) CreateCompany(c *gin.Context, company models.Company) error {

	logger := logging.GetLogger(c).
		WithField(constants.ReqID, requestid.Get(c)).
		WithField(constants.Interface, "Repository").
		WithField(constants.Method, "CreateCompany")

	_, err := r.db.ExecContext(c.Request.Context(), insertCompany, company.ID, company.Name, company.Description, company.AmountOfEmployees, company.Registered, company.Type)
	if err != nil {
		logger.Errorf("repository: CreateCompany ID [%s]", err.Error())
		return err
	}

	logger.Debugf("created company with ID: [%s]", company.ID)
	return nil
}

func (r repository) GetCompany(c *gin.Context, id string) (models.Company, error) {

	logger := logging.GetLogger(c).
		WithField(constants.ReqID, requestid.Get(c)).
		WithField(constants.Interface, "Repository").
		WithField(constants.Method, "GetCompany")

	var company models.Company
	err := r.db.GetContext(c.Request.Context(), &company, getCompany, id)

	switch {
	case err == sql.ErrNoRows:
		logger.Errorf("no rows found for ID: [%s]", id)
		return models.Company{}, err
	case err != nil:
		logger.Errorf("repository: GetCompany ID [%s]", err.Error())
		return models.Company{}, err
	}

	logger.Debugf("found company for ID: [%s]", id)
	return company, nil
}

func (r repository) DeleteCompany(c *gin.Context, id string) error {

	logger := logging.GetLogger(c).
		WithField(constants.ReqID, requestid.Get(c)).
		WithField(constants.Interface, "Repository").
		WithField(constants.Method, "DeleteCompany")

	result, err := r.db.ExecContext(c.Request.Context(), deleteCompany, id)
	if err != nil {
		logger.Errorf("repository: DeleteCompany ID [%s] error: %s", id, err.Error())
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no company found for ID [%s]", id)
	}

	logger.Debugf("deleted company with ID %s", id)
	return nil
}

func (r repository) CheckCompanyExistsByName(c *gin.Context, name string) (bool, error) {

	logger := logging.GetLogger(c).
		WithField(constants.ReqID, requestid.Get(c)).
		WithField(constants.Interface, "Repository").
		WithField(constants.Method, "CheckCompanyExistsByName")

	var exists bool
	err := r.db.GetContext(c.Request.Context(), &exists, checkCompanyExistsByName, name)
	if err != nil {
		logger.Errorf("repository: CheckCompanyExistsByName name [%s] error: %s", name, err.Error())
		return false, err
	}

	logger.Debugf("company exists with name: [%s]", name)
	return exists, nil
}

func (r repository) CheckCompanyExistsByID(c *gin.Context, id string) (bool, error) {

	logger := logging.GetLogger(c).
		WithField(constants.ReqID, requestid.Get(c)).
		WithField(constants.Interface, "Repository").
		WithField(constants.Method, "CheckCompanyExistsByID")

	var exists bool
	err := r.db.GetContext(c.Request.Context(), &exists, checkCompanyExistsByID, id)
	if err != nil {
		logger.Errorf("repository: CheckCompanyExistsByID ID [%s] error: %s", id, err.Error())
		return false, err
	}

	logger.Debugf("company exists with ID: [%s]", id)
	return exists, nil
}

func (r repository) UpdateCompany(c *gin.Context, updateFields map[string]interface{}, id string) error {

	logger := logging.GetLogger(c).
		WithField(constants.ReqID, requestid.Get(c)).
		WithField(constants.Interface, "Repository").
		WithField(constants.Method, "PatchCompany")

	sql, args := buildUpdateSql(c, id, updateFields)
	_, err := r.db.DB.Exec(sql, args...)
	if err != nil {
		logger.Errorf("repository: PatchCompany ID [%s] error: %s", id, err.Error())
		return err
	}

	logger.Debugf("updated company with ID: [%s]", id)
	return nil
}

func buildUpdateSql(c *gin.Context, id string, updateFields map[string]interface{}) (string, []interface{}) {
	var (
		setValues   []string
		args        []interface{}
		fieldsCount int = 1
	)

	for field, value := range updateFields {
		setValues = append(setValues, fmt.Sprintf(` %s = $%d `, field, fieldsCount))
		args = append(args, value)
		fieldsCount++
	}
	args = append(args, id)
	setClause := strings.Join(setValues, ", ")

	return fmt.Sprintf(`UPDATE companies SET %s  WHERE id = $%d `, setClause, fieldsCount), args
}
