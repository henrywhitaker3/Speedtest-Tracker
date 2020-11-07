import * as React from 'react';
import { OnChangeCallback } from './eventManager';
import { ToastContent, ToastOptions, Id, ToastContainerProps, UpdateOptions, ClearWaitingQueueParams } from '../types';
declare const toast: {
    (content: ToastContent, options?: ToastOptions | undefined): React.ReactText;
    success(content: ToastContent, options?: ToastOptions | undefined): React.ReactText;
    info(content: ToastContent, options?: ToastOptions | undefined): React.ReactText;
    error(content: ToastContent, options?: ToastOptions | undefined): React.ReactText;
    warning(content: ToastContent, options?: ToastOptions | undefined): React.ReactText;
    dark(content: ToastContent, options?: ToastOptions | undefined): React.ReactText;
    /**
     * Maybe I should remove warning in favor of warn, I don't know
     */
    warn: (content: ToastContent, options?: ToastOptions | undefined) => React.ReactText;
    /**
     * Remove toast programmaticaly
     */
    dismiss(id?: string | number | undefined): false | void;
    /**
     * Clear waiting queue when limit is used
     */
    clearWaitingQueue(params?: ClearWaitingQueueParams): false | void;
    /**
     * return true if one container is displaying the toast
     */
    isActive(id: Id): boolean;
    update(toastId: Id, options?: UpdateOptions): void;
    /**
     * Used for controlled progress bar.
     */
    done(id: Id): void;
    /**
     * Track changes. The callback get the number of toast displayed
     *
     */
    onChange(callback: OnChangeCallback): () => void;
    /**
     * Configure the ToastContainer when lazy mounted
     */
    configure(config?: ToastContainerProps): void;
    POSITION: {
        TOP_LEFT: import("../types").ToastPosition;
        TOP_RIGHT: import("../types").ToastPosition;
        TOP_CENTER: import("../types").ToastPosition;
        BOTTOM_LEFT: import("../types").ToastPosition;
        BOTTOM_RIGHT: import("../types").ToastPosition;
        BOTTOM_CENTER: import("../types").ToastPosition;
    };
    TYPE: {
        INFO: import("../types").TypeOptions;
        SUCCESS: import("../types").TypeOptions;
        WARNING: import("../types").TypeOptions;
        ERROR: import("../types").TypeOptions;
        DEFAULT: import("../types").TypeOptions;
        DARK: import("../types").TypeOptions;
    };
};
export { toast };
