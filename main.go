package main

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

const ApplicationName = "strx"

var urlStore Store

func main() {
	app := fiber.New()
	app.Use(logger.New())

	// Display all aliases
	app.Get("/", func(c *fiber.Ctx) error {
		all := urlStore.All()

		switch c.Accepts("html", "json") {
		case "json":
			return c.JSON(all)
		default:
			c.Set("Content-Type", "text/html; charset=utf-8")
			return HTML(all, c)
		}
	})

	// Create a new alias
	app.Post("/create", func(c *fiber.Ctx) error {
		type Request struct {
			URL   string `json:"url"`
			Alias string `json:"alias,omitempty"`
		}

		req := new(Request)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		if req.URL == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "URL is required"})
		}

		alias := req.Alias
		if alias == "" {
			word1, word2 := randomWords()
			alias = word1 + "-" + word2
		}

		urlStore.Set(alias, req.URL)
		return c.JSON(fiber.Map{"alias": alias, "url": req.URL})
	})

	// Resolve an alias
	app.Get("/:alias", func(c *fiber.Ctx) error {
		alias := c.Params("alias")

		url, exists := urlStore.Get(alias)
		if !exists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Alias not found"})
		}
		return c.Redirect(url, fiber.StatusMovedPermanently)
	})

	port := flag.String("port", "3000", "Port to run the server on")
	store := flag.String("store", "file", "Store type: memory or file")

	flag.Parse()

	if *store == "file" {
		fileStore := NewFilesStore(ApplicationName)
		log.Printf("Using file store: %s", fileStore.directoryPath)
		urlStore = fileStore
	} else {
		log.Printf("Using in-memory store")
		urlStore = NewInMemoryStore()
	}

	app.Listen(":" + *port)
}
