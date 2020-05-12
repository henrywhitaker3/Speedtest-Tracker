/// <reference types="react" />
import { ToastTransitionProps } from '../types';
export interface CSSTransitionProps {
    /**
     * Css class to apply when toast enter
     */
    enter: string;
    /**
     * Css class to apply when toast leave
     */
    exit: string;
    /**
     * Define the duration of the transition in ms
     * You can also pass an array `[enterDuration, exitDuration]`
     * `Default: 750`
     */
    duration?: number | [number, number];
    /**
     * Append current toast position to the classname.
     * For instance `myclass--top-center`...
     * `Default: false`
     */
    appendPosition?: boolean;
    /**
     * Collapse toast smoothly when animation end
     * `Default: true`
     */
    collapse?: boolean;
    /**
     * Collapse transition duration
     * `Default: 300`
     */
    collapseDuration?: number;
}
export declare function cssTransition({ enter, exit, duration, appendPosition, collapse, collapseDuration }: CSSTransitionProps): ({ children, position, preventExitTransition, done, ...props }: ToastTransitionProps) => JSX.Element;
