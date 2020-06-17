var isString = require('lodash/isString');

/**
 * Extracts parts of a string, beginning at the character at the specified position,
 * and returns the specified number of characters.
 * The substr() does not change the original string.
 *
 * @customNeeds Validate string type for preventing SyntaxError
 *
 * @since 1.0.0
 * @category String
 *
 * @param {String} string -  The string to extract
 * @param {Number} start -  The position where to start the extraction. First character is at index 0
 * @param {Number?} length -  Optional. The number of characters to extract. If omitted, it extracts the rest of the string
 *
 * @returns {String} Returns extract part of a string
 *
 * @example
 *
 * famulus.substr('Hello World!', 0, 5)
 * // => 'Hello'
 *
 * famulus.substr('Hello World!', 6)
 * // => 'World!'
 */
function substr(string, start, length) {
    if (!isString(string)) {
        return string;
    }

    return string.substr(start, length);
}

module.exports = substr;
