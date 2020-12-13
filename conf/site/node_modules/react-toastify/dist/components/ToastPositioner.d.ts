import * as React from 'react';
import { ToastProps } from '../types';
declare type ToastPositionerProps = Pick<ToastProps, 'style' | 'in'> & {
    className?: string;
};
export declare const ToastPositioner: React.FC<ToastPositionerProps>;
export {};
