package app

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

type LoginThrottle struct {
}

type Payload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// New : To create a new LoginThrottle Service
func New() (*LoginThrottle, error) {
	return &LoginThrottle{}, nil
}

func (app LoginThrottle) Login(c echo.Context) error {
	var parsed Payload

	err := json.NewDecoder(c.Request().Body).Decode(&parsed)
	if err != nil {
		return c.JSON(400, "Bad request")
	}

	if parsed.Username == "sample" && parsed.Password == "test" {
		return c.JSON(200, "OK")
	}
	return c.NoContent(401)
}
