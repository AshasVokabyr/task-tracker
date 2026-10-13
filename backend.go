package main

import (
	"errors"
	"fmt"
	"time"
)

var errNotFound = errors.New("not found")

type Task struct {
	id          int
	name        string
	priority    string
	deadline    string
	assignee    string
	author      string
	update_date time.Time
}

type Storage interface {
	CreateTask(name, priority, deadline, assignee, author string) Task
	GetTask(i int) (Task, error)
	UpdateTask(i int, name, priority, deadline, assignee, author string) (int, error)
	DeleteTask(i int) (int, error)
}

type MemoryStorage struct {
	tasks  map[int]Task
	nextID int
}

func (ms *MemoryStorage) CreateTask(name, priority, deadline, assignee, author string) Task {
	ms.nextID++
	task := Task{
		id:          ms.nextID,
		name:        name,
		priority:    priority,
		deadline:    deadline,
		assignee:    assignee,
		author:      author,
		update_date: time.Now(),
	}
	ms.tasks[task.id] = task

	return task
}

func (ms *MemoryStorage) GetTask(id int) (Task, error) {
	_, exists := ms.tasks[id]
	if !exists {
		return Task{}, errNotFound
	}
	return ms.tasks[id], nil
}

func (ms *MemoryStorage) UpdateTask(id int, name, priority, deadline, assignee, author string) (int, error) {
	_, exists := ms.tasks[id]
	if !exists {
		return 404, errNotFound
	}
	task := Task{
		id:          id,
		name:        name,
		priority:    priority,
		deadline:    deadline,
		assignee:    assignee,
		author:      author,
		update_date: time.Now(),
	}
	ms.tasks[id] = task

	return 200, nil
}

func (ms *MemoryStorage) DeleteTask(id int) (int, error) {
	_, exists := ms.tasks[id]
	if !exists {
		return 404, errNotFound
	}
	delete(ms.tasks, id)
	return 200, nil
}

func main() {
	var store Storage = &MemoryStorage{
		tasks: make(map[int]Task),
	}

	task := store.CreateTask("Fix bug", "High", "2026-10-15", "Bob", "Alice")
	fmt.Printf("CreateTask: %+v\n", task)

	foundTask, err := store.GetTask(task.id)
	fmt.Printf("GetTask: %+v, err: %v\n", foundTask, err)

	statusUpdate, err := store.UpdateTask(task.id, "Fix bug immediately", "Critical", "2026-10-14", "Bob", "Alice")
	updatedTask, _ := store.GetTask(task.id)
	fmt.Printf("UpdateTask status: %d, err: %v, updated task: %+v\n", statusUpdate, err, updatedTask)

	statusDelete, err := store.DeleteTask(task.id)
	_, errAfterDelete := store.GetTask(task.id)
	fmt.Printf("DeleteTask status: %d, err: %v, get after delete err: %v\n", statusDelete, err, errAfterDelete)
}
