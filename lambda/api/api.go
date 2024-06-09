package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) *ApiHandler {
	return &ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registeruser types.User

	err := json.Unmarshal([]byte(request.Body), &registeruser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Bad request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if registeruser.Username == "" || registeruser.Password == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request - Empty fields",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	// Check if user exists
	userExists, err := api.dbStore.DoesUserExists(registeruser.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error - Could not check user",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if userExists {
		return events.APIGatewayProxyResponse{
			Body:       "User already exists",
			StatusCode: http.StatusNotFound,
		}, err
	}

	user, err := types.NewUser(registeruser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error - Could not create user",
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	err = api.dbStore.InsertUser(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Could not create user",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("Could not insert user %w", err)
	}

	return events.APIGatewayProxyResponse{
		Body:       "User created",
		StatusCode: http.StatusCreated,
	}, nil
}

func (api ApiHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var loginUser types.User
	err := json.Unmarshal([]byte(request.Body), &loginUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Bad request",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	var user *types.User
	user, err = api.dbStore.GetUserByName(loginUser.Username)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	if user == nil {
		return events.APIGatewayProxyResponse{
			Body:       "User not found",
			StatusCode: http.StatusNotFound,
		}, nil
	}

	match := types.ValidatePassword(user.Password, loginUser.Password)
	if !match {
		return events.APIGatewayProxyResponse{
			Body:       "Not Authorized",
			StatusCode: http.StatusUnauthorized,
		}, nil
	}

  //TODO: Generate jwt token

	return events.APIGatewayProxyResponse{
		Body: "Success",
		Headers: map[string]string{
			"Authorization": "Beraer ....",
		},
		StatusCode: http.StatusOK,
	}, nil
}
