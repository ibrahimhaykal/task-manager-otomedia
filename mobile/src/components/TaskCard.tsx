import React from 'react';
import { Text, TouchableOpacity, View } from 'react-native';
import { taskCardStyles as styles } from '../styles/taskCard.styles';
import { Task, TaskStatus } from '../types/task';

interface TaskCardProps {
  task: Task;
  onEdit: (task: Task) => void;
  onDelete: (task: Task) => void;
}

const STATUS_CONFIG: Record<TaskStatus, { label: string; bgColor: string; textColor: string }> = {
  todo: {
    label: 'To Do',
    bgColor: '#F1F5F9',
    textColor: '#475569',
  },
  in_progress: {
    label: 'In Progress',
    bgColor: '#FEF3C7',
    textColor: '#B45309',
  },
  done: {
    label: 'Done',
    bgColor: '#DCFCE7',
    textColor: '#15803D',
  },
};

function formatDate(isoString: string): string {
  const date = new Date(isoString);
  if (Number.isNaN(date.getTime())) return '-';
  return date.toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  });
}

export function TaskCard({ task, onEdit, onDelete }: TaskCardProps) {
  const statusStyle = STATUS_CONFIG[task.status] || STATUS_CONFIG.todo;

  return (
    <View style={styles.card}>
      <View style={styles.headerRow}>
        <Text style={styles.title} numberOfLines={2}>
          {task.title}
        </Text>
        <View style={[styles.badge, { backgroundColor: statusStyle.bgColor }]}>
          <Text style={[styles.badgeText, { color: statusStyle.textColor }]}>
            {statusStyle.label}
          </Text>
        </View>
      </View>

      {task.description ? (
        <Text style={styles.description} numberOfLines={3}>
          {task.description}
        </Text>
      ) : null}

      <View style={styles.metaRow}>
        <Text style={styles.metaText}>
          Assignee: <Text style={styles.metaBold}>{task.assignee || 'Unassigned'}</Text>
        </Text>
        <Text style={styles.metaText}>{formatDate(task.created_at)}</Text>
      </View>

      <View style={styles.actionRow}>
        <TouchableOpacity
          style={[styles.actionBtn, styles.editBtn]}
          onPress={() => onEdit(task)}
          activeOpacity={0.7}
        >
          <Text style={styles.editBtnText}>Edit</Text>
        </TouchableOpacity>

        <TouchableOpacity
          style={[styles.actionBtn, styles.deleteBtn]}
          onPress={() => onDelete(task)}
          activeOpacity={0.7}
        >
          <Text style={styles.deleteBtnText}>Hapus</Text>
        </TouchableOpacity>
      </View>
    </View>
  );
}
