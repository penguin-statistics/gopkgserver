package main

import (
	"html/template"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/helmet/v2"
)

const (
	HtmlContent = `<html>
  <head>
    <meta name="go-import" content="exusiai.dev/{{ .pkg }} git https://github.com/penguin-statistics/{{ .pkg }}.git">
    <meta http-equiv="refresh" content="0; url=https://pkg.go.dev/exusiai.dev/{{ .pkg }}">
  </head>
</html>`
)

func main() {
	app := fiber.New()
	app.Use(helmet.New(helmet.Config{
		ContentSecurityPolicy: "default-src 'none'; sandbox allow-scripts",
		ReferrerPolicy:        "no-referrer",
	}))

	pkgMatcher := regexp.MustCompile(`^[a-z0-9-]+$`)

	tmpl := template.Must(template.New("index").Parse(HtmlContent))

	app.Get("/:pkg", func(c *fiber.Ctx) error {
		pkg := c.Params("pkg")

		err := validation.Validate(pkg,
			validation.Required,
			validation.Length(1, 64),
			validation.NotIn("favicon.ico", "robots.txt"),
			validation.Match(pkgMatcher),
		)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		c.Set("Content-Type", "text/html; charset=utf-8")
		c.Set("Cache-Control", "public, max-age=86400")
		c.Set("Vary", "Accept-Encoding")

		return tmpl.Execute(c, fiber.Map{
			"pkg": pkg,
		})
	})

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
