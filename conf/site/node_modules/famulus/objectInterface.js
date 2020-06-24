/**
 * Interface for building an object by configuration
 *
 * @customNeeds - have an interface for building an object based on configuration
 *
 * @since 2.1.0
 * @category Object
 *
 * @param {Array} config - Keys with configuration
 * [
 *   'key/value' - "OR" if no value, set value after "/",
 *   'key|this.firstName + " " + this.lastName' - set value from the expression after "|" which is bind to the passed object,
 *   'key:new Date()' - set value from the expression after ":"
 * ]
 *
 * @returns {Function}
 *
 * @example
 *
 * var email = objectInterface(['body', 'count/1', 'sender|this.firstName + " " + this.lastName', 'isRead: false', 'created: new Date()'])
 * // => function
 *
 * email({body: 'Hello world!', count: '', firstName: 'Vasyl', lastName: 'Stokolosa', another: ''})
 * // => {body: "Hello world!", count: 1, created: Mon Jul 09 2018 10:31:08, isRead: false, sender: "Vasyl Stokolosa"}
 */
function objectInterface(config) {
    return function(obj) {
        var result = {};

        for (var i = 0; i < config.length; i++) {
            var OR, NEXT, REAL;

            if ((OR = config[i].split('/')) && OR[1]) {
                result[OR[0]] = obj[OR[0]] || Function('return ' + OR[1])();
            }
            else if ((NEXT = config[i].split('|')) && NEXT[1]) {
                result[NEXT[0]] = Function('return ' + NEXT[1]).call(obj);
            }
            else if ((REAL = config[i].split(':')) && REAL[1]) {
                result[REAL[0]] = Function('return ' + REAL[1])();
            }
            else {
                result[config[i]] = obj[config[i]];
            }
        }

        return result;
    }
}

module.exports = objectInterface;
