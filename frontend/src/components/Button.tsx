import { ComponentChildren, h } from 'preact';
import classnames from 'classnames';

type ButtonProps = {
  onClick?: h.JSX.MouseEventHandler<HTMLButtonElement>;
  children?: ComponentChildren
  className?: string;
};

export default function Button({ onClick, children, className }: ButtonProps) {
  return (
    <button
      class={classnames("bg-purple-700 py-2 px-3 rounded-md text-sm font-semibold hover:bg-purple-600 hover:shadow-lg active:shadow-lg hover:shadow-purple-600/20 active:shadow-purple-600/50 transition-all", className)}
      onClick={onClick}
    >
      {children}
    </button>
  );
}
