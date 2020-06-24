import test from 'ava';
import dateDifference from '../dateDifference';

test('module should be a function', t => {
    t.is(typeof dateDifference, 'function');
});

test('should return difference in days', t => {
    t.is(dateDifference(new Date('06-20-2018'), new Date('06-26-2018'), 'days'), 6);
});

test('should return difference in hours', t => {
    t.is(dateDifference(new Date('06-20-2018'), new Date('06-26-2018'), 'hours'), 144);
    t.is(dateDifference(new Date('06-26-2018 10:00'), new Date('06-26-2018 08:00'), 'hours'), 2)
});

test('should return difference in minutes', t => {
    t.is(dateDifference(new Date('06-20-2018'), new Date('06-26-2018'), 'minutes'), 8640);
});

test('should return difference in milliseconds', t => {
    t.is(dateDifference(new Date('06-26-2018'), new Date('06-20-2018'), 'milliseconds'), 518400000);
});

test('should return all difference types in object', t => {
    const diffResult = { days: 6, hours: 2, minutes: 10, milliseconds: 526200000 };
    t.deepEqual(dateDifference(new Date('06-26-2018 10:10'), new Date('06-20-2018 08:00'), 'all'), diffResult);
});
