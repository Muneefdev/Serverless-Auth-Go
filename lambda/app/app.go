package app

import (
	"lambda-func/api"
	"lambda-func/database"
)

type App struct {
	ApiHanlder api.ApiHandler
}

func NewApp() *App {
	db := database.NewDynamoDBClient()
	apiHanlder := api.NewApiHandler(db)

	return &App{
		ApiHanlder: *apiHanlder,
	}
}
