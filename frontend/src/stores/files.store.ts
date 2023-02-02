import { create } from 'zustand';
import { data } from '../../wailsjs/go/models';

type FilesState = {
  files: data.File[];
  addFiles: (files: data.File[]) => void;
};

export const useFilesStore = create<FilesState>(set => ({
  files: [],
  addFiles: (files: data.File[]) => set(state => ({ files: [...state.files, ...files] }))
}));
