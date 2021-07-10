package server

import (
	"errors"
	"github.com/VladKornilov/fondy_payments_go/internal/database"
	"net/url"
	"os"
	"strconv"
)

type AppConfig struct {
	SiteUrl *url.URL
	MerchantId int
	MerchantPassword string
}
func NewAppConfig() (*AppConfig, error) {
	// get parameters from ENV
	rawurl, exists := os.LookupEnv("SITE_URL")
	if !exists {
		return nil, errors.New("SITE_URL was not specified")
	}
	merchID, exists := os.LookupEnv("MERCHANT_ID")
	if !exists {
		return nil, errors.New("MERCHANT_ID was not specified")
	}
	merchPass, exists := os.LookupEnv("MERCHANT_PASSWORD")
	if !exists {
		return nil, errors.New("MERCHANT_PASSWORD was not specified")
	}

	strurl, err := url.Parse(rawurl)
	if err != nil { return nil, err }
	id, err := strconv.Atoi(merchID)
	if err != nil { return nil, err }

	cfg := new(AppConfig)
	cfg.SiteUrl = strurl
	cfg.MerchantId = id
	cfg.MerchantPassword= merchPass
	return cfg, nil
}

type Application struct {
	db database.Database
	//HttpClient
	config *AppConfig
}

func CreateApplication() (*Application, error) {

	db, err := database.OpenDatabase()
	if err != nil { return nil, err }
	cfg, err := NewAppConfig()
	if err != nil { return nil, err }

	app := new(Application)
	app.db = db
	app.config = cfg
	return app, nil
}

func (a Application) Close() error {
	return a.db.Close()
}