import { Placement } from './Overlay';
import { BsPrefixPropsWithChildren, BsPrefixRefForwardingComponent } from './helpers';
export interface TooltipProps extends BsPrefixPropsWithChildren {
    style?: any;
    id: string;
    placement?: Placement;
    arrowProps?: {
        ref: any;
        style: object;
    };
    show?: boolean;
    popper?: any;
}
declare type Tooltip = BsPrefixRefForwardingComponent<'div', TooltipProps>;
declare const Tooltip: Tooltip;
export default Tooltip;
