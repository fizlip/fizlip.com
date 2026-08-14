package main

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Count struct {
	Count int
}

type Article struct {
	Title string
	Date  string
	Body  template.HTML
}

func main() {
	e := echo.New()
	e.Renderer = NewTemplates()
	e.Use(middleware.Logger())

	// Serve CSS files from /css
	e.Static("/css", "css")
	// Serve static assets (images, etc.)
	e.Static("/static", "static")

	count := Count{Count: 0}

	e.GET("/", func(c echo.Context) error {
		count.Count++
		return c.Render(200, "index", count)
	})

	e.GET("/about", func(c echo.Context) error {
		return c.Render(200, "about", nil)
	})

	e.GET("/projects", func(c echo.Context) error {
		return c.Render(200, "projects", nil)
	})

	e.GET("/blog", func(c echo.Context) error {
		return c.Render(200, "blog", nil)
	})

	e.GET("/articles/:year/:month/:slug", func(c echo.Context) error {
		year := c.Param("year")
		month := c.Param("month")
		slug := c.Param("slug")

		// Reject path traversal attempts
		for _, p := range []string{year, month, slug} {
			if strings.Contains(p, "..") || strings.Contains(p, "/") || strings.Contains(p, "\\") {
				return c.String(400, "invalid path")
			}
		}

		path := filepath.Join("articles", year, month, slug+".html")
		content, err := os.ReadFile(path)
		if err != nil {
			return c.String(404, "article not found")
		}

		article := Article{
			Title: strings.ReplaceAll(slug, "-", " "),
			Date:  month + "/" + year,
			Body:  template.HTML(content),
		}

		return c.Render(200, "article", article)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
