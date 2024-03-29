import classNames from 'classnames';
import { useState } from "react";
import { data } from "../../../wailsjs/go/models";
import { DownloadFile } from '../../../wailsjs/go/main/App';

type PeerFilesProps = {
  files: data.PeerFile[];
  peerId: string;
};

export default function PeerFiles({ files, peerId }: PeerFilesProps) {
  const [selectedFile, setSelectedFile] = useState<string | null>(null);

  const handleClick = (file: data.PeerFile) => {
    DownloadFile(peerId, file.id)
  }
  return (
    <div className="h-full grid grid-cols-[1fr_300px]">
      <div
        className={classNames('files-grid overflow-y-auto pr-6', {
          'col-span-full': !selectedFile
        })}
      >
        {files.map(file => (
          <button className="aspect-square rounded-md bg-zinc-900/50 flex items-end justify-center text-xs px-2 py-1 break-all text-center hover:bg-purple-600/30" onClick={() => handleClick(file)} key={file.id}>
            {file.name}
          </button>
        ))}
      </div>
      {selectedFile ? <div className="">Selected file</div> : null}
    </div>
  );
}
