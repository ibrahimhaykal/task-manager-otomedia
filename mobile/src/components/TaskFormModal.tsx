import React, { useEffect, useState } from 'react';
import {
  ActivityIndicator,
  KeyboardAvoidingView,
  Modal,
  Platform,
  ScrollView,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';
import { ApiError } from '../api/taskApi';
import { taskFormModalStyles as styles } from '../styles/taskFormModal.styles';
import { Task, TaskPayload, TaskStatus } from '../types/task';

interface TaskFormModalProps {
  visible: boolean;
  initialTask: Task | null;
  existingAssignees: string[];
  onClose: () => void;
  onSubmit: (payload: TaskPayload, taskId?: string) => Promise<void>;
}

const STATUS_OPTIONS: { value: TaskStatus; label: string }[] = [
  { value: 'todo', label: 'To Do' },
  { value: 'in_progress', label: 'In Progress' },
  { value: 'done', label: 'Done' },
];

export function TaskFormModal({
  visible,
  initialTask,
  existingAssignees,
  onClose,
  onSubmit,
}: TaskFormModalProps) {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [status, setStatus] = useState<TaskStatus>('todo');
  const [assignee, setAssignee] = useState('');
  const [assigneeFocused, setAssigneeFocused] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [titleError, setTitleError] = useState<string | null>(null);

  const isEditMode = Boolean(initialTask);

  const filteredAssignees = existingAssignees.filter((name) =>
    name.toLowerCase().includes(assignee.trim().toLowerCase()),
  );

  useEffect(() => {
    if (visible) {
      if (initialTask) {
        setTitle(initialTask.title);
        setDescription(initialTask.description || '');
        setStatus(initialTask.status);
        setAssignee(initialTask.assignee || '');
      } else {
        setTitle('');
        setDescription('');
        setStatus('todo');
        setAssignee('');
      }
      setAssigneeFocused(false);
      setErrorMessage(null);
      setTitleError(null);
    }
  }, [visible, initialTask]);

  const handleSave = async () => {
    const trimmedTitle = title.trim();
    if (!trimmedTitle) {
      setErrorMessage('Judul task wajib diisi.');
      return;
    }

    setSubmitting(true);
    setErrorMessage(null);

    try {
      await onSubmit(
        {
          title: trimmedTitle,
          description: description.trim(),
          status,
          assignee: assignee.trim(),
        },
        initialTask?.id,
      );
      onClose();
    } catch (err) {
      if (err instanceof ApiError && err.status === 409) {
        setTitleError(
          'Judul task sudah digunakan oleh task aktif lain. Gunakan judul yang berbeda.',
        );
      } else if (err instanceof Error) {
        setErrorMessage(err.message);
      } else {
        setErrorMessage('Gagal menyimpan task.');
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Modal visible={visible} animationType="slide" transparent onRequestClose={onClose}>
      <KeyboardAvoidingView
        style={styles.overlay}
        behavior={Platform.OS === 'ios' ? 'padding' : undefined}
      >
        <View style={styles.sheet}>
          <View style={styles.sheetHeader}>
            <Text style={styles.sheetTitle}>{isEditMode ? 'Edit Task' : 'Tambah Task Baru'}</Text>
            <TouchableOpacity onPress={onClose} disabled={submitting}>
              <Text style={styles.closeText}>Tutup</Text>
            </TouchableOpacity>
          </View>

          <ScrollView showsVerticalScrollIndicator={false} keyboardShouldPersistTaps="handled">
            {errorMessage ? (
              <View style={styles.errorBanner}>
                <Text style={styles.errorText}>{errorMessage}</Text>
              </View>
            ) : null}

            <Text style={styles.label}>Judul Task *</Text>
            <TextInput
              style={[styles.input, titleError && styles.inputError]}
              placeholder="Contoh: Integrasi Redis Cache"
              placeholderTextColor="#94A3B8"
              value={title}
              onChangeText={(val) => {
                setTitle(val);
                if (titleError) setTitleError(null);
                if (errorMessage) setErrorMessage(null);
              }}
              editable={!submitting}
            />
            {titleError ? <Text style={styles.fieldErrorText}>{titleError}</Text> : null}

            <Text style={styles.label}>Deskripsi</Text>
            <TextInput
              style={[styles.input, styles.textArea]}
              placeholder="Detail pekerjaan (opsional)"
              placeholderTextColor="#94A3B8"
              value={description}
              onChangeText={setDescription}
              multiline
              numberOfLines={3}
              textAlignVertical="top"
              editable={!submitting}
            />

            <Text style={styles.label}>Status</Text>
            <View style={styles.statusRow}>
              {STATUS_OPTIONS.map((opt) => {
                const selected = status === opt.value;
                return (
                  <TouchableOpacity
                    key={opt.value}
                    style={[styles.statusChip, selected && styles.statusChipSelected]}
                    onPress={() => setStatus(opt.value)}
                    disabled={submitting}
                  >
                    <Text
                      style={[styles.statusChipText, selected && styles.statusChipTextSelected]}
                    >
                      {opt.label}
                    </Text>
                  </TouchableOpacity>
                );
              })}
            </View>

            <Text style={styles.label}>Assignee</Text>
            <View style={styles.assigneeWrapper}>
              <View style={styles.assigneeInputRow}>
                <TextInput
                  style={styles.assigneeInput}
                  placeholder="Pilih atau ketik nama assignee..."
                  placeholderTextColor="#94A3B8"
                  value={assignee}
                  onFocus={() => setAssigneeFocused(true)}
                  onChangeText={(val) => {
                    setAssignee(val);
                    if (!assigneeFocused) setAssigneeFocused(true);
                  }}
                  editable={!submitting}
                />
                {assignee.length > 0 ? (
                  <TouchableOpacity
                    onPress={() => {
                      setAssignee('');
                      setAssigneeFocused(true);
                    }}
                    disabled={submitting}
                  >
                    <Text style={styles.clearAssigneeText}>Hapus</Text>
                  </TouchableOpacity>
                ) : null}
              </View>

              {assigneeFocused && filteredAssignees.length > 0 ? (
                <View style={styles.suggestionBox}>
                  {filteredAssignees.map((name) => {
                    const isExactMatch = name.toLowerCase() === assignee.trim().toLowerCase();
                    return (
                      <TouchableOpacity
                        key={name}
                        style={[styles.suggestionItem, isExactMatch && styles.suggestionItemActive]}
                        onPress={() => {
                          setAssignee(name);
                          setAssigneeFocused(false);
                        }}
                      >
                        <Text
                          style={[
                            styles.suggestionText,
                            isExactMatch && styles.suggestionTextActive,
                          ]}
                        >
                          {name}
                        </Text>
                        <Text style={styles.suggestionHint}>Pilih</Text>
                      </TouchableOpacity>
                    );
                  })}
                </View>
              ) : null}
            </View>

            <View style={styles.footer}>
              <TouchableOpacity
                style={[styles.btn, styles.cancelBtn]}
                onPress={onClose}
                disabled={submitting}
              >
                <Text style={styles.cancelBtnText}>Batal</Text>
              </TouchableOpacity>

              <TouchableOpacity
                style={[styles.btn, styles.saveBtn, submitting && styles.btnDisabled]}
                onPress={handleSave}
                disabled={submitting}
              >
                {submitting ? (
                  <ActivityIndicator color="#FFFFFF" size="small" />
                ) : (
                  <Text style={styles.saveBtnText}>
                    {isEditMode ? 'Simpan Perubahan' : 'Buat Task'}
                  </Text>
                )}
              </TouchableOpacity>
            </View>
          </ScrollView>
        </View>
      </KeyboardAvoidingView>
    </Modal>
  );
}
