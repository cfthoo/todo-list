package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cfthoo/todo-app/pkg/db/model"
)

// TodoList handles all of the database actions
type TodoList struct {
	DB *sql.DB
}

// Create will insert a task into the database
func (t *TodoList) Create(ctx context.Context, task *model.Task) (*model.Task, error) {
	if task == nil {
		return nil, errors.New("task can not be nil")
	}
	completed := false
	now := time.Now()
	var lastInsertId int64
	statement := "INSERT INTO tasks ( name, created_by,complete, created_at, modified_at) VALUES ($1, $2,$3,$4,$5) RETURNING id"
	err := t.DB.QueryRow(statement, task.Name, task.CreatedBy, completed, now, now).Scan(&lastInsertId)
	if err != nil {
		fmt.Println("sss:", err)
		return nil, err
	}

	task.ID = int(lastInsertId)
	task.CreatedAt = now
	task.ModifiedAt = now

	return task, nil
}

// FetchAll returns all tasks
// Here i implmented the select by created_by/userId , due to we have to login with
// google/fb/github. Therefore each user can only select their own task.
func (t *TodoList) FetchAll(ctx context.Context, userId string) ([]model.Task, error) {

	statement := "SELECT * FROM tasks WHERE created_by=$1"
	rows, err := t.DB.Query(statement, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := []model.Task{}

	for rows.Next() {
		var t model.Task
		err := rows.Scan(&t.ID, &t.Name, &t.Complete, &t.CreatedBy, &t.CreatedAt, &t.ModifiedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

// FetchByID returns an task by the id
func (t *TodoList) FetchByID(ctx context.Context, id int) (*model.Task, error) {
	task := &model.Task{}
	statement := "SELECT * FROM tasks WHERE id=$1"
	rows, err := t.DB.Query(statement, id)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		err = rows.Scan(&task.ID, &task.Name, &task.Complete, &task.CreatedBy, &task.CreatedAt, &task.ModifiedAt)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("No record found")
	}

	return task, nil
}

// Update will update task
func (t *TodoList) Update(ctx context.Context, task *model.Task, userId string) (*model.Task, error) {

	now := time.Now()
	statement := "UPDATE tasks SET name=$1 , modified_at=$2 WHERE id=$3 and created_by=$4"
	_, err := t.DB.Exec(statement, task.Name, now, task.ID)
	if err != nil {
		return nil, err
	}

	res, err := t.FetchByID(context.Background(), task.ID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Update will update task
func (t *TodoList) MarkComplete(ctx context.Context, id int, userId string) (*model.Task, error) {

	now := time.Now()
	statement := "UPDATE tasks SET complete=true , modified_at=$1 WHERE id=$2 and created_by=$3"
	_, err := t.DB.Exec(statement, now, id, userId)
	if err != nil {
		return nil, err
	}

	res, err := t.FetchByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Delete will delete a task
func (t *TodoList) Delete(ctx context.Context, id int, userId string) error {
	if id < 0 {
		return fmt.Errorf("invalid id")
	}

	statement := "DELETE from tasks WHERE id=$1 and created_by=$2"
	_, err := t.DB.Exec(statement, id, userId)
	if err != nil {
		return err
	}
	return nil

}
