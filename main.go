package main

import (
	"AliceChessServer/handlers"
	"html/template"
	"io"

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
		panic(err)
	}

	renderObj := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	echoObj.Renderer = renderObj

	echoObj.GET("/", handObj.GetMain)
	echoObj.GET("/game", handObj.GetGame)

	echoObj.GET("/register", handObj.GetReg)
	echoObj.POST("/register", handObj.PostReg)

	echoObj.GET("/login", handObj.NotImplemented)
	echoObj.POST("/login", handObj.NotImplemented)

	echoObj.Static("/static", "./static")

	if err := echoObj.Start(":1323"); err != nil {
		echoObj.Logger.Error("failed to start server", "error", err)
	}
}
