package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/jrdnull/microservices_technical_test/article_service/internal"
)

func errorHandler(err error, c echo.Context) {
	if err == nil {
		return
	}

	var (
		cause = err
		code  = 0
		msg   = http.StatusText(http.StatusInternalServerError)
	)
	for cause != nil {
		switch v := cause.(type) {
		case *internal.Error:
			switch v.Code {
			case internal.ErrorCodeValidation:
				code = http.StatusBadRequest
			case internal.ErrorCodeUnauthorized:
				code = http.StatusUnauthorized
			}
			msg = err.Error()

		case *echo.HTTPError:
			code = v.Code
			msg = fmt.Sprintf("%v", v.Message)

		default:
			if cause == context.DeadlineExceeded {
				code = http.StatusGatewayTimeout
				msg = http.StatusText(code)
			}
		}
		if code > 0 {
			break
		}
		cause = errors.Unwrap(cause)
	}
	if code == 0 {
		code = http.StatusInternalServerError
	}
	if err := c.JSON(code, map[string]interface{}{
		"error": msg,
	}); err != nil {
		log.Print(err)
	}
}
