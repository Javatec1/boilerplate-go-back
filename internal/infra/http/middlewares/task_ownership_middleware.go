package middlewares

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Javatec1/boilerplate-go-back/internal/app"
	"github.com/Javatec1/boilerplate-go-back/internal/domain"
	"github.com/Javatec1/boilerplate-go-back/internal/infra/http/controllers"
	"github.com/go-chi/chi/v5"
)

func TaskOwnershipMiddleware(ts app.TaskService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Витягуємо taskId із URL
			taskIdStr := chi.URLParam(r, "taskId")
			taskId, err := strconv.ParseUint(taskIdStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid task ID", http.StatusBadRequest)
				return
			}

			// Отримуємо задачу через TaskService
			taskObj, err := ts.Find(taskId)
			if err != nil {
				http.Error(w, "Task not found", http.StatusNotFound)
				return
			}
			task := taskObj.(domain.Task)

			// Отримуємо користувача з контексту (додано AuthMw)
			user, ok := r.Context().Value(controllers.UserKey).(domain.User)
			if !ok {
				http.Error(w, "User not found in context", http.StatusInternalServerError)
				return
			}

			// Перевіряємо належність
			if task.UserId != user.Id {
				http.Error(w, "Forbidden: Task does not belong to user", http.StatusForbidden)
				return
			}

			// Додаємо задачу до контексту
			ctx := context.WithValue(r.Context(), controllers.TaskKey, task)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
