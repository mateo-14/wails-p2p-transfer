import Button from '../../components/Button';
import { AddFiles } from '../../../wailsjs/go/main/App';
import { useFilesStore } from "../../stores/files.store";

export default function Files() {
  const store = useFilesStore();
  
  const handleClick = () => {
    AddFiles().then(
      (files) => {
        if (files) {
          store.addFiles(files);
        }
      }
    )
  }
  return (
    <div className="flex flex-col h-full">
      <div className="bg-zinc-900/25 py-4 px-6 flex items-center justify-between">
        <h2 className="text-xl">You shared files</h2>
        <Button onClick={handleClick}>Add file/s</Button>
      </div>
      <div className="flex-1 px-6 py-3">
        <ul>
          {store.files.map(file => (
            <li key={file.id}>{file.name}</li>
          ))}
        </ul>
      </div>
    </div>
  );
}
