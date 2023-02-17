package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	er "github.com/kumareswaramoorthi/companies/api/errors"
	"github.com/kumareswaramoorthi/companies/api/models"
	"github.com/kumareswaramoorthi/companies/api/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type CompanyServiceTestSuite struct {
	suite.Suite
	mockCtrl              *gomock.Controller
	mockCompanyRepository *mocks.MockRepository
	CompanyService        Company
	context               *gin.Context
}

func TestCompanyService(t *testing.T) {
	suite.Run(t, new(CompanyServiceTestSuite))
}

func (suite *CompanyServiceTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockCompanyRepository = mocks.NewMockRepository(suite.mockCtrl)
	suite.CompanyService = NewService(suite.mockCompanyRepository)
	suite.context, _ = gin.CreateTestContext(httptest.NewRecorder())
	suite.context.Request, _ = http.NewRequest("GET", "", nil)

}

var id string = "041d2027-e6fa-4d6d-836d-eedb235c82bc"

func (suite *CompanyServiceTestSuite) TestGetCompanySuccess() {
	var expectedCompany models.Company = models.Company{
		ID:                id,
		Name:              "xyz",
		Description:       "test company",
		AmountOfEmployees: 100,
		Registered:        true,
		Type:              "Corporations"}

	suite.mockCompanyRepository.EXPECT().GetCompany(suite.context, id).Return(expectedCompany, nil)
	company, err := suite.CompanyService.GetCompany(suite.context, id)
	suite.Nil(err)
	suite.Equal(expectedCompany, company)
}

func (suite *CompanyServiceTestSuite) TestGetCompanyFail() {
	suite.mockCompanyRepository.EXPECT().GetCompany(suite.context, id).Return(models.Company{}, errors.New("something went wrong"))
	_, err := suite.CompanyService.GetCompany(suite.context, id)
	suite.NotNil(err)
}

func (suite *CompanyServiceTestSuite) TestDeleteCompanySuccess() {
	suite.mockCompanyRepository.EXPECT().CheckCompanyExistsByID(suite.context, id).Return(true, nil)
	suite.mockCompanyRepository.EXPECT().DeleteCompany(suite.context, id).Return(nil)
	err := suite.CompanyService.DeleteCompany(suite.context, id)
	suite.Nil(err)
}

func (suite *CompanyServiceTestSuite) TestDeleteCompanyFailIfIDNonExists() {
	suite.mockCompanyRepository.EXPECT().CheckCompanyExistsByID(suite.context, id).Return(false, nil)
	err := suite.CompanyService.DeleteCompany(suite.context, id)
	suite.NotNil(err)
	suite.Equal(err, er.ErrNoCompanyRecordsFoundByID)
}

func (suite *CompanyServiceTestSuite) TestDeleteCompanyFailIfDBErr() {
	suite.mockCompanyRepository.EXPECT().CheckCompanyExistsByID(suite.context, id).Return(true, nil)
	suite.mockCompanyRepository.EXPECT().DeleteCompany(suite.context, id).Return(er.ErrUnableToDeleteCompany)
	err := suite.CompanyService.DeleteCompany(suite.context, id)
	suite.NotNil(err)
	suite.Equal(err, er.ErrUnableToDeleteCompany)
}

func (suite *CompanyServiceTestSuite) TestDeleteCompanyFailIfCheckErr() {
	suite.mockCompanyRepository.EXPECT().CheckCompanyExistsByID(suite.context, id).Return(false, er.ErrInternalServerError)
	err := suite.CompanyService.DeleteCompany(suite.context, id)
	suite.NotNil(err)
	suite.Equal(err, er.ErrInternalServerError)
}

