package middleware

import (
	"context"
	"crud_service/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type contextKey string

const RolesKey contextKey = "roles"
const UserIDs contextKey = "ids"

func RolesMiddleware(UserService service.UserService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userIDStr := c.Request().Header.Get("X-User-ID")

			if userIDStr == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing X-User-ID header",
				})
			}

			userId, err := strconv.Atoi(userIDStr)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid X-User-ID header",
				})
			}

			log.Println(userId, ": middleware user id")

			roles, err := UserService.GetRoles(c.Request().Context(), userId)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "user not found",
				})
			}

			ctx := context.WithValue(c.Request().Context(), RolesKey, roles)
			ctx = context.WithValue(ctx, UserIDs, userId)

			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
