import classnames from 'classnames';

type ButtonProps = {
  onClick?: React.MouseEventHandler<HTMLButtonElement>
  children?: React.ReactNode
  className?: string;
  disabled?: boolean
};

export default function Button({ onClick, children, className, disabled }: ButtonProps) {
  return (
    <button
      className={classnames("bg-purple-700 py-2 px-3 rounded-md text-sm font-semibold hover:bg-purple-600 hover:shadow-lg active:shadow-lg hover:shadow-purple-600/20 active:shadow-purple-600/50 transition-all disabled:bg-zinc-600 disabled:shadow-none", className)}
      onClick={onClick}
      disabled={disabled}
    >
      {children}
    </button>
  );
}
