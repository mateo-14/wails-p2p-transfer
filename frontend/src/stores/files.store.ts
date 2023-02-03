import { create } from 'zustand';
import { data } from '../../wailsjs/go/models';

type FilesState = {
  files: data.File[];
  addFiles: (files: data.File[]) => void;
  setFiles: (files: data.File[]) => void;
  removeFile: (id: data.PeerFile['id']) => void;
};

export const useFilesStore = create<FilesState>(set => ({
  files: [],
  addFiles: (files: data.File[]) => set(state => ({ files: [...state.files, ...files] })),
  setFiles: (files: data.File[]) => set({ files }),
  removeFile: (id: data.PeerFile['id']) => set(state => ({ files: state.files.filter(file => file.id !== id) })),
}));
