package repository

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/kumareswaramoorthi/companies/api/models"
	"github.com/stretchr/testify/suite"
)

const (
	TestGetCompany               = `SELECT * FROM companies WHERE id  = $1`
	TestInsertCompany            = `INSERT INTO companies (id,name,description,amount_of_employees,registered,type) VALUES ($1,$2,$3,$4,$5,$6)`
	TestDeleteCompany            = `DELETE  FROM companies WHERE id  = $1`
	TestcheckCompanyExistsByName = `SELECT EXISTS(SELECT 1 FROM companies where name = $1)`
	TestcheckCompanyExistsByID   = `SELECT EXISTS(SELECT 1 FROM companies where id = $1)`
)

type RepositoryTestSuite struct {
	suite.Suite
	mockCtrl   *gomock.Controller
	sqlMock    sqlmock.Sqlmock
	repository Repository
	context    *gin.Context
	recorder   *httptest.ResponseRecorder
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.context.Request, _ = http.NewRequest("GET", "", nil)
	suite.sqlMock = mock
	suite.repository = NewRepository(sqlxDB)
}

func (suite *RepositoryTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *RepositoryTestSuite) TestShouldReturnNewInstanceOfRepository() {
	db, _, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repository := NewRepository(sqlxDB)
	suite.NotNil(repository)
}

func (suite *RepositoryTestSuite) TestGetCompanySuccess() {
	id := "041d2027-e6fa-4d6d-836d-eedb235c82bc"
	rows := sqlmock.NewRows([]string{"id", "name", "description", "amount_of_employees", "registered", "type"}).
		AddRow(id, "xyz", "test company", 100, true, "Corporations")
	suite.sqlMock.ExpectQuery(regexp.QuoteMeta(TestGetCompany)).
		WillReturnRows(rows)
	if err := suite.sqlMock.ExpectationsWereMet(); err != nil {
		suite.Error(errors.New("there were unfulfilled expectations"), err)
	}
	var expectedCompany models.Company = models.Company{
		ID:                "041d2027-e6fa-4d6d-836d-eedb235c82bc",
		Name:              "xyz",
		Description:       "test company",
		AmountOfEmployees: 100,
		Registered:        true,
		Type:              "Corporations"}

	company, err := suite.repository.GetCompany(suite.context, id)
	suite.Equal(company, expectedCompany)
	suite.Nil(err)

}

func (suite *RepositoryTestSuite) TestGetCompanyShouldFailWhenDatabaseQueryFails() {
	id := "041d2027-e6fa-4d6d-836d-eedb235c82bc"
	dbErr := errors.New("ID invalid identifier")
	suite.sqlMock.ExpectQuery(regexp.QuoteMeta(TestGetCompany)).WillReturnError(dbErr)
	if err := suite.sqlMock.ExpectationsWereMet(); err != nil {
		suite.Error(errors.New("there were unfulfilled expectations"), err)
	}
	_, err := suite.repository.GetCompany(suite.context, id)

	suite.Equal(dbErr, err)
}

func (suite *RepositoryTestSuite) TestCreateCompanySuccess() {
	inputdetails := models.Company{
		ID:                "041d2027-e6fa-4d6d-836d-eedb235c82bc",
		Name:              "xyz",
		Description:       "test company",
		AmountOfEmployees: 100,
		Registered:        true,
		Type:              "Corporations"}

	suite.sqlMock.ExpectExec(regexp.QuoteMeta(TestInsertCompany)).
		WithArgs(inputdetails.ID, inputdetails.Name, inputdetails.Description, inputdetails.AmountOfEmployees, inputdetails.Registered, inputdetails.Type).WillReturnResult(sqlmock.NewResult(1, 1))
	if err := suite.sqlMock.ExpectationsWereMet(); err != nil {
		suite.Error(errors.New("there were unfulfilled expectations"), err)
	}
	err := suite.repository.CreateCompany(suite.context, inputdetails)
	suite.Nil(err)
}

func (suite *RepositoryTestSuite) TestCreateCompanyShouldFailWhenDatabaseQueryFails() {

	inputdetails := models.Company{
		ID:                "041d2027-e6fa-4d6d-836d-eedb235c82bc",
		Name:              "xyz",
		Description:       "test company",
		AmountOfEmployees: 100,
		Registered:        true,
		Type:              "Corporations"}

	dbErr := errors.New("ID invalid identifier")
	suite.sqlMock.ExpectExec(regexp.QuoteMeta(TestInsertCompany)).
		WithArgs(inputdetails.ID, inputdetails.Name, inputdetails.Description, inputdetails.AmountOfEmployees, inputdetails.Registered, inputdetails.Type).WillReturnError(dbErr)
	if err := suite.sqlMock.ExpectationsWereMet(); err != nil {
		suite.Error(errors.New("there were unfulfilled expectations"), err)
	}
	err := suite.repository.CreateCompany(suite.context, inputdetails)
	suite.Equal(dbErr, err)
}

func (suite *RepositoryTestSuite) TestDeleteCompanySuccess() {
	id := "041d2027-e6fa-4d6d-836d-eedb235c82bc"
	suite.sqlMock.ExpectExec(regexp.QuoteMeta(TestDeleteCompany)).
		WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))
	if err := suite.sqlMock.ExpectationsWereMet(); err != nil {
		suite.Error(errors.New("there were unfulfilled expectations"), err)
	}
	err := suite.repository.DeleteCompany(suite.context, id)
	suite.Nil(err)
}

func (suite *RepositoryTestSuite) TestCheckCompanyExistsByNameSuccess() {
	name := "xyz"
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1)
	suite.sqlMock.ExpectQuery(regexp.QuoteMeta(TestcheckCompanyExistsByName)).
		WillReturnRows(rows)
	if err := suite.sqlMock.ExpectationsWereMet(); err != nil {
		suite.Error(errors.New("there were unfulfilled expectations"), err)
	}

	exists, err := suite.repository.CheckCompanyExistsByName(suite.context, name)
	suite.Equal(exists, true)
	suite.Nil(err)
}

func (suite *RepositoryTestSuite) TestCheckCompanyExistsByIDSuccess() {
	id := "041d2027-e6fa-4d6d-836d-eedb235c82bc"
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1)
	suite.sqlMock.ExpectQuery(regexp.QuoteMeta(TestcheckCompanyExistsByID)).
		WillReturnRows(rows)
	if err := suite.sqlMock.ExpectationsWereMet(); err != nil {
		suite.Error(errors.New("there were unfulfilled expectations"), err)
	}

	exists, err := suite.repository.CheckCompanyExistsByID(suite.context, id)
	suite.Equal(exists, true)
	suite.Nil(err)
}

func (suite *RepositoryTestSuite) TestUpdateCompanySuccess() {
	id := "041d2027-e6fa-4d6d-836d-eedb235c82bc"
	name := "xyz"
	suite.sqlMock.ExpectExec(regexp.QuoteMeta(`UPDATE companies SET  name = $1   WHERE id = $2 `)).
		WithArgs(name, id).WillReturnResult(sqlmock.NewResult(1, 1))
	if err := suite.sqlMock.ExpectationsWereMet(); err != nil {
		suite.Error(errors.New("there were unfulfilled expectations"), err)
	}
	req := make(map[string]interface{})
	req["name"] = name
	err := suite.repository.UpdateCompany(suite.context, req, id)
	suite.Nil(err)
}
