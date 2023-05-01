package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/cfthoo/todo-app/pkg/db/model"
)

type mockTodoListDAO struct {
	tasks []model.Task
	task  model.Task
	err   error
}

func (m *mockTodoListDAO) Create(ctx context.Context, task *model.Task) (*model.Task, error) {
	//m.tasks = append(m.tasks, *task)
	return task, nil
}

func (m *mockTodoListDAO) FetchAll(ctx context.Context, userId string) ([]model.Task, error) {
	return m.tasks, m.err
}

func (m *mockTodoListDAO) FetchByID(ctx context.Context, id int) (*model.Task, error) {
	//m.tasks = append(m.tasks, *task)
	return &m.task, nil
}

func (m *mockTodoListDAO) Update(ctx context.Context, task *model.Task) (*model.Task, error) {
	//m.tasks = append(m.tasks, *task)
	return task, nil
}

func (m *mockTodoListDAO) Delete(ctx context.Context, id int) error {
	//m.tasks = append(m.tasks, *task)
	return nil
}

func TestHandler_Create(t *testing.T) {
	type fields struct {
		TodoListDAO DAO
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		status int
		body   interface{}
	}{
		{
			name: "created",
			fields: fields{
				TodoListDAO: &mockTodoListDAO{
					tasks: []model.Task{
						{
							ID:        1,
							Name:      "task1",
							CreatedBy: "1234",
						},
					},
				},
			},
			args: args{
				req: func() *http.Request {
					u := model.Task{
						Name: "task1",
						//	LastName:  "testison",
					}
					enc, _ := json.Marshal(u)
					return httptest.NewRequest(http.MethodPost, "http://www.google.com", bytes.NewReader(enc))
				}(),
			},
			status: http.StatusCreated,
			body: model.Task{
				Name: "task1",
				//	LastName:  "testison",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				TodoListDAO: tt.fields.TodoListDAO,
			}
			writer := httptest.NewRecorder()
			handler := h.Create()
			handler.ServeHTTP(writer, tt.args.req)

			if writer.Result().StatusCode != tt.status {
				t.Errorf("Handler.Create() = %v, want %v", writer.Result().StatusCode, tt.status)
				return
			}

			var bodyMap map[string]interface{}
			if err := json.NewDecoder(writer.Body).Decode(&bodyMap); err != nil {
				t.Errorf("Handler.Create() = json body decode error %v", err)
				return
			}

			var wantBodyMap map[string]interface{}
			if enc, err := json.Marshal(tt.body); err == nil {
				_ = json.Unmarshal(enc, &wantBodyMap)
			}

			if !reflect.DeepEqual(bodyMap, wantBodyMap) {
				t.Errorf("Handler.Create() = %v, want %v", bodyMap, wantBodyMap)
			}
		})
	}
}

func TestHandler_List(t *testing.T) {
	type fields struct {
		TodoListDAO DAO
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		status int
		body   interface{}
	}{
		{
			name: "fetched",
			fields: fields{
				TodoListDAO: &mockTodoListDAO{
					tasks: []model.Task{
						{
							ID:        1,
							Name:      "task1",
							CreatedBy: "1234",
						},
					},
				},
			},
			args: args{
				req: func() *http.Request {
					return httptest.NewRequest(http.MethodGet, "http://www.google.com/", strings.NewReader(""))
				}(),
			},
			status: http.StatusOK,
			body: []*model.Task{
				{
					ID:        1,
					Name:      "task1",
					CreatedBy: "1234",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				TodoListDAO: tt.fields.TodoListDAO,
			}
			writer := httptest.NewRecorder()
			handler := h.List()
			handler.ServeHTTP(writer, tt.args.req)

			if writer.Result().StatusCode != tt.status {
				t.Errorf("Handler.List() = %v, want %v", writer.Result().StatusCode, tt.status)
				return
			}

			var bodyMap []interface{}
			if err := json.NewDecoder(writer.Body).Decode(&bodyMap); err != nil {
				t.Errorf("Handler.List() = json body decode error %v", err)
				return
			}

			var wantBodyMap []interface{}
			if enc, err := json.Marshal(tt.body); err == nil {
				_ = json.Unmarshal(enc, &wantBodyMap)
			}

			if !reflect.DeepEqual(bodyMap, wantBodyMap) {
				t.Errorf("Handler.List() = %v, want %v", bodyMap, wantBodyMap)
			}
		})
	}
}

