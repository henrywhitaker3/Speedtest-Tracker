import { BsPrefixProps, BsPrefixRefForwardingComponent } from './helpers';
declare const MediaBody: BsPrefixRefForwardingComponent<any, {}>;
declare type MediaProps = BsPrefixProps;
declare type Media = BsPrefixRefForwardingComponent<'div', MediaProps> & {
    Body: typeof MediaBody;
};
declare const Media: Media;
export default Media;
