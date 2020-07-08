/// <reference types="react" />
import FigureImage from './FigureImage';
import FigureCaption from './FigureCaption';
declare type Figure = React.ComponentType & {
    Image: typeof FigureImage;
    Caption: typeof FigureCaption;
};
declare const Figure: Figure;
export default Figure;
