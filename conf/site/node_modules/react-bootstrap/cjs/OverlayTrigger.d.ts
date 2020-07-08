import PropTypes from 'prop-types';
import React from 'react';
import { OverlayChildren, OverlayProps } from './Overlay';
export declare type OverlayTriggerType = 'hover' | 'click' | 'focus';
export interface OverlayTriggerProps extends Omit<OverlayProps, 'children' | 'target'> {
    children: React.ReactNode;
    trigger?: OverlayTriggerType | OverlayTriggerType[];
    delay?: number | {
        show: number;
        hide: number;
    };
    defaultShow?: boolean;
    flip?: boolean;
    overlay: OverlayChildren;
    target?: never;
    onHide?: never;
    show?: never;
}
declare function OverlayTrigger({ trigger, overlay, children, popperConfig, defaultShow, delay: propsDelay, placement, flip, ...props }: OverlayTriggerProps): JSX.Element;
declare namespace OverlayTrigger {
    var propTypes: {
        children: PropTypes.Validator<PropTypes.ReactElementLike>;
        /**
         * Specify which action or actions trigger Overlay visibility
         *
         * @type {'hover' | 'click' |'focus' | Array<'hover' | 'click' |'focus'>}
         */
        trigger: PropTypes.Requireable<string | (string | null | undefined)[]>;
        /**
         * A millisecond delay amount to show and hide the Overlay once triggered
         */
        delay: PropTypes.Requireable<number | PropTypes.InferProps<{
            show: PropTypes.Requireable<number>;
            hide: PropTypes.Requireable<number>;
        }>>;
        /**
         * The initial visibility state of the Overlay. For more nuanced visibility
         * control, consider using the Overlay component directly.
         */
        defaultShow: PropTypes.Requireable<boolean>;
        /**
          The initial flip state of the Overlay.
         */
        flip: PropTypes.Requireable<boolean>;
        /**
         * An element or text to overlay next to the target.
         */
        overlay: PropTypes.Requireable<((...args: any[]) => any) | PropTypes.ReactElementLike>;
        /**
         * A Popper.js config object passed to the the underlying popper instance.
         */
        popperConfig: PropTypes.Requireable<object>;
        /**
         * @private
         */
        target: PropTypes.Requireable<null>;
        /**
         * @private
         */
        onHide: PropTypes.Requireable<null>;
        /**
         * @private
         */
        show: PropTypes.Requireable<null>;
        /**
         * The placement of the Overlay in relation to it's `target`.
         */
        placement: PropTypes.Requireable<string>;
    };
    var defaultProps: {
        defaultShow: boolean;
        trigger: string[];
    };
}
export default OverlayTrigger;
