import { Options, Modifier, Instance, VirtualElement } from "./types";
export * from "./types";
export * from "./enums";
declare type PopperGeneratorArgs = {
    defaultModifiers?: Array<Modifier<any>>;
    defaultOptions?: Partial<Options>;
};
export declare function popperGenerator(generatorOptions?: PopperGeneratorArgs): (reference: Element | VirtualElement, popper: HTMLElement, options?: Partial<Options>) => Instance;
export declare const createPopper: (reference: Element | VirtualElement, popper: HTMLElement, options?: Partial<Options>) => Instance;
