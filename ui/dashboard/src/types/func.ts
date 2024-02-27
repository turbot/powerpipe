export type AsyncAction<T> = () => Promise<T>;
export type AsyncNoop = () => Promise<void>;
export type Noop = () => void;
