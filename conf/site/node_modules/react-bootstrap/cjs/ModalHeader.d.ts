import React from 'react';
import { BsPrefixAndClassNameOnlyProps } from './helpers';
export interface ModalHeaderProps extends React.PropsWithChildren<BsPrefixAndClassNameOnlyProps> {
    closeLabel?: string;
    closeButton?: boolean;
    onHide?: () => void;
}
declare const ModalHeader: React.ForwardRefExoticComponent<ModalHeaderProps & React.RefAttributes<HTMLDivElement>>;
export default ModalHeader;
