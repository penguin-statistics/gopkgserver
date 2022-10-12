package main

import (
	"html/template"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
)

const (
	HtmlContent = `<html>
  <head>
    <meta name="go-import" content="github.com/penguin-statistics/{{ .pkg }} git https://github.com/penguin-statistics/{{ .pkg }}.git">
  </head>
</html>`
)

func main() {
	app := fiber.New()

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

		return tmpl.Execute(c, fiber.Map{
			"pkg": pkg,
		})
	})

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
