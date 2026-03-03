package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

const fileName = "tasks.json"

func main() {
	add := flag.String("add", "", "Добавить новую задачу")
	list := flag.Bool("list", false, "Показать список задач")
	del := flag.Int("del", 0, "Удалить задачу по ID")
	done := flag.Int("done", 0, "Отметить задачу как выполненную по ID")
	flag.Parse()

	switch {
	case *add != "":
		addTask(*add)
	case *list:
		listTasks()
	case *del != 0:
		deleteTask(*del)
	case *done != 0:
		completeTask(*done)
	default:
		fmt.Println("Использование:")
		fmt.Println("  -add \"текст\"   Добавить задачу")
		fmt.Println("  -list           Показать список")
		fmt.Println("  -del ID         Удалить задачу")
		fmt.Println("  -done ID        Выполнить задачу")
	}
}

func addTask(text string) {
	tasks, _ := loadTasks()

	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	newTask := Task{
		ID:        maxID + 1,
		Text:      text,
		Done:      false,
		CreatedAt: time.Now(),
	}
	tasks = append(tasks, newTask)
	saveTasks(tasks)
	fmt.Printf("Задача добавлена: %s\n", text)
}

func listTasks() {
	tasks, _ := loadTasks()
	if len(tasks) == 0 {
		fmt.Println("Список задач пуст.")
		return
	}
	for _, t := range tasks {
		status := " "
		if t.Done {
			status = "X"
		}
		fmt.Printf("[%d] [%s] %s\n", t.ID, status, t.Text)
	}
}

func saveTasks(tasks []Task) {
	data, _ := json.MarshalIndent(tasks, "", " ")
	os.WriteFile(fileName, data, 0644)
}

func loadTasks() ([]Task, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return []Task{}, nil
	}
	var tasks []Task
	json.Unmarshal(data, &tasks)
	return tasks, nil
}

func deleteTask(id int) {
	tasks, _ := loadTasks()
	newTasks := []Task{}
	found := false

	for _, t := range tasks {
		if t.ID != id {
			newTasks = append(newTasks, t)
		} else {
			found = true
		}
	}
	if found {
		saveTasks(newTasks)
		fmt.Printf("Задача №%d удалена.\n", id)
	} else {
		fmt.Printf("Задача с ID %d не найдена.\n", id)
	}
}

func completeTask(id int) {
	tasks, _ := loadTasks()
	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = true
			found = true
			break
		}
	}
	if found {
		saveTasks(tasks)
		fmt.Printf("Задача №%d отмечена как выполненная!\n", id)
	} else {
		fmt.Printf("Задача с ID %d не найдена.\n", id)
	}
}
