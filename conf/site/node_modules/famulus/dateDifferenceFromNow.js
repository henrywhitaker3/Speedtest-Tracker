var dates = require('./.internal/dates');

/**
 * Difference between now and date which is passed, in formats 'milliseconds', 'days', 'hours', 'minutes'
 *
 * @customNeeds -
 *
 * @since 2.0.0
 * @category Date
 *
 * @param {Date} date - The Date to inspect
 * @param {String} differenceType - [ 'days', 'hours', 'minutes', 'milliseconds', 'all' ]
 *
 * @returns {Number|Object} Returns the numeric value or object depends on passed differenceType param
 *
 * @example
 * example result for now is new Date('12-26-2017')
 *
 * famulus.dateDifferenceFromNow(new Date('12-20-2017'), 'milliseconds')
 * // => 555261242
 *
 * famulus.dateDifferenceFromNow(new Date('12-20-2017'), 'days')
 * // => 6
 *
 * famulus.dateDifferenceFromNow(new Date('12-20-2017'), 'hours')
 * // => 156
 */
function dateDifferenceFromNow(date, differenceType) {
    var now = new Date(),
        diffMilliseconds = Math.abs(date - now);

    switch(differenceType) {
        case 'days':
            return dates._getDaysDiff(diffMilliseconds);
        case 'hours':
            return dates._differenceInHours(diffMilliseconds);
        case 'minutes':
            return dates._differenceInMinutes(diffMilliseconds);
        case 'milliseconds':
            return diffMilliseconds;

        default:
            return {
                days: dates._getDaysDiff(diffMilliseconds),
                hours: dates._getHoursDiff(diffMilliseconds),
                minutes: dates._getMinutesDiff(diffMilliseconds),
                milliseconds: diffMilliseconds
            }
    }
}

module.exports = dateDifferenceFromNow;
