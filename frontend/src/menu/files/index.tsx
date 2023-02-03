import Button from '../../components/Button';
import { useFilesStore } from '../../stores/files.store';
import { Item, Menu, useContextMenu } from 'react-contexify';
import { data } from '../../../wailsjs/go/models';
import { addFiles, removeFile } from '../../services/filesService';
import { useState } from 'react';
import classNames from 'classnames';

export default function Files() {
  const store = useFilesStore();
  const [isLoading, setIsLoading] = useState<boolean>(false);

  const handleClick = () => {
    addFiles(files => {
      console.log(files)
      store.addFiles(files);
      setIsLoading(false);
    }).then((hasSelectedFiles) => {
      if (hasSelectedFiles)
        setIsLoading(true);
    })
  };

  return (
    <div className="flex flex-col h-full">
      <div className="bg-zinc-900/25 py-4 px-6 flex items-center justify-between">
        <h2 className="text-xl">You shared files</h2>
        <Button
          onClick={handleClick}
          disabled={isLoading}
          className={classNames({ 'cursor-wait': isLoading })}
        >
          Add file/s
        </Button>
      </div>
      <div className="flex-1 px-6 py-3">
        <ul>
          {store.files.map(file => (
            <FileItem file={file} key={file.id} />
          ))}
        </ul>
      </div>
      {isLoading ? (
        <div className="bg-zinc-900/25 py-2 px-2 text-sm font-semibold cursor-wait">
          <p className="animate-pulse">Hashing and storing files. Wait please.</p>
        </div>
      ) : null}
    </div>
  );
}

type FileItemProps = {
  file: data.File;
};

function FileItem({ file }: FileItemProps) {
  const menuID = `file-${file.id}`;
  const { show } = useContextMenu({ id: menuID });
  const removeFileFromStore = useFilesStore(state => state.removeFile);

  const handleRemove = () => {
    removeFile(file.id).then(() => {
      removeFileFromStore(file.id);
    });
  };

  return (
    <li
      className="hover:bg-white/20 py-0.5 my-1 px-2 cursor-default rounded-md"
      onContextMenu={e => show({ event: e })}
    >
      {file.name}
      <Menu id={menuID}>
        <Item onClick={handleRemove}>Remove</Item>
      </Menu>
    </li>
  );
}
