package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

func TracingMiddleware(name string) fiber.Handler {
	tracer := otel.Tracer(name)

	return func(c *fiber.Ctx) error {
		// Start a new span
		ctx, span := tracer.Start(c.UserContext(), c.Path())
		defer span.End()
		// Pass the context with the span to the next handler
		c.SetUserContext(ctx)
		return c.Next()
	}
}
