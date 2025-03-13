package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	path string
}

type Task struct {
	Desc     string `json:"desc"`
	Start    string `json:"start"`
	Deadline string `json:"deadline"`
	Status   string `json:"status"`
}

func (t *Task) completed() {
	t.Status = "Completed!"
}

func (t *Task) in_process() {
	t.Status = "In process!"
}

func (t *Task) abandoned() {
	t.Status = "Later..."
}

func isValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func readTasks(path string) ([]Task, error) {
	var tasks []Task

	file, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}

	if len(file) == 0 {
		return []Task{}, nil
	}

	if err := json.Unmarshal(file, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func saveTasks(path string, tasks []Task) error {
	file, err := json.MarshalIndent(tasks, "", "	")

	if err != nil {
		return err
	}
	return os.WriteFile(path, file, 0644)
}

func addNewTask(path, taskStr string) error {
	parts := strings.Split(taskStr, ",")
	if len(parts) != 3 {
		return fmt.Errorf("Invalid task format! Expected: 'description, start, deadline'")
	}

	newTask := Task{
		Desc:     parts[0],
		Start:    parts[1],
		Deadline: parts[2],
		Status:   "Later...",
	}

	if !isValidDate(newTask.Start) || !isValidDate(newTask.Deadline) {
		return fmt.Errorf("Invalid date fmt! Expected: 'YYYY-MM-DD'")
	}

	tasks, err := readTasks(path)

	if err != nil {
		return err
	}

	tasks = append(tasks, newTask)

	return saveTasks(path, tasks)

}

func outputTasks(path string) error {
	tasks, err := readTasks(path)
	if err != nil {
		return err
	}

	tasksCount := len(tasks)
	fmt.Printf("You got %d tasks:\n", tasksCount)
	if tasksCount != 0 {

		fmt.Println("id")
		for id, task := range tasks {
			fmt.Printf("%d |	Desc: %s		Started: %s		Dedline: %s		Status: %s\n", id+1, task.Desc, task.Start, task.Deadline, task.Status)
		}
	}
	return nil
}

func deleteTask(path, taskToDel string) error {
	tasks, err := readTasks(path)

	if err != nil {
		return err
	}

	var updatedTasks []Task

	for _, task := range tasks {
		if task.Desc == taskToDel {
			continue
		}
		updatedTasks = append(updatedTasks, task)
	}

	return saveTasks(path, updatedTasks)
}

func main() {

	cfg := Config{}
	var addTask string
	var showTasks bool
	var taskToDelete string

	flag.StringVar(&cfg.path, "path", "tasks.json", "Path to tasks file")
	flag.StringVar(&addTask, "add_task", "", "Add new task in fmt: 'description, start, deadline'")
	flag.BoolVar(&showTasks, "show", false, "Output all your tasks")
	flag.StringVar(&taskToDelete, "del", "", "Delete ALL tasks w this desc")

	flag.Parse()

	if addTask != "" {
		if err := addNewTask(cfg.path, addTask); err != nil {
			fmt.Printf("Error adding task: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Task added!")
	}

	if *&showTasks {
		err := outputTasks(cfg.path)
		if err != nil {
			fmt.Printf("Error showing tasks: %v\n", err)
			os.Exit(1)
		}
	}

	if taskToDelete != "" {
		err := deleteTask(cfg.path, taskToDelete)
		if err != nil {
			fmt.Printf("Error while deleting task '%s': %v\n", taskToDelete, err)
			os.Exit(1)
		}
	}
}
