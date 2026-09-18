package v1_test

import (
	"testing"

	"github.com/bagus-aulia/inventhier/config"
	http_mocks "github.com/bagus-aulia/inventhier/internal/adapters/helpers/http_client/mocks"
	"github.com/stretchr/testify/suite"
)

type suiteAPI struct {
	suite.Suite
	client *http_mocks.HttpClient
	cfg    *config.Config
}

func (s *suiteAPI) SetupSuite() {
	s.cfg = &config.Config{
		ClientTimeout:        10,
		PaymentServiceURL:    "http://abc.com",
		PaymentServiceAPIKey: "",
	}
}

func (s *suiteAPI) SetupTest() {
	s.client = &http_mocks.HttpClient{}
}

func (s *suiteAPI) TearDownTest() {
	s.client = nil
}

func (s *suiteAPI) TearDownSuite() {
	s.cfg = nil
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(suiteAPI))
}