func TestHandler_FetchByID(t *testing.T) {
	type fields struct {
		TodoListDAO DAO
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		status int
		body   interface{}
	}{
		{
			name: "fetched",
			fields: fields{
				TodoListDAO: &mockTodoListDAO{
					tasks: []model.Task{
						{
							ID:        1,
							Name:      "task1",
							CreatedBy: "1234",
						},
					},
				},
			},
			args: args{
				req: func() *http.Request {
					return httptest.NewRequest(http.MethodGet, "http://www.google.com/1234", strings.NewReader(""))
				}(),
			},
			status: http.StatusOK,
			body: []*model.Task{
				{
					ID:        1,
					Name:      "task1",
					CreatedBy: "1234",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				TodoListDAO: tt.fields.TodoListDAO,
			}
			writer := httptest.NewRecorder()
			handler := h.FetchByID()
			handler.ServeHTTP(writer, tt.args.req)

			if writer.Result().StatusCode != tt.status {
				t.Errorf("Handler.FetchByID() = %v, want %v", writer.Result().StatusCode, tt.status)
				return
			}

			var bodyMap map[string]interface{}
			if err := json.NewDecoder(writer.Body).Decode(&bodyMap); err != nil {
				t.Errorf("Handler.FetchByID() = json body decode error %v", err)
				return
			}

			var wantBodyMap map[string]interface{}
			if enc, err := json.Marshal(tt.body); err == nil {
				_ = json.Unmarshal(enc, &wantBodyMap)
			}

		})
	}
}

func TestHandler_Update(t *testing.T) {
	type fields struct {
		TodoListDAO DAO
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		status int
		body   interface{}
	}{
		{
			name: "updated",
			fields: fields{
				TodoListDAO: &mockTodoListDAO{
					tasks: []model.Task{
						{
							ID:        1,
							Name:      "task1",
							CreatedBy: "1234",
						},
					},
				},
			},
			args: args{
				req: func() *http.Request {
					u := model.Task{
						ID:   1,
						Name: "task1",
					}
					enc, _ := json.Marshal(u)
					return httptest.NewRequest(http.MethodPatch, "http://www.google.com/1234", bytes.NewReader(enc))
				}(),
			},
			status: http.StatusOK,
			body: model.Task{
				ID:   1,
				Name: "task1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				TodoListDAO: tt.fields.TodoListDAO,
			}
			writer := httptest.NewRecorder()
			handler := h.Update()
			handler.ServeHTTP(writer, tt.args.req)

			if writer.Result().StatusCode != tt.status {
				t.Errorf("Handler.Update() = %v, want %v", writer.Result().StatusCode, tt.status)
				return
			}

			var bodyMap map[string]interface{}
			if err := json.NewDecoder(writer.Body).Decode(&bodyMap); err != nil {
				t.Errorf("Handler.Update() = json body decode error %v", err)
				return
			}

			var wantBodyMap map[string]interface{}
			if enc, err := json.Marshal(tt.body); err == nil {
				_ = json.Unmarshal(enc, &wantBodyMap)
			}

			if !reflect.DeepEqual(bodyMap, wantBodyMap) {
				t.Errorf("Handler.Update() = %v, want %v", bodyMap, wantBodyMap)
			}
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	type fields struct {
		TodoListDAO DAO
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		status int
	}{
		{
			name: "deleted",
			fields: fields{
				TodoListDAO: &mockTodoListDAO{
					tasks: []model.Task{
						{
							ID:        1,
							Name:      "task1",
							CreatedBy: "1234",
						},
					},
				},
			},
			args: args{
				req: httptest.NewRequest(http.MethodDelete, "http://www.google.com/1", nil),
			},
			status: http.StatusNoContent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				TodoListDAO: tt.fields.TodoListDAO,
			}
			writer := httptest.NewRecorder()
			handler := h.Delete()
			handler.ServeHTTP(writer, tt.args.req)

			if writer.Result().StatusCode != tt.status {
				t.Errorf("Handler.Delete() = %v, want %v", writer.Result().StatusCode, tt.status)
				return
			}
		})
	}
}
