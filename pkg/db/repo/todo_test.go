package repo

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cfthoo/todo-app/pkg/db/model"
)

func TestTodo_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	// create a new TodoList instance
	todoList := &TodoList{DB: db}

	// mock the database response
	task := &model.Task{Name: "test task", CreatedBy: "123"}
	mock.ExpectQuery("^INSERT INTO tasks").
		WithArgs(task.Name, task.CreatedBy, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// call the Create method
	result, err := todoList.Create(context.Background(), task)

	// check if there are any errors
	if err != nil {
		t.Fatalf("failed to create todo list: %v", err)
	}

	// check if the returned result is correct
	if result.ID != 1 {
		t.Errorf("expected ID to be 1, but got %d", result.ID)
	}

	// check if the CreatedAt and ModifiedAt fields are set
	if result.CreatedAt.IsZero() {
		t.Error("expected CreatedAt field to be set, but it's zero")
	}
	if result.ModifiedAt.IsZero() {
		t.Error("expected ModifiedAt field to be set, but it's zero")
	}

	// check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTodo_FetchAll(t *testing.T) {

	userId := "123"
	tasks := []model.Task{
		{
			ID:         1,
			Name:       "Task 1",
			CreatedBy:  userId,
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		},
		{
			ID:         2,
			Name:       "Task 2",
			CreatedBy:  userId,
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	todoList := &TodoList{
		DB: db,
	}

	columns := []string{"id", "name", "created_by", "created_at", "modified_at"}
	rows := sqlmock.NewRows(columns)

	for _, task := range tasks {
		rows.AddRow(task.ID, task.Name, task.CreatedBy, task.CreatedAt, task.ModifiedAt)
	}

	mock.ExpectQuery("SELECT \\* FROM tasks WHERE created_by=\\$1").
		WithArgs(userId).
		WillReturnRows(rows)

	result, err := todoList.FetchAll(context.Background(), userId)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if len(result) != len(tasks) {
		t.Errorf("Result length is %v, expected %v", len(result), len(tasks))
		return
	}

	for i, task := range tasks {
		if result[i].ID != task.ID || result[i].Name != task.Name || result[i].CreatedBy != task.CreatedBy {
			t.Errorf("Unexpected result. Got %v, expected %v", result[i], task)
			return
		}
	}
}

func TestTodo_FetchByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Set up the expected rows to be returned by the mock
	expectedTask := &model.Task{ID: 1, Name: "Task 1", CreatedBy: "user1", CreatedAt: time.Now(), ModifiedAt: time.Now()}
	rows := sqlmock.NewRows([]string{"id", "name", "created_by", "created_at", "modified_at"}).
		AddRow(expectedTask.ID, expectedTask.Name, expectedTask.CreatedBy, expectedTask.CreatedAt, expectedTask.ModifiedAt)

	// Set up the mock query and result
	mock.ExpectQuery("SELECT \\* FROM tasks WHERE id=\\$1").WithArgs(1).WillReturnRows(rows)

	// Call the function being tested
	list := &TodoList{DB: db}
	task, err := list.FetchByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("FetchByID returned an error: %v", err)
	}

	// Verify that the mock was called as expected
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Unfulfilled expectations: %s", err)
	}

	// Verify that the returned task matches the expected task
	if !reflect.DeepEqual(task, expectedTask) {
		t.Fatalf("Returned task %+v does not match expected task %+v", task, expectedTask)
	}
}

func TestTodo_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	// Set up the expected task and mock query result
	now := time.Now()
	expectedTask := &model.Task{ID: 1, Name: "Updated Task 1", CreatedBy: "user1", CreatedAt: now, ModifiedAt: now}
	mock.ExpectExec("UPDATE tasks SET name=\\$1 , modified_at=\\$2 WHERE id=\\$3").
		WithArgs(expectedTask.Name, now, expectedTask.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT \\* FROM tasks WHERE id=\\$1").
		WithArgs(expectedTask.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_by", "created_at", "modified_at"}).
			AddRow(expectedTask.ID, expectedTask.Name, expectedTask.CreatedBy, expectedTask.CreatedAt, expectedTask.ModifiedAt))

	// Call the function being tested
	list := &TodoList{DB: db}
	res, err := list.Update(context.Background(), expectedTask)
	if err != nil {
		t.Fatalf("Update returned an error: %v", err)
	}

	// Verify that the mock was called as expected
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Unfulfilled expectations: %s", err)
	}

	// Verify that the returned task matches the expected task
	if !reflect.DeepEqual(res, expectedTask) {
		t.Fatalf("Returned task %+v does not match expected task %+v", res, expectedTask)
	}
}

func TestTodo_Delete(t *testing.T) {
	// create mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	// create TodoList with mock database
	todoList := &TodoList{DB: db}

	// set up test case
	id := 1
	mock.ExpectExec("DELETE from tasks WHERE id=\\$1").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// call method
	err = todoList.Delete(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// assert expectations were met
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}
