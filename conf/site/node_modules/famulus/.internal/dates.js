var MILLISECONDS_IN_DAY = 86400000;
var MILLISECONDS_IN_HOUR = 3600000;
var MILLISECONDS_IN_MINUTE = 60000;

var HOURS_IN_DAY = 24;
var MINUTES_IN_DAY = 1440;
var MINUTES_IN_HOUR = 60;

/**
 * @param {Number} milliseconds
 * @return {Number}
 * @private
 */
function _getDaysDiff(milliseconds) {
    return Math.abs(Math.floor(milliseconds / MILLISECONDS_IN_DAY));
}

/**
 * @param {Number} milliseconds
 * @return {Number}
 * @private
 */
function _getHoursDiff(milliseconds) {
    return Math.abs(Math.floor((milliseconds % MILLISECONDS_IN_DAY) / MILLISECONDS_IN_HOUR));
}

/**
 * @param {Number} milliseconds
 * @return {Number}
 * @private
 */
function _getMinutesDiff(milliseconds) {
    return Math.abs(Math.round(((milliseconds % MILLISECONDS_IN_DAY) % MILLISECONDS_IN_HOUR) / MILLISECONDS_IN_MINUTE));
}

/**
 * @param {Number} milliseconds
 * @return {Number}
 * @private
 */
function _differenceInHours(milliseconds) {
    var days = _getDaysDiff(milliseconds);

    if (days !== 0) {
        return _getHoursDiff(milliseconds) + (days * HOURS_IN_DAY);
    }

    return _getHoursDiff(milliseconds);
}

/**
 * @param {Number} milliseconds
 * @return {Number}
 * @private
 */
function _differenceInMinutes(milliseconds) {
    var days = _getDaysDiff(milliseconds),
        hours = _getHoursDiff(milliseconds);

    if (days !== 0) {
        days = days * MINUTES_IN_DAY;
    }

    if (hours !== 0) {
        hours = hours * MINUTES_IN_HOUR;
    }

    return _getMinutesDiff(milliseconds) + days + hours;
}

module.exports = {
    _getDaysDiff: _getDaysDiff,
    _getHoursDiff: _getHoursDiff,
    _getMinutesDiff: _getMinutesDiff,
    _differenceInHours: _differenceInHours,
    _differenceInMinutes: _differenceInMinutes
};
