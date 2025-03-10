package delivery

import (
	"github.com/MaximKlimenko/scheduler/internal/storages"
	"github.com/gofiber/fiber/v2"
)

type Repository struct {
	DB storages.Storage
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	app.Post("/jobs", r.CreateJob)
	app.Get("/jobs", r.GetJobs)
	app.Get("/jobs/:id", r.GetJobByID)
	app.Delete("/jobs/:id", r.CancelJob)
}
