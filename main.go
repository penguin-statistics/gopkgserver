package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const (
	HtmlContent = `<html>
  <head>
    <meta name="go-import" content="github.com/%s/%s git https://github.com/%s/%s.git">
  </head>
</html>`
)

func main() {
	app := fiber.New()

	app.Get("/:pkg", func(c *fiber.Ctx) error {
		project := "penguin-statistics"
		pkg := c.Params("pkg")

		if project == "" || pkg == "" {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.SendString(fmt.Sprintf(HtmlContent, project, pkg, project, pkg))
	})

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
