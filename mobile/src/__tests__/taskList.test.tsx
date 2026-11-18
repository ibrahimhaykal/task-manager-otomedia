import { describe, expect, it, jest } from '@jest/globals';
import React from 'react';
import { fireEvent, render } from '@testing-library/react-native';
import { TaskCard } from '../components/TaskCard';
import { Task } from '../types/task';

const mockTask: Task = {
  id: '11111111-1111-4111-8111-111111111101',
  title: 'Setup Go Gin Boilerplate',
  description: 'Initialize project structure and Redis client',
  status: 'in_progress',
  assignee: 'Pak Erwin Otomedia',
  created_at: '2026-11-18T10:00:00Z',
  updated_at: '2026-11-18T10:00:00Z',
};

describe('TaskCard Component', () => {
  it('renders task title, status badge, assignee, and triggers edit/delete callbacks', async () => {
    const onEdit = jest.fn();
    const onDelete = jest.fn();

    const { getByText } = await render(
      <TaskCard task={mockTask} onEdit={onEdit} onDelete={onDelete} />,
    );

    expect(getByText('Setup Go Gin Boilerplate')).toBeTruthy();
    expect(getByText('In Progress')).toBeTruthy();
    expect(getByText('Pak Erwin Otomedia')).toBeTruthy();

    fireEvent.press(getByText('Edit'));
    expect(onEdit).toHaveBeenCalledWith(mockTask);

    fireEvent.press(getByText('Hapus'));
    expect(onDelete).toHaveBeenCalledWith(mockTask);
  });
});
