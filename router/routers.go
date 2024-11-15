package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type BaseController interface {
	Routes() []Route
	Group() string
}

type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// ContentType accepted content type
	ContentType string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// NewRouter returns a new router.
func NewRouter(controllers ...BaseController) *gin.Engine {
	router := gin.Default()
	return NewRouterWithGinEngine(router, controllers...)
}

// NewRouterWithGinEngine add routes to existing gin engine.
func NewRouterWithGinEngine(router *gin.Engine, controllers ...BaseController) *gin.Engine {
	for _, controller := range controllers {
		addRoutes(router, controller.Routes())
	}
	return router
}

// NewRouter add routes to existing gin engine.
func addRoutes(router *gin.Engine, routes []Route) {
	for _, route := range routes {
		if route.HandlerFunc == nil {
			route.HandlerFunc = DefaultHandleFunc
		}
		var handler = route.HandlerFunc
		if len(route.ContentType) != 0 {
			handler = contentTypeCheck(route.ContentType, handler)
		}
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, handler)
		case http.MethodPost:
			router.POST(route.Pattern, handler)
		case http.MethodPut:
			router.PUT(route.Pattern, handler)
		case http.MethodPatch:
			router.PATCH(route.Pattern, handler)
		case http.MethodDelete:
			router.DELETE(route.Pattern, handler)
		}
	}
}

func DefaultHandleFunc(c *gin.Context) {
	c.String(http.StatusNotImplemented, "501 not implemented")
}

func contentTypeCheck(contentType string, next gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.Contains(contentType, ctx.GetHeader("Content-Type")) {
			next(ctx)
		} else {
			ctx.String(http.StatusUnsupportedMediaType, "")
		}
	}
}
