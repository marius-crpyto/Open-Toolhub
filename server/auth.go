package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/marius-crpyto/toolhub/logger"
)

var (
	ERROR_AUTH_FAIL         = Response{Code: -1, Message: "auth fail"}
	ERROR_PATH_NOT_REGISTER = Response{Code: -2, Message: "path not register"}
)

func (s *Server) initFilter(group string, apiPath *ApiPath) {
	if apiPath.Auth == nil {
		apiPath.Auth = func(c *gin.Context) error {
			return nil
		}
	}

	s.filterMapping[group+apiPath.Path] = apiPath.Auth

}

func (s *Server) AuthMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentPath := c.FullPath()
		if handle, ok := s.filterMapping[currentPath]; !ok {
			log.Error("path not register", zap.String("path", currentPath))
			c.JSON(http.StatusOK, ERROR_PATH_NOT_REGISTER)
			c.Abort()
			return
		} else {
			if err := handle(c); err != nil {
				log.Error("failed to check auth", zap.String("error", err.Error()))
				c.JSON(http.StatusOK, ERROR_AUTH_FAIL)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
