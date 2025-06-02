package controllers

import (
	"log"
	"net/http"

	"github.com/Javatec1/boilerplate-go-back/internal/app"
	"github.com/Javatec1/boilerplate-go-back/internal/domain"
	"github.com/Javatec1/boilerplate-go-back/internal/infra/http/requests"
	"github.com/Javatec1/boilerplate-go-back/internal/infra/http/resources"
)

type TaskController struct {
	taskService app.TaskService
}

func NewTaskController(ts app.TaskService) TaskController {
	return TaskController{
		taskService: ts,
	}
}

func (c TaskController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController: %s", err)
			BadRequest(w, err)
			return
		}

		task.Status = domain.TaskNew
		user := r.Context().Value(UserKey).(domain.User)
		task.UserId = user.Id

		task, err = c.taskService.Save(task)
		if err != nil {
			log.Printf("TaskController: %s", err)
			InternalServerError(w, err)
			return
		}

		var taskDto resources.TaskDto
		taskDto = taskDto.DomainToDto(task)
		Success(w, taskDto)
	}
}

func (c TaskController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task := r.Context().Value(TaskKey).(domain.Task)
		var taskDto resources.TaskDto
		taskDto = taskDto.DomainToDto(task)
		Success(w, taskDto)
	}
}

func (c TaskController) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		tasks, err := c.taskService.FindAll(user.Id)
		if err != nil {
			log.Printf("TaskController.FindAll(c.taskService.FindAll): %s", err)
			InternalServerError(w, err)
			return
		}

		var taskDto resources.TaskDto
		tasksDto := taskDto.DomainToDtoCollection(tasks)
		Success(w, tasksDto)
	}
}

func (c TaskController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController: %s", err)
			BadRequest(w, err)
			return
		}

		taskExists := r.Context().Value(TaskKey).(domain.Task)
		taskExists.Title = task.Title
		taskExists.Description = task.Description
		taskExists.Date = task.Date

		task, err = c.taskService.Update(taskExists)
		if err != nil {
			log.Printf("TaskController: %s", err)
			InternalServerError(w, err)
			return
		}

		var taskDto resources.TaskDto
		taskDto = taskDto.DomainToDto(task)
		Success(w, taskDto)
	}
}

func (c TaskController) UpdateStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, err := requests.Bind(r, requests.UpdateTaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController.UpdateStatus: %s", err.Error())
			BadRequest(w, err)
			return
		}

		taskExists := r.Context().Value(TaskKey).(domain.Task)
		taskExists.Status = task.Status

		updatedTask, err := c.taskService.UpdateStatus(taskExists.Id, taskExists.Status)
		if err != nil {
			log.Printf("TaskController.UpdateStatus: %s", err.Error())
			InternalServerError(w, err)
			return
		}

		var taskDto resources.TaskDto
		taskDto = taskDto.DomainToDto(updatedTask)
		Success(w, taskDto)
	}
}

func (c TaskController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task := r.Context().Value(TaskKey).(domain.Task)
		err := c.taskService.Delete(task.Id)
		if err != nil {
			log.Printf("TaskController: %s", err)
			InternalServerError(w, err)
			return
		}
		noContent(w)
	}
}
