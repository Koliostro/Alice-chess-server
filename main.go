package main

import (
	"AliceChessServer/handlers"
	"html/template"
	"io"
	"log"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (templRender *TemplateRenderer) Render(context *echo.Context, writer io.Writer, name string, data any) error {
	return templRender.templates.ExecuteTemplate(writer, name, data)
}

func main() {
	echoObj := echo.New()
	echoObj.Use(middleware.RequestLogger())

	handObj, err := handlers.NewGenericHandler()

	if err != nil {
		log.Fatal("Unnable to connect to DB")
	}

	renderObj := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	echoObj.Renderer = renderObj

	echoObj.GET("/", handObj.GetMain)

	authGroup := echoObj.Group("/auth")

	authGroup.GET("/register", handObj.GetReg)
	authGroup.POST("/register", handObj.PostReg)
	authGroup.GET("/login", handObj.GetLogin)
	authGroup.GET("/login/logout", handObj.GetLogOut)
	authGroup.POST("/login", handObj.PostLogin)

	gameGroup := echoObj.Group("/games")

	gameGroup.GET("/connectionMenu", handObj.GetConnectionMenu)
	gameGroup.POST("/closeGame/:id", handObj.PostCloseGame)
	gameGroup.GET("/createGame", handObj.GetCreateRoom)
	gameGroup.GET("/connect/:id", handObj.NotImplemented)
	gameGroup.GET("/:id", handObj.GetwaitingRoom)
	gameGroup.GET("/:id/getTurn", handObj.GetGameState)
	gameGroup.POST("/:id/newState", handObj.PostNewState)

	echoObj.Static("/static", "./static")

	if err := echoObj.Start(":1323"); err != nil {
		echoObj.Logger.Error("failed to start server", "error", err)
	}
}
