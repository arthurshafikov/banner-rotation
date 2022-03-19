package tests

import (
	"context"
	"fmt"
	httppkg "net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/thewolf27/banner-rotation/internal/config"
	"github.com/thewolf27/banner-rotation/internal/core"
	"github.com/thewolf27/banner-rotation/internal/repository"
	"github.com/thewolf27/banner-rotation/internal/services"
	"github.com/thewolf27/banner-rotation/internal/transport/http"
	"github.com/thewolf27/banner-rotation/internal/transport/http/handler"
	"github.com/thewolf27/banner-rotation/pkg/postgres"
)

var (
	r *require.Assertions
)

type APITestSuite struct {
	suite.Suite

	db *sqlx.DB

	server  *http.Server
	handler *handler.Handler
	repos   *repository.Repository
	config  *config.Config

	ctx       context.Context
	ctxCancel context.CancelFunc
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	r = s.Require()
	s.ctx, s.ctxCancel = context.WithCancel(context.Background())

	s.config = &config.Config{
		DatabaseConfig: config.DatabaseConfig{
			DSN: "host=db user=homestead password=secret dbname=homestead sslmode=disable",
		},
		MultihandedBanditConfig: config.MultihandedBanditConfig{
			EGreedValue: 0.1,
		},
		ServerConfig: config.ServerConfig{
			Port: "8999",
		},
	}

	s.db = postgres.NewSqlxDb(s.ctx, s.config.DSN)
	s.repos = repository.NewRepository(s.db)
	services := services.NewServices(services.Dependencies{
		Repository:  s.repos,
		EGreedValue: s.config.MultihandedBanditConfig.EGreedValue,
	})

	s.handler = handler.NewHandler(s.ctx, services, http.NewRequestParser())
	s.server = http.NewServer(s.ctx, s.handler)
	go func() {
		s.server.Serve(s.config.ServerConfig.Port)
	}()
}

func (s *APITestSuite) TearDownTest() {
	r.NoError(s.resetDatabase())
}

func (s *APITestSuite) TearDownSuite() {
	s.ctxCancel()
}

func (s *APITestSuite) resetDatabase() error {
	tables := []string{
		core.BannersTable,
		core.SlotsTable,
		core.BannerSlotTable,
		core.SocialGroupTable,
		core.BannerSlotSocialGroupTable,
	}
	_, err := s.db.Exec(fmt.Sprintf(`TRUNCATE TABLE %s`, strings.Join(tables, ", ")))

	return err
}

func (s *APITestSuite) makeGetRequest(path string, urlValues url.Values) (*httppkg.Response, error) {
	return s.makeRequest(httppkg.MethodGet, path, urlValues)
}

func (s *APITestSuite) makePostRequest(path string, urlValues url.Values) (*httppkg.Response, error) {
	return s.makeRequest(httppkg.MethodPost, path, urlValues)
}

func (s *APITestSuite) makeDeleteRequest(path string, urlValues url.Values) (*httppkg.Response, error) {
	return s.makeRequest(httppkg.MethodDelete, path, urlValues)
}

func (s *APITestSuite) makeRequest(method string, path string, urlValues url.Values) (*httppkg.Response, error) {
	req, err := httppkg.NewRequest(
		method,
		fmt.Sprintf("http://integration:%v/%s?%s", s.config.Port, path, urlValues.Encode()),
		nil,
	)
	if err != nil {
		return nil, err
	}
	client := &httppkg.Client{}

	return client.Do(req)
}
