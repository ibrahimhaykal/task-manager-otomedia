package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/handler"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/model"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/repository"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/response"
)

type mockTaskService struct {
	listFunc   func(ctx context.Context, filter model.TaskFilter) (*model.TaskListResponse, bool, error)
	getFunc    func(ctx context.Context, id string) (*model.Task, error)
	createFunc func(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error)
	updateFunc func(ctx context.Context, id string, req model.UpdateTaskRequest) (*model.Task, error)
	deleteFunc func(ctx context.Context, id string) error
}

func (m *mockTaskService) ListTasks(ctx context.Context, filter model.TaskFilter) (*model.TaskListResponse, bool, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, filter)
	}
	return nil, false, nil
}

func (m *mockTaskService) GetTaskByID(ctx context.Context, id string) (*model.Task, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockTaskService) CreateTask(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockTaskService) UpdateTask(ctx context.Context, id string, req model.UpdateTaskRequest) (*model.Task, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, req)
	}
	return nil, nil
}

func (m *mockTaskService) DeleteTask(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func setupTestRouter(svc *mockTaskService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api")
	handler.NewTaskHandler(svc).RegisterRoutes(api)
	return r
}

func TestCreateTask_DuplicateReturns409(t *testing.T) {
	svc := &mockTaskService{
		createFunc: func(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error) {
			return nil, repository.ErrDuplicateTask
		},
	}

	router := setupTestRouter(svc)
	body := bytes.NewBufferString(`{"title":"Setup Go Gin Boilerplate"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/tasks", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409 Conflict, got %d", w.Code)
	}

	var resp response.APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Success || resp.Error == nil || resp.Error.Code != "DUPLICATE_TITLE" {
		t.Errorf("unexpected error payload: %+v", resp.Error)
	}
}

func TestUpdateTask_DuplicateReturns409(t *testing.T) {
	svc := &mockTaskService{
		updateFunc: func(ctx context.Context, id string, req model.UpdateTaskRequest) (*model.Task, error) {
			return nil, repository.ErrDuplicateTask
		},
	}

	router := setupTestRouter(svc)
	body := bytes.NewBufferString(`{"title":"Setup Go Gin Boilerplate","status":"todo"}`)
	req, _ := http.NewRequest(http.MethodPut, "/api/tasks/11111111-1111-4111-8111-111111111102", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409 Conflict, got %d", w.Code)
	}
}

func TestGetTaskByID_InvalidUUIDReturns400(t *testing.T) {
	router := setupTestRouter(&mockTaskService{})
	req, _ := http.NewRequest(http.MethodGet, "/api/tasks/bukan-uuid", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d", w.Code)
	}
}

func TestGetTaskByID_NotFoundReturns404(t *testing.T) {
	svc := &mockTaskService{
		getFunc: func(ctx context.Context, id string) (*model.Task, error) {
			return nil, repository.ErrTaskNotFound
		},
	}

	router := setupTestRouter(svc)
	req, _ := http.NewRequest(http.MethodGet, "/api/tasks/11111111-1111-4111-8111-111111111999", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 Not Found, got %d", w.Code)
	}
}

func TestListTasks_SetsXCacheHeader(t *testing.T) {
	svc := &mockTaskService{
		listFunc: func(ctx context.Context, filter model.TaskFilter) (*model.TaskListResponse, bool, error) {
			return &model.TaskListResponse{
				Tasks:      []model.Task{},
				Pagination: model.PaginationMeta{Page: 1, Limit: 10},
			}, true, nil
		},
	}

	router := setupTestRouter(svc)
	req, _ := http.NewRequest(http.MethodGet, "/api/tasks", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", w.Code)
	}
	if cacheHeader := w.Header().Get("X-Cache"); cacheHeader != "HIT" {
		t.Errorf("expected X-Cache header 'HIT', got '%s'", cacheHeader)
	}
}
