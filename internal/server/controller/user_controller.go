package controller

import (
	"errors"
	"log"
	"net/http"
	"sika/internal/domain"
	"sika/internal/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userSrv service.UserService) *UserController {
	return &UserController{
		userService: userSrv,
	}
}

func (u *UserController) GetUser(c echo.Context) error {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"message": "invalid user id"},
		)
	}

	user, err := u.userService.Get(c.Request().Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{"message": "invalid user id"},
			)
		default:
			log.Println("[USER] ", err)
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{"message": "server internal error"},
			)
		}

	}

	return c.JSON(
		http.StatusOK,
		map[string]any{"user": user},
	)
}
