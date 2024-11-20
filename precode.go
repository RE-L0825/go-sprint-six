package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
func getAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allTasks := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		allTasks = append(allTasks, task)
	}
	if err := json.NewEncoder(w).Encode(allTasks);
	err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func SendTasks(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask);
	err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	newID := strconv.Itoa(len(tasks) + 1)
	newTask.ID = newID
	tasks[newID] = newTask
	w.WriteHeader(http.StatusCreated)
}

func GetTaskID(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	task, exists := tasks[taskID]
	if !exists {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task);
	err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteTaskID(w http.ResponseWriter, r *http.Request) {
	taskID :=chi.URLParam(r, "id")
	_, exists := tasks[taskID]
	if !exists {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	delete(tasks, taskID)
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	
	r.Get("/tasks", getAllTasks)
	r.Post("/tasks", SendTasks)
	r.Get("/tasks/{id}", GetTaskID)
	r.Delete("/tasks/{id}", DeleteTaskID)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
