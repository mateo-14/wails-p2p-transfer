import { AddFiles, RemoveSharedFile } from '../../wailsjs/go/main/App';
import { data } from '../../wailsjs/go/models';
import { EventsOnce } from '../../wailsjs/runtime';

type AddFileCallback = (files: data.File[]) => void;

export function addFiles(cb: AddFileCallback) {
  return AddFiles().then(eventID => {
    if (eventID) {
      EventsOnce(eventID, (cb));
      return true;
    }

    return false;
  });
}

export function removeFile(id: data.PeerFile['id']) {
  return RemoveSharedFile(id);
}
