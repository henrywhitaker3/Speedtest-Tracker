import glob = require("glob");

/**
 * Performs an asynchronous glob search.
 * @param pattern Pattern or patterns to be matched.
 * @param cb The callback invoked when the search completes.
 */
declare function G(pattern: string|string[], cb: (err: Error | null, matches: string[]) => void): void;

/**
 * Performs an asynchronous glob search.
 * @param pattern Pattern or patterns to be matched.
 * @param options The glob options to use.
 * @param cb The callback invoked when the search completes.
 */
declare function G(pattern: string|string[], options: glob.IOptions, cb: (err: Error | null, matches: string[]) => void): void;

declare namespace G {

    /**
     * Performs an synchronous glob search.
     * @param pattern Pattern or patterns to be matched.
     * @param options The glob options to use.
     * @returns The file paths matched by the glob patterns.
     */
    function sync(pattern: string|string[], options?: glob.IOptions): string[];
}

export = G;