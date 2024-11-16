package exception

import (
	"device-management/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/richzw/gin-error"
	"net/http"
	"strconv"
	"strings"
)

var (
	InvalidUserIdParam = &BadRequest{"userId is invalid or not provided"}
	errorList          = []HttpResponseError{InvalidUserIdParam}
)

func AddErrorHandler(engine *gin.Engine) {

	for _, e := range errorList {
		engine.Use(err.Error(err.NewErrMap(e).StatusCode(e.StatusCode())))
	}

}

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Type is %t", err)
				switch err.(type) {
				case HttpResponseError:
					apiError := err.(HttpResponseError)
					writeJson(ctx, apiError, (apiError).StatusCode())
				case *strconv.NumError:
					apiError := err.(*strconv.NumError)
					writeJson(ctx, model.NewApiError("400", "NumError", "Invalid number format", apiError.Error()), 400)
				case validator.ValidationErrors:
					apiError := err.(validator.ValidationErrors)
					writeJson(ctx, model.NewApiError("400", "ValidationError", "Missing required field, or wrong field value", apiError.Error()), 400)
				case error:
					apiError := err.(error)
					if apiError.Error() == "record not found" {
						writeJson(ctx, model.NewApiError("404", "record not found", "Invalid number format", apiError.Error()), 404)
					} else if strings.Contains(apiError.Error(), "invalid UUID length") {
						writeJson(ctx, model.NewApiError("400", "invalid uuid", "Invalid uuid format", apiError.Error()), 400)
					} else if strings.Contains(apiError.Error(), "duplicate key value violates unique constraint ") {
						writeJson(ctx, model.NewApiError("409", "Unique key violation", "Unique key violation", apiError.Error()), http.StatusConflict)
					} else {
						writeJson(ctx, model.NewApiError("500", apiError.Error(), "InternalServerError", apiError.Error()), 500)
					}
				default:
					writeJson(ctx, model.NewApiError("500", err.(string), "InternalServerError", err.(string)), 500)
				}
			}
		}()
		ctx.Next()
	}
}

func writeJson(ctx *gin.Context, data interface{}, code int) {
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.AbortWithStatusJSON(code, data)
}
