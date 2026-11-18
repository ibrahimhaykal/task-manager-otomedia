import { StyleSheet } from 'react-native';

export const taskFormModalStyles = StyleSheet.create({
  overlay: {
    flex: 1,
    backgroundColor: 'rgba(15, 23, 42, 0.5)',
    justifyContent: 'flex-end',
  },
  sheet: {
    backgroundColor: '#FFFFFF',
    borderTopLeftRadius: 20,
    borderTopRightRadius: 20,
    padding: 20,
    maxHeight: '90%',
  },
  sheetHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  sheetTitle: {
    fontSize: 18,
    fontWeight: '700',
    color: '#0F172A',
  },
  closeText: {
    fontSize: 14,
    fontWeight: '600',
    color: '#64748B',
  },
  errorBanner: {
    backgroundColor: '#FEF2F2',
    borderWidth: 1,
    borderColor: '#FECACA',
    borderRadius: 10,
    padding: 12,
    marginBottom: 14,
  },
  errorText: {
    color: '#B91C1C',
    fontSize: 13,
    fontWeight: '500',
    lineHeight: 18,
  },
  label: {
    fontSize: 13,
    fontWeight: '600',
    color: '#334155',
    marginBottom: 6,
  },
  input: {
    backgroundColor: '#F8FAFC',
    borderWidth: 1,
    borderColor: '#CBD5E1',
    borderRadius: 10,
    paddingHorizontal: 12,
    paddingVertical: 10,
    fontSize: 14,
    color: '#0F172A',
    marginBottom: 14,
  },
  textArea: {
    minHeight: 80,
  },
  statusRow: {
    flexDirection: 'row',
    gap: 8,
    marginBottom: 14,
  },
  statusChip: {
    flex: 1,
    paddingVertical: 10,
    borderRadius: 10,
    borderWidth: 1,
    borderColor: '#CBD5E1',
    backgroundColor: '#F8FAFC',
    alignItems: 'center',
  },
  statusChipSelected: {
    backgroundColor: '#2563EB',
    borderColor: '#2563EB',
  },
  statusChipText: {
    fontSize: 13,
    fontWeight: '600',
    color: '#475569',
  },
  statusChipTextSelected: {
    color: '#FFFFFF',
  },
  assigneeWrapper: {
    marginBottom: 14,
  },
  assigneeInputRow: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#F8FAFC',
    borderWidth: 1,
    borderColor: '#CBD5E1',
    borderRadius: 10,
    paddingHorizontal: 12,
  },
  assigneeInput: {
    flex: 1,
    paddingVertical: 10,
    fontSize: 14,
    color: '#0F172A',
  },
  clearAssigneeText: {
    fontSize: 12,
    fontWeight: '600',
    color: '#64748B',
    paddingLeft: 8,
  },
  suggestionBox: {
    marginTop: 4,
    backgroundColor: '#FFFFFF',
    borderWidth: 1,
    borderColor: '#CBD5E1',
    borderRadius: 10,
    overflow: 'hidden',
  },
  suggestionItem: {
    paddingHorizontal: 12,
    paddingVertical: 10,
    borderBottomWidth: 1,
    borderBottomColor: '#F1F5F9',
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  suggestionItemActive: {
    backgroundColor: '#EFF6FF',
  },
  suggestionText: {
    fontSize: 14,
    color: '#334155',
  },
  suggestionTextActive: {
    color: '#1D4ED8',
    fontWeight: '600',
  },
  suggestionHint: {
    fontSize: 11,
    color: '#94A3B8',
  },
  footer: {
    flexDirection: 'row',
    gap: 10,
    marginTop: 8,
    marginBottom: 12,
  },
  btn: {
    flex: 1,
    paddingVertical: 12,
    borderRadius: 10,
    alignItems: 'center',
    justifyContent: 'center',
  },
  cancelBtn: {
    backgroundColor: '#F1F5F9',
  },
  cancelBtnText: {
    color: '#475569',
    fontWeight: '600',
    fontSize: 14,
  },
  saveBtn: {
    backgroundColor: '#2563EB',
  },
  saveBtnText: {
    color: '#FFFFFF',
    fontWeight: '600',
    fontSize: 14,
  },
  btnDisabled: {
    opacity: 0.65,
  },
  inputError: {
    borderColor: '#DC2626',
    backgroundColor: '#FEF2F2',
    marginBottom: 4,
  },
  fieldErrorText: {
    color: '#DC2626',
    fontSize: 12,
    fontWeight: '500',
    marginBottom: 14,
  },
});
