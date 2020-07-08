import React from 'react';
import { Variant } from './types';
import { TransitionType } from './helpers';
export interface AlertProps extends React.HTMLProps<HTMLDivElement> {
    bsPrefix?: string;
    variant?: Variant;
    dismissible?: boolean;
    show?: boolean;
    onClose?: (a: any, b: any) => void;
    closeLabel?: string;
    transition?: TransitionType;
}
declare const AlertHeading: import("./helpers").BsPrefixRefForwardingComponent<any, {}>;
declare const AlertLink: import("./helpers").BsPrefixRefForwardingComponent<any, {}>;
declare type Alert = React.ForwardRefExoticComponent<AlertProps> & {
    Link: typeof AlertLink;
    Heading: typeof AlertHeading;
};
declare const Alert: Alert;
export default Alert;
