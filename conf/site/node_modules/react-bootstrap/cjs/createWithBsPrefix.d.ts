import React from 'react';
import { BsPrefixRefForwardingComponent } from './helpers';
interface BsPrefixOptions {
    displayName?: string;
    Component?: React.ElementType;
    defaultProps?: any;
}
export default function createWithBsPrefix(prefix: string, { displayName, Component, defaultProps, }?: BsPrefixOptions): BsPrefixRefForwardingComponent<any>;
export {};
