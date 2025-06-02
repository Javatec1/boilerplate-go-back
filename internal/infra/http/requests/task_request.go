package requests

import (
	"time"

	"github.com/Javatec1/boilerplate-go-back/internal/domain"
)

type TaskRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description *string `json:"description"`
	Date        *int64  `json:"date"`
}

func (r TaskRequest) ToDomainModel() (interface{}, error) {
	var date time.Time
	if r.Date != nil {
		date = time.Unix(*r.Date, 0)
	}

	return domain.Task{
		Title:       r.Title,
		Description: r.Description,
		Date:        &date,
	}, nil
}

type UpdateTaskRequest struct {
	Status domain.TaskStatus `json:"status" validate:"required,oneof=NEW IN_PROGRESS COMPLETE"`
}

func (r UpdateTaskRequest) ToDomainModel() (interface{}, error) {
	return domain.Task{
		Status: r.Status,
	}, nil
}
