import classNames from 'classnames';
import { useState } from "react";
// import { DownloadFile } from "../../../wailsjs/go/main/App";
import { main } from "../../../wailsjs/go/models";

type PeerFilesProps = {
  files: main.PeerFile[];
  peerId: string;
};

export default function PeerFiles({ files, peerId }: PeerFilesProps) {
  const [selectedFile, setSelectedFile] = useState<string | null>(null);

  const handleClick = (file: main.PeerFile) => {
    // DownloadFile(peerId, file.path)
  }
  return (
    <div className="h-full grid grid-cols-[1fr_300px]">
      <div
        className={classNames('files-grid overflow-y-auto pr-6', {
          'col-span-full': !selectedFile
        })}
      >
        {files.map(file => (
          <button className="aspect-square rounded-md bg-zinc-900/50 flex items-end justify-center text-xs px-2 py-1 break-all text-center hover:bg-purple-600/30 select-none" onClick={() => handleClick(file)} key={file.id}>
            {file.name}
          </button>
        ))}
      </div>
      {selectedFile ? <div className="">Selected file</div> : null}
    </div>
  );
}
