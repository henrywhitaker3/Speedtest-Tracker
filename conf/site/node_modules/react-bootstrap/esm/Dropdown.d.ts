import React from 'react';
import DropdownItem from './DropdownItem';
import DropdownMenu from './DropdownMenu';
import DropdownToggle from './DropdownToggle';
import { BsPrefixPropsWithChildren, BsPrefixRefForwardingComponent, SelectCallback } from './helpers';
declare const DropdownHeader: BsPrefixRefForwardingComponent<"div", {}>;
declare const DropdownDivider: BsPrefixRefForwardingComponent<"div", {}>;
export interface DropdownProps extends BsPrefixPropsWithChildren {
    drop?: 'up' | 'left' | 'right' | 'down';
    alignRight?: boolean;
    show?: boolean;
    flip?: boolean;
    onToggle?: (isOpen: boolean, event: React.SyntheticEvent<Dropdown>, metadata: {
        source: 'select' | 'click' | 'rootClose' | 'keydown';
    }) => void;
    focusFirstItemOnShow?: boolean | 'keyboard';
    onSelect?: SelectCallback;
    navbar?: boolean;
}
declare type Dropdown = BsPrefixRefForwardingComponent<'div', DropdownProps> & {
    Toggle: typeof DropdownToggle;
    Menu: typeof DropdownMenu;
    Item: typeof DropdownItem;
    Divider: typeof DropdownDivider;
    Header: typeof DropdownHeader;
};
declare const Dropdown: Dropdown;
export default Dropdown;
