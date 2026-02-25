package http

import (
	"context"
	"fmt"
	"sika/internal/server/controller"
	"time"

	"github.com/labstack/echo/v4"
)

type HttpServer interface {
	Listen() error
	ShutDown()
}

type httpServer struct {
	host           string
	port           string
	handler        *echo.Echo
	userController controller.UserController
}

func NewHttpServer(host, port string, userCtrl controller.UserController) *httpServer {
	return &httpServer{
		host:           host,
		port:           port,
		handler:        echo.New(),
		userController: userCtrl,
	}
}

func (h *httpServer) Listen() error {
	h.registerRoutes()

	address := fmt.Sprintf("%s:%s", h.host, h.port)
	return h.handler.Start(address)
}

func (h *httpServer) registerRoutes() {
	h.handler.GET("/:id", h.userController.GetUser)
}

func (h *httpServer) ShutDown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h.handler.Shutdown(ctx)
}
