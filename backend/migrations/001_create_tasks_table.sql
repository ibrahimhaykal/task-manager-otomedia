CREATE DATABASE IF NOT EXISTS task_management;
USE task_management;

CREATE TABLE IF NOT EXISTS tasks (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
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

INSERT INTO tasks (title, description, status, assignee) VALUES
('Setup Go Gin Boilerplate', 'Initialize project structure, database connection, and Redis client', 'done', 'Ibrahim'),
('Implement Task Filtering & Pagination', 'Support status, keyword, assignee, page, limit, and sort params', 'in_progress', 'Ibrahim'),
('Add Redis Cache Invalidation', 'Cache GET /api/tasks for 60s and invalidate on mutations', 'todo', 'Erwin'),
('Build React Native Edit Modal', 'Allow updating task fields and refresh list automatically', 'todo', 'Ibrahim');