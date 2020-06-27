(function (global, factory) {
    typeof exports === 'object' && typeof module !== 'undefined'
        ? module.exports = factory(require('papaparse'), require('lodash/uniqBy'), require('lodash/isFunction'), require('famulus/isValuesUnique'))
        : typeof define === 'function' && define.amd 
            ? define(['papaparse', 'lodash/uniqBy', 'lodash/isFunction', 'famulus/isValuesUnique'], factory)
            : (global.myBundle = factory(global.Papa,global._uniqBy,global._isFunction, global.isValuesUnique));
}(this, (function (Papa, _uniqBy, _isFunction, isValuesUnique) {
    'use strict';

    Papa = Papa && Papa.hasOwnProperty('default') ? Papa['default'] : Papa;
    isValuesUnique = isValuesUnique && isValuesUnique.hasOwnProperty('default') ? isValuesUnique['default'] : isValuesUnique;
    _uniqBy = _uniqBy && _uniqBy.hasOwnProperty('default') ? _uniqBy['default'] : _uniqBy;
    _isFunction = _isFunction && _isFunction.hasOwnProperty('default') ? _isFunction['default'] : _isFunction;

    /**
     * @param {File} csvFile 
     * @param {Object} config 
     */
    function CSVFileValidator(csvFile, config) {
        return new Promise(function(resolve, reject) {
            Papa.parse(csvFile, {
                complete: function(results) {
                    resolve(_prepareDataAndValidateFile(results.data, config));
                },
                error: function(error, file) {
                    reject({ error: error, file: file });
                }
            });
        })
    }

    /**
     * @param {Array} csvData 
     * @param {Object} config 
     * @private
     */
    function _prepareDataAndValidateFile(csvData, config) {
        const file = {
            inValidMessages: [],
            data: []
        };

        csvData.forEach(function(row, rowIndex) {
            const columnData = {};
            const headers = [];

            for (let i = 0; i < config.headers.length; i++) {
                const data = config.headers[i];

                if (!data.optional) {
                    headers.push(data);
                }
            }

            if (row.length < headers.length) {
                return;
            }

            row.forEach(function(columnValue, columnIndex) {
                const valueConfig = config.headers[columnIndex];

                if (!valueConfig) {
                    return;
                }

                // header validation
                if (rowIndex === 0) {
                    if (valueConfig.name !== columnValue) {
                        file.inValidMessages.push(
                            _isFunction(valueConfig.headerError)
                                ? valueConfig.headerError(columnValue)
                                : 'Header name ' + columnValue + ' is not correct or missing'
                        );
                    }

                    return;
                }

                if (valueConfig.required && !columnValue.length) {
                    file.inValidMessages.push(
                        _isFunction(valueConfig.requiredError)
                            ? valueConfig.requiredError(valueConfig.name, rowIndex + 1, columnIndex + 1)
                            : String(valueConfig.name + ' is required in the ' + (rowIndex + 1) + ' row / ' + (columnIndex + 1) + ' column')
                    );
                } else if (valueConfig.validate && !valueConfig.validate(columnValue)) {
                    file.inValidMessages.push(
                        _isFunction(valueConfig.validateError)
                            ? valueConfig.validateError(valueConfig.name, rowIndex + 1, columnIndex + 1)
                            : String(valueConfig.name + ' is not valid in the ' + (rowIndex + 1) + ' row / ' + (columnIndex + 1) + ' column')
                    );
                }

                if (valueConfig.optional) {
                    columnData[valueConfig.inputName] = columnValue;
                }

                if (valueConfig.isArray) {
                    columnData[valueConfig.inputName] = columnValue.split(',').map(function(value) { 
                        return value.trim();
                    });
                } else {
                    columnData[valueConfig.inputName] = columnValue;
                }
            });

            file.data.push(columnData);
        });

        _checkUniqueFields(file, config);

        return file;
    }

    /**
     * @param {Object} file 
     * @param {Object} config
     * @private 
     */
    function _checkUniqueFields(file, config) {
        if (!file.data.length) {
            return;
        }

        config.headers
            .filter(function(header) {
                return header.unique
            })
            .forEach(function(header) {
                if (!isValuesUnique(file.data, header.inputName)) {
                    file.inValidMessages.push(
                        _isFunction(header.uniqueError)
                            ? header.uniqueError(header.name)
                            : String(header.name + ' is not unique')
                    );
                }
            });
    };

    return CSVFileValidator;
})));
