package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ibrahimhaykal/otomedia-task-manager/internal/model"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/repository"
	"github.com/ibrahimhaykal/otomedia-task-manager/internal/service"
)

type mockTaskRepository struct {
	findAllFunc    func(ctx context.Context, filter model.TaskFilter) ([]model.Task, int64, error)
	findByIDFunc   func(ctx context.Context, id string) (*model.Task, error)
	createFunc     func(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error)
	updateFunc     func(ctx context.Context, id string, req model.UpdateTaskRequest) (*model.Task, error)
	softDeleteFunc func(ctx context.Context, id string) error
	findAllCalls   int
}

func (m *mockTaskRepository) FindAll(ctx context.Context, filter model.TaskFilter) ([]model.Task, int64, error) {
	m.findAllCalls++
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx, filter)
	}
	return nil, 0, nil
}

func (m *mockTaskRepository) FindByID(ctx context.Context, id string) (*model.Task, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockTaskRepository) Create(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockTaskRepository) Update(ctx context.Context, id string, req model.UpdateTaskRequest) (*model.Task, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, req)
	}
	return nil, nil
}

func (m *mockTaskRepository) SoftDelete(ctx context.Context, id string) error {
	if m.softDeleteFunc != nil {
		return m.softDeleteFunc(ctx, id)
	}
	return nil
}

type mockTaskCache struct {
	data            map[string]*model.TaskListResponse
	invalidateCalls int
}

func newMockTaskCache() *mockTaskCache {
	return &mockTaskCache{data: make(map[string]*model.TaskListResponse)}
}

func (m *mockTaskCache) GetTaskList(ctx context.Context, key string) (*model.TaskListResponse, error) {
	if res, ok := m.data[key]; ok {
		return res, nil
	}
	return nil, nil
}

func (m *mockTaskCache) SetTaskList(ctx context.Context, key string, data *model.TaskListResponse, ttl time.Duration) error {
	m.data[key] = data
	return nil
}

func (m *mockTaskCache) InvalidateTaskLists(ctx context.Context) error {
	m.invalidateCalls++
	m.data = make(map[string]*model.TaskListResponse)
	return nil
}

func TestCreateTask_ValidationError(t *testing.T) {
	svc := service.NewTaskService(&mockTaskRepository{}, newMockTaskCache())

	_, err := svc.CreateTask(context.Background(), model.CreateTaskRequest{
		Title: "   ",
	})

	var valErr *service.ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("expected ValidationError for empty title, got %v", err)
	}
}

func TestCreateTask_SuccessAndInvalidateCache(t *testing.T) {
	mockCache := newMockTaskCache()
	mockRepo := &mockTaskRepository{
		createFunc: func(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error) {
			if req.Status != model.StatusTodo {
				t.Errorf("expected default status 'todo', got %s", req.Status)
			}
			return &model.Task{
				ID:     "11111111-1111-4111-8111-111111111199",
				Title:  req.Title,
				Status: req.Status,
			}, nil
		},
	}

	svc := service.NewTaskService(mockRepo, mockCache)
	created, err := svc.CreateTask(context.Background(), model.CreateTaskRequest{
		Title: "  New Task  ",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.Title != "New Task" {
		t.Errorf("expected trimmed title 'New Task', got '%s'", created.Title)
	}
	if mockCache.invalidateCalls != 1 {
		t.Errorf("expected cache invalidation to be called once, got %d", mockCache.invalidateCalls)
	}
}

func TestListTasks_CacheHitAndMiss(t *testing.T) {
	mockCache := newMockTaskCache()
	mockRepo := &mockTaskRepository{
		findAllFunc: func(ctx context.Context, filter model.TaskFilter) ([]model.Task, int64, error) {
			return []model.Task{{ID: "11111111-1111-4111-8111-111111111101", Title: "Task 1"}}, 15, nil
		},
	}

	svc := service.NewTaskService(mockRepo, mockCache)
	filter := model.TaskFilter{Page: 1, Limit: 10}

	// first call repo with cache miss
	res1, hit1, err := svc.ListTasks(context.Background(), filter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hit1 {
		t.Error("expected cache MISS on first call")
	}
	if res1.Pagination.TotalPages != 2 {
		t.Errorf("expected 2 total pages for 15 items with limit 10, got %d", res1.Pagination.TotalPages)
	}
	if mockRepo.findAllCalls != 1 {
		t.Errorf("expected repo FindAll to be called once, got %d", mockRepo.findAllCalls)
	}

	// second call should hit cache
	_, hit2, err := svc.ListTasks(context.Background(), filter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hit2 {
		t.Error("expected cache HIT on second call")
	}
	if mockRepo.findAllCalls != 1 {
		t.Errorf("expected repo FindAll calls to remain 1, got %d", mockRepo.findAllCalls)
	}
}

func TestCreateTask_PropagatesDuplicateError(t *testing.T) {
	mockRepo := &mockTaskRepository{
		createFunc: func(ctx context.Context, req model.CreateTaskRequest) (*model.Task, error) {
			return nil, repository.ErrDuplicateTask
		},
	}

	svc := service.NewTaskService(mockRepo, newMockTaskCache())
	_, err := svc.CreateTask(context.Background(), model.CreateTaskRequest{Title: "Existing Title"})

	if !errors.Is(err, repository.ErrDuplicateTask) {
		t.Fatalf("expected ErrDuplicateTask, got %v", err)
	}
}

func TestUpdateTask_SuccessAndInvalidateCache(t *testing.T) {
	mockCache := newMockTaskCache()
	mockRepo := &mockTaskRepository{
		updateFunc: func(ctx context.Context, id string, req model.UpdateTaskRequest) (*model.Task, error) {
			return &model.Task{
				ID:     id,
				Title:  req.Title,
				Status: req.Status,
			}, nil
		},
	}

	svc := service.NewTaskService(mockRepo, mockCache)
	updated, err := svc.UpdateTask(context.Background(), "11111111-1111-4111-8111-111111111101", model.UpdateTaskRequest{
		Title:  "  Updated Title  ",
		Status: "DONE", // testing case insensitivity
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Title != "Updated Title" || updated.Status != model.StatusDone {
		t.Errorf("unexpected updated task values: %+v", updated)
	}
	if mockCache.invalidateCalls != 1 {
		t.Errorf("expected cache invalidation once, got %d", mockCache.invalidateCalls)
	}
}

func TestDeleteTask_NotFoundDoesNotInvalidateCache(t *testing.T) {
	mockCache := newMockTaskCache()
	mockRepo := &mockTaskRepository{
		softDeleteFunc: func(ctx context.Context, id string) error {
			return repository.ErrTaskNotFound
		},
	}

	svc := service.NewTaskService(mockRepo, mockCache)
	err := svc.DeleteTask(context.Background(), "11111111-1111-4111-8111-111111111999")

	if !errors.Is(err, repository.ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
	if mockCache.invalidateCalls != 0 {
		t.Errorf("expected 0 cache invalidation calls on failed delete, got %d", mockCache.invalidateCalls)
	}
}
