import { create } from 'zustand';

type FilesState = {
  files: File[];
}

export const useFilesStore = create<FilesState>(set => ({files: []}));
