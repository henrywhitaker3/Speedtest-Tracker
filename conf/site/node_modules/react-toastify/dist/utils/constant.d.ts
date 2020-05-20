export declare type TypeOptions = 'info' | 'success' | 'warning' | 'error' | 'default';
export declare type Positons = 'top-left' | 'top-right' | 'top-center' | 'bottom-left' | 'bottom-right' | 'bottom-center';
export declare const POSITION: {
    readonly TOP_LEFT: "top-left";
    readonly TOP_RIGHT: "top-right";
    readonly TOP_CENTER: "top-center";
    readonly BOTTOM_LEFT: "bottom-left";
    readonly BOTTOM_RIGHT: "bottom-right";
    readonly BOTTOM_CENTER: "bottom-center";
};
export declare const TYPE: {
    INFO: string;
    SUCCESS: string;
    WARNING: string;
    ERROR: string;
    DEFAULT: string;
    DARK: string;
};
export declare const enum DEFAULT {
    COLLAPSE_DURATION = 300,
    DEBOUNCE_DURATION = 50,
    CSS_NAMESPACE = "Toastify"
}
