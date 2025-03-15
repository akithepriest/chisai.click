package middlewares

import (
	"github.com/akithepriest/chisai.click/internal"
	"github.com/labstack/echo"
)

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func (c echo.Context) error {
		internal.Logger.LogInfo().Fields(map[string]interface{}{
			"method": c.Request().Method,
			"path":    c.Request().URL.Path,
			"query":  c.Request().URL.RawQuery,
		   }).Msg("Request")

		if err := next(c); err != nil {
			internal.Logger.LogError().Fields(map[string]interface{} {
				"method": c.Request().Method,
				"path": c.Request().URL.Path,
				"query": c.Request().URL.RawQuery,
				"error": err.Error(),
			}).Msg("Response")
			return err
		}
		return nil
	}
}