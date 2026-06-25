package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ANSI color codes for terminal output
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorCyan   = "\033[36m"
)

// PerformanceLogger is a middleware that logs the execution time of each request.
// It highlights slow requests to help identify bottlenecks.
func PerformanceLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Execute the next middleware/handler
		err := c.Next()

		latency := time.Since(start)
		method := c.Method()
		path := c.Path()

		var color string
		var speed string

		// Categorize performance
		if latency < 50*time.Millisecond {
			color = colorGreen
			speed = "FAST"
		} else if latency < 200*time.Millisecond {
			color = colorYellow
			speed = "MODERATE"
		} else {
			color = colorRed
			speed = "SLOW"
		}

		// Log the custom performance format
		log.Printf("%s[PERFORMANCE] %s | %-8s | %s %s | Took: %v%s",
			color,
			time.Now().Format("15:04:05"),
			speed,
			method,
			path,
			latency,
			colorReset,
		)

		return err
	}
}
