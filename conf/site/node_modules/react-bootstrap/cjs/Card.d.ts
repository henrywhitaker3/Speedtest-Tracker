import CardImg from './CardImg';
import { BsPrefixPropsWithChildren, BsPrefixRefForwardingComponent } from './helpers';
import { Color, Variant } from './types';
declare const CardBody: BsPrefixRefForwardingComponent<any, {}>;
declare const CardTitle: BsPrefixRefForwardingComponent<any, {}>;
declare const CardSubtitle: BsPrefixRefForwardingComponent<any, {}>;
declare const CardLink: BsPrefixRefForwardingComponent<any, {}>;
declare const CardText: BsPrefixRefForwardingComponent<any, {}>;
declare const CardHeader: BsPrefixRefForwardingComponent<any, {}>;
declare const CardFooter: BsPrefixRefForwardingComponent<any, {}>;
declare const CardImgOverlay: BsPrefixRefForwardingComponent<any, {}>;
export interface CardProps extends BsPrefixPropsWithChildren {
    bg?: Variant;
    text?: Color;
    border?: Variant;
    body?: boolean;
}
declare type Card = BsPrefixRefForwardingComponent<'div', CardProps> & {
    Img: typeof CardImg;
    Title: typeof CardTitle;
    Subtitle: typeof CardSubtitle;
    Body: typeof CardBody;
    Link: typeof CardLink;
    Text: typeof CardText;
    Header: typeof CardHeader;
    Footer: typeof CardFooter;
    ImgOverlay: typeof CardImgOverlay;
};
declare const Card: Card;
export default Card;
