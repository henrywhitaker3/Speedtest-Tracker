import test from 'ava';
import dateDifferenceFromNow from '../dateDifferenceFromNow';

test('module should be a function', t => {
    t.is(typeof dateDifferenceFromNow, 'function');
});

test('should return Number: days', t => {
    t.is(typeof dateDifferenceFromNow(new Date('12-20-2017'), 'days'), 'number');
});

test('should return Number: hours', t => {
    t.is(typeof dateDifferenceFromNow(new Date('12-20-2017'), 'hours'), 'number');
});

test('should return Number: minutes', t => {
    t.is(typeof dateDifferenceFromNow(new Date('12-20-2017'), 'minutes'), 'number');
});

test('should return Number: milliseconds', t => {
    t.is(typeof dateDifferenceFromNow(new Date('12-20-2017'), 'milliseconds'), 'number');
});

test('should return Object: all', t => {
    t.deepEqual(typeof dateDifferenceFromNow(new Date('12-20-2017'), 'all'), 'object');
});
