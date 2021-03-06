package main

import (
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"example.com/realworld/httpservice"
	"example.com/realworld/stor"
	"github.com/brianvoe/gofakeit"
	"github.com/labstack/echo"
	"gopkg.in/gavv/httpexpect.v2"
)

var serverURL string

func TestMain(m *testing.M) {
	gofakeit.Seed(0)
	e := echo.New()
	st, err := stor.New()
	if err != nil {
		log.Panic(err)
	}
	if err := st.Migrate(); err != nil {
		log.Panic(err)
	}
	if err := st.Clear(); err != nil {
		log.Panic(err)
	}
	s := httpservice.Service{Stor: st}
	if err := s.SetupAPI(e); err != nil {
		log.Panic(err)
	}
	srv := httptest.NewServer(e)
	serverURL = srv.URL
	resultCode := m.Run()
	srv.Close()
	os.Exit(resultCode)
}

func TestArticles(t *testing.T) {
	e := httpexpect.New(t, serverURL)
	r := e.GET("/api/articles").Expect().JSON()
	r.Path("$.articles").Array().Length().Equal(0)
}
