var _ = require('lodash');

/**
 * Sort an array by the name of existing property and add a first element into array
 *
 * @customNeeds Sort array by name and add first element.
 * For e.g. user names - [ {name: 'All'}, {name: 'Aron'}, {name: 'Bob'} ]
 *
 * @since 1.1.0
 * @category Array
 *
 * @param {Array} array - The array to sort and add
 * @param {String} sortBy - Name of the property from an array by which array will be sorted
 * @param {*} element - Element which is added into an array
  *
 * @returns {Array} Returns the new array
 *
 * @example
 *
 * famulus.sortAndAddFirstElement([{name:'Bob'}, {name:'Aron'}], 'name', {name:'All'})
 * // => [ {name:'All'}, {name:'Aron'}, {name:'Bob'} ]
 */
function sortAndAddFirstElement(array, sortBy, element) {
    return _(array)
        .sortBy(sortBy)
        .unshift(element)
        .value();
}

module.exports = sortAndAddFirstElement;
