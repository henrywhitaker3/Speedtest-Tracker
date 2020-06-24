var _uniqBy = require('lodash/uniqBy');

/**
 * Checking if values are unique
 *
 * @customNeeds 
 * For e.g. [{email:'api@test.com'}, {email:'api@test.com'}] - email is not valid
 *
 * @since 1.3.0
 * @category Array
 *
 * @param {Array} array - The array of objects 
 * @param {String} keyName - Name of the object property from an array in which unique will be checking  
 *
 * @returns {Boolean} Returns true if values are unique and false if not
 *
 * @example
 *
 * famulus.isValuesUnique([{email:'api@test.com'}, {email:'api@test.com'}], 'email')
 * // => false
 */
function isValuesUnique(array, keyName) {
    return _uniqBy(array, keyName).length === array.length;
}

module.exports = isValuesUnique;
