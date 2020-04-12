"use strict";

require(["dojo/node!fs", "dojox/json/ref", "dojo/_base/kernel"], function(fs, ref, kernel) {
    const nodeRequire = kernel.global.require && kernel.global.require.nodeRequire;
    const HIGHLIGHT_DIR = dojo.config.highlightJsDir;
    const CWD = dojo.config.cwd;
    const LANGS_W_DEPS = ['arduino.js'];

    const cloneDeep = nodeRequire(`${CWD}/lodash.cloneDeep.js`);
    const hljs = nodeRequire(`${HIGHLIGHT_DIR}/highlight.js`);

    /**
     * Translate any RegExp objects that may exist in a language definition into a string representation.
     *
     * @param {Object} lang
     * @param {number} nestingLevel
     */
    function regexToStr(lang, nestingLevel = 0) {
        // Max recursion level
        if (nestingLevel > 15) {
            return;
        }

        for (const key in lang) {
            if (lang[key] instanceof RegExp) {
                lang[key] = lang[key].source;
            } else if (typeof lang[key] === 'object') {
                regexToStr(lang[key], nestingLevel + 1);
            }
        }
    }

    /**
     * PCRE does not support the `\uXXXX` syntax, so we must use `\x{XXXX}` instead.
     *
     * @param {string} s
     *
     * @see https://www.regular-expressions.info/unicode.html#codepoint
     *
     * @returns {string}
     */
    function jsUnicodeToPhpUnicode(s) {
        return s.replace(/\\u([0-9A-Fa-f]+)/g, "\\x{$1}");
    }

    /**
     * Load a language and export it as a translated string to STDOUT.
     *
     * @param {string} lang
     */
    function exportLang(lang) {
        const x = nodeRequire(`${HIGHLIGHT_DIR}/languages/${lang}.js`);
        const l = cloneDeep(x(hljs));

        regexToStr(l);
        hljs.registerLanguage(lang, x);

        console.log(lang);
        console.log(jsUnicodeToPhpUnicode(dojox.json.ref.toJson(l)));
    }

    fs.readdir(`${HIGHLIGHT_DIR}/languages/`,function(err, files) {
        if (err) {
            throw err;
        }

        // Load all of the languages that don't extend other languages
        files.forEach(function(file) {
            if (file === ".DS_Store" || LANGS_W_DEPS.indexOf(file) >= 0) {
                return;
            }

            exportLang(file.replace(/\.js$/, ""));
        });

        // These languages extend other languages, so we need to make sure that
        // they are loaded *after* all the standard languages are loaded.
        LANGS_W_DEPS.forEach(function(file) {
            exportLang(file.replace(/\.js$/, ""));
        });
    });
});