func (suite *CompanyServiceTestSuite) TestUpdateCompanySuccess() {
	var expectedCompany models.Company = models.Company{
		ID:                id,
		Name:              "xyz",
		Description:       "test company",
		AmountOfEmployees: 100,
		Registered:        true,
		Type:              "Corporations"}

	req := make(map[string]interface{})
	req["name"] = "xyz"
	suite.mockCompanyRepository.EXPECT().CheckCompanyExistsByID(suite.context, id).Return(true, nil)
	suite.mockCompanyRepository.EXPECT().UpdateCompany(suite.context, req, id).Return(nil)
	suite.mockCompanyRepository.EXPECT().GetCompany(suite.context, id).Return(expectedCompany, nil)
	company, err := suite.CompanyService.UpdateCompany(suite.context, id, req)
	suite.Nil(err)
	suite.Equal(expectedCompany, company)
}

func (suite *CompanyServiceTestSuite) TestUpdateCompanyFailIfDBErr() {
	req := make(map[string]interface{})
	req["name"] = "xyz"
	suite.mockCompanyRepository.EXPECT().CheckCompanyExistsByID(suite.context, id).Return(false, errors.New("something went wrong"))
	_, err := suite.CompanyService.UpdateCompany(suite.context, id, req)
	suite.NotNil(err)
	suite.Equal(err, er.ErrInternalServerError)
}

func (suite *CompanyServiceTestSuite) TestUpdateCompanyFailIfIDNonExists() {
	req := make(map[string]interface{})
	req["name"] = "xyz"
	suite.mockCompanyRepository.EXPECT().CheckCompanyExistsByID(suite.context, id).Return(false, nil)
	_, err := suite.CompanyService.UpdateCompany(suite.context, id, req)
	suite.NotNil(err)
	suite.Equal(err, er.ErrNoCompanyRecordsFoundByID)
}

func (suite *CompanyServiceTestSuite) TestUpdateCompanyFailIfDBUpdateErr() {
	req := make(map[string]interface{})
	req["name"] = "xyz"
	suite.mockCompanyRepository.EXPECT().CheckCompanyExistsByID(suite.context, id).Return(true, nil)
	suite.mockCompanyRepository.EXPECT().UpdateCompany(suite.context, req, id).Return(errors.New("something went wrong"))
	_, err := suite.CompanyService.UpdateCompany(suite.context, id, req)
	suite.NotNil(err)
	suite.Equal(err, er.ErrUnableToUpdateCompany)
}

func (suite *CompanyServiceTestSuite) TestUpdateCompanyFailIfDBFetchErr() {
	req := make(map[string]interface{})
	req["name"] = "xyz"
	suite.mockCompanyRepository.EXPECT().CheckCompanyExistsByID(suite.context, id).Return(true, nil)
	suite.mockCompanyRepository.EXPECT().UpdateCompany(suite.context, req, id).Return(nil)
	suite.mockCompanyRepository.EXPECT().GetCompany(suite.context, id).Return(models.Company{}, errors.New("something went wrong"))
	_, err := suite.CompanyService.UpdateCompany(suite.context, id, req)
	suite.NotNil(err)
	suite.Equal(err, er.ErrInternalServerError)
}

func (suite *CompanyServiceTestSuite) TestCreateCompanyFailsIfNameExists() {
	var req models.Company = models.Company{
		ID:                id,
		Name:              "xyz",
		Description:       "test company",
		AmountOfEmployees: 100,
		Registered:        true,
		Type:              "Corporations"}

	suite.mockCompanyRepository.EXPECT().CheckCompanyExistsByName(suite.context, req.Name).Return(true, nil)
	_, err := suite.CompanyService.CreateCompany(suite.context, req)
	suite.NotNil(err)
	suite.Equal(err, er.ErrRecordAlreadyExistsForGivenName)
}

func (suite *CompanyServiceTestSuite) TestCreateCompanyFailsIfNameCheckErr() {
	var req models.Company = models.Company{
		ID:                id,
		Name:              "xyz",
		Description:       "test company",
		AmountOfEmployees: 100,
		Registered:        true,
		Type:              "Corporations"}

	suite.mockCompanyRepository.EXPECT().CheckCompanyExistsByName(suite.context, req.Name).Return(true, errors.New("something went wrong"))
	_, err := suite.CompanyService.CreateCompany(suite.context, req)
	suite.NotNil(err)
	suite.Equal(err, er.ErrInternalServerError)
}
