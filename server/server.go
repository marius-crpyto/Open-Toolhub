package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"toolhub/logger"
)

type Server struct {
	engine *gin.Engine
	logger *logger.Logger

	apiPathList   []ApiPath
	filterMapping map[string]FilterHandler
}

func NewServer(ginMode string, corsAllowOrigins, allowHeaders []string, logger *logger.Logger) *Server {
	gin.SetMode(ginMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(NoCache())
	engine.Use(CORS(corsAllowOrigins, allowHeaders))
	engine.Use(RequestID())
	engine.Use(AccessLog(logger))

	return &Server{
		engine:        engine,
		logger:        logger,
		apiPathList:   []ApiPath{},
		filterMapping: map[string]FilterHandler{},
	}
}

func (s *Server) Router(group string) {
	s.engine.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Working!")
	})

	if group == "" {
		group = "/api"
	}

	api := s.engine.Group(group)
	api.Use(s.AuthMiddleware(s.logger))

	for _, apiPath := range s.apiPathList {
		s.logger.Info("register api path", zap.String("path", apiPath.Path), zap.String("method", apiPath.Method))
		s.initFilter(api.BasePath(), &apiPath)
		api.Handle(apiPath.Method, apiPath.Path, apiPath.Handler)
	}

}

func (s *Server) Run(addr string) error {
	s.logger.Info("server run", zap.String("addr", addr))
	return s.engine.Run(addr)
}

func (s *Server) AddApiPaths(paths ...ApiPath) {
	s.apiPathList = append(s.apiPathList, paths...)
}

func (s *Server) AddApiPath(path ApiPath) {
	s.apiPathList = append(s.apiPathList, path)
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}
