import Button from '../../components/Button';
import { AddFiles } from '../../../wailsjs/go/main/App';

export default function Files() {
  return (
    <div className="flex flex-col h-full">
      <div className="bg-zinc-900/25 py-4 px-6 flex items-center justify-between">
        <h2 className="text-xl">You shared files</h2>
        <Button onClick={() => AddFiles()}>Add file/s</Button>
      </div>
      <div className="flex-1 px-6 py-3"></div>
    </div>
  );
}
