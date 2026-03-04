package main

import (
	"learning-app-mobile-bna/Backend/initializers"
	"learning-app-mobile-bna/Backend/model"
)

func init() {
	initializers.LoadEnvs()
	initializers.ConnectDB()
}

func main() {

	initializers.DB.AutoMigrate(
		&model.Language{},
		&model.User{},
		&model.Word{},
		&model.List{},
		&model.ListWord{},
		&model.Params{})
}
