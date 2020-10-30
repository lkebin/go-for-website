// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package app

import (
	"buhaoyong/app/api"
	"buhaoyong/app/website"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// Injectors from wire.go:

func SetupComponent() (*Component, error) {
	router := mux.NewRouter()
	component := New(router)
	return component, nil
}

func SetupAPI(c *Component) (api.Repository, error) {
	config, err := apiConfigProvider()
	if err != nil {
		return nil, err
	}
	router := c.Router
	repository := api.New(config, router)
	return repository, nil
}

func SetupWebsite(c *Component) (website.Repository, error) {
	config, err := websiteConfigProvider()
	if err != nil {
		return nil, err
	}
	router := c.Router
	repository := website.New(config, router)
	return repository, nil
}

// wire.go:

func apiConfigProvider() (*api.Config, error) {
	var c api.Config
	key := "api"
	if !viper.IsSet(key) {
		return nil, fmt.Errorf("missing %s config", key)
	}
	err := viper.UnmarshalKey(key, &c)
	if err != nil {
		return nil, fmt.Errorf("can not decode api config: %w", err)
	}
	return &c, nil
}

func websiteConfigProvider() (*website.Config, error) {
	var c website.Config
	key := "website"
	if !viper.IsSet(key) {
		return nil, fmt.Errorf("missing %s config", key)
	}
	err := viper.UnmarshalKey(key, &c)
	if err != nil {
		return nil, fmt.Errorf("can not decode website config: %w", err)
	}
	return &c, nil
}
