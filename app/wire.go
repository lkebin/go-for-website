// +build wireinject

package app

import (
	"buhaoyong/app/api"
	"buhaoyong/app/website"
	"buhaoyong/pkg/db"
	"buhaoyong/pkg/service/post"
	"fmt"

	"github.com/google/wire"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var (
	dbSet = wire.NewSet(db.New, dbConfigProvider)
)

func SetupComponent() (*Component, error) {
	wire.Build(
		New,
		mux.NewRouter,
		dbSet,
	)

	return nil, nil
}

func SetupAPI(c *Component) (api.Repository, error) {
	wire.Build(
		api.New,
		apiConfigProvider,
		wire.FieldsOf(&c, "Router"),
	)

	return nil, nil
}

func SetupWebsite(c *Component) (website.Repository, error) {

	wire.Build(
		website.New,
		websiteConfigProvider,
		postServiceProvider,
		wire.FieldsOf(&c, "Router"),
	)

	return nil, nil
}

func apiConfigProvider() (*api.Config, error) {
	var c api.Config
	key := "api"
	if !viper.IsSet(key) {
		return nil, fmt.Errorf("missing %s config", key)
	}

	if err := viper.UnmarshalKey(key, &c); err != nil {
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

	if err := viper.UnmarshalKey(key, &c); err != nil {
		return nil, fmt.Errorf("can not decode website config: %w", err)
	}
	return &c, nil
}

func dbConfigProvider() (*db.Config, error) {
	var c db.Config
	key := "database"
	if !viper.IsSet(key) {
		return nil, fmt.Errorf("missing %s config", key)
	}

	if err := viper.UnmarshalKey(key, &c); err != nil {
		return nil, fmt.Errorf("can not decode database config: %w", err)
	}
	return &c, nil
}

func postServiceProvider(c *Component) (post.Repository, error) {
	wire.Build(
		post.New,
		wire.FieldsOf(&c, "DB"),
	)

	return nil, nil
}
