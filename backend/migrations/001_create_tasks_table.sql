CREATE DATABASE IF NOT EXISTS task_management;
USE task_management;

CREATE TABLE IF NOT EXISTS tasks (
    id CHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NULL,
    status ENUM('todo', 'in_progress', 'done') NOT NULL DEFAULT 'todo',
    assignee VARCHAR(100) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    active_title VARCHAR(255) GENERATED ALWAYS AS (IF(deleted_at IS NULL, title, NULL)) STORED,
    PRIMARY KEY (id),
    UNIQUE KEY uq_tasks_active_title (active_title),
    INDEX idx_tasks_deleted_status (deleted_at, status),
    INDEX idx_tasks_assignee (assignee),
    INDEX idx_tasks_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO tasks (id, title, description, status, assignee) VALUES
('11111111-1111-4111-8111-111111111101', 'Setup Go Gin Boilerplate', 'Initialize project structure, database connection, and Redis client', 'done', 'Ibrahim'),
('11111111-1111-4111-8111-111111111102', 'Implement Task Filtering & Pagination', 'Support status, keyword, assignee, page, limit, and sort params', 'in_progress', 'Ibrahim'),
('11111111-1111-4111-8111-111111111103', 'Add Redis Cache Invalidation', 'Cache GET /api/tasks for 60s and invalidate on mutations', 'todo', 'Erwin'),
('11111111-1111-4111-8111-111111111104', 'Build React Native Edit Modal', 'Allow updating task fields and refresh list automatically', 'todo', 'Ibrahim');