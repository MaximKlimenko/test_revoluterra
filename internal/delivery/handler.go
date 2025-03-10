package delivery

import (
	"time"

	"github.com/MaximKlimenko/scheduler/internal/storages"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type JobReq struct {
	Description string    `json:"description"`
	ExecuteAt   time.Time `json:"executeAt"`
}

func (r *Repository) CreateJob(ctx *fiber.Ctx) error {
	var input JobReq
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
			"data":  input,
		})
	}

	if input.ExecuteAt.Before(time.Now()) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Переданное время выполнения находится в прошлом",
		})
	}

	job := &storages.Job{
		ID:          uuid.New().String(),
		Description: input.Description,
		ExecuteAt:   input.ExecuteAt,
		Status:      storages.Scheduled,
	}

	if err := r.DB.CreateJob(job); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ошибка в создании задачи",
		})
	}
	go r.scheduleJob(job)
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Задача успешно создана!",
	})
}

func (r *Repository) GetJobs(ctx *fiber.Ctx) error {
	jobs, err := r.DB.GetJobs()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"jobs": jobs,
	})
}

func (r *Repository) GetJobByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	job, err := r.DB.GetJobByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err": "Задача не найдена",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"job": job,
	})
}

func (r *Repository) CancelJob(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := r.DB.CancelJob(id); err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"err": err,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Задача успешно отменена!",
	})
}

func (r *Repository) scheduleJob(job *storages.Job) {
	time.Sleep(time.Until(job.ExecuteAt))
	r.DB.UpdateStatus(job.ID)
}
