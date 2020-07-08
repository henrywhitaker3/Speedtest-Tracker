import ToastHeader from './ToastHeader';
import { BsPrefixPropsWithChildren, BsPrefixRefForwardingComponent, TransitionComponent } from './helpers';
export interface ToastProps extends BsPrefixPropsWithChildren {
    animation?: boolean;
    autohide?: boolean;
    delay?: number;
    onClose?: () => void;
    show?: boolean;
    transition?: TransitionComponent;
}
declare const _default: BsPrefixRefForwardingComponent<"div", ToastProps> & {
    Body: BsPrefixRefForwardingComponent<any, {}>;
    Header: ToastHeader;
};
export default _default;
