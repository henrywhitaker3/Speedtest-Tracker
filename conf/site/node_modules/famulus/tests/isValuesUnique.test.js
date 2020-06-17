import test from 'ava';
import isValuesUnique from './../isValuesUnique';

const notUniqueData = [{email: 'api@test.com'}, {email: 'api@test.com'}];
const uniqueData = [{email: 'api@test.com'}, {email: 'api_1@test.com'}];

test('module should be a function', t => {
    t.is(typeof isValuesUnique, 'function');
});

test('should return true when values are unique', t => {
    t.is(isValuesUnique(uniqueData, 'email'), true);
});

test('should return false when values are not unique', t => {
    t.is(isValuesUnique(notUniqueData, 'email'), false);
});
