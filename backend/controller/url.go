package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/homma509/9rece/backend/domain/model"
	"github.com/homma509/9rece/backend/usecase"
)

// URLController URLコントローラのインターフェース
type URLController interface {
	Get(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type urlController struct {
	urlUsecase usecase.URLUsecase
}

type request struct {
}

// Response Response構造体
type Response struct {
	URL string `json:"url"`
}

// NewURLController URLコントローラを生成します
func NewURLController(u usecase.URLUsecase) URLController {
	return &urlController{
		urlUsecase: u,
	}
}

// Get URLの取得
func (c *urlController) Get(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	_, err := parseRequest(event)
	if err != nil {
		return response(
			http.StatusBadRequest,
			errorResponseBody(err.Error()),
		), nil
	}

	// URL生成
	url, err := c.urlUsecase.Get(ctx)
	if err != nil {
		return response(
			http.StatusInternalServerError,
			errorResponseBody(err.Error()),
		), nil
	}

	// レスポンスのJSON生成
	body, err := json.Marshal(url)
	if err != nil {
		return response(
			http.StatusInternalServerError,
			errorResponseBody(err.Error()),
		), nil
	}

	return response(http.StatusOK, string(body)), nil
}

func parseRequest(req events.APIGatewayProxyRequest) (*request, error) {
	if req.HTTPMethod != http.MethodGet {
		return nil, fmt.Errorf("use GET request")
	}

	return &request{}, nil
}

func response(code int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func responseBody(url *model.URL) (string, error) {
	resp, err := json.Marshal(Response{URL: url.URL})
	if err != nil {
		return "", err
	}

	return string(resp), nil
}

func errorResponseBody(msg string) string {
	return fmt.Sprintf("{\"message\":\"%s\"}", msg)
}
