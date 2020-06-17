import test from 'ava';
import substr from './../substr';

test('module should be a function', t => {
    t.is(typeof substr, 'function');
});

test('should substring "Hello"', t => {
    t.is(substr('Hello World!', 0, 5), 'Hello');
});

test('should substring "World!"', t => {
    t.is(substr('Hello World!', 6), 'World!');
});

test('should return value if it is not correct', t => {
    t.deepEqual(substr({}, 0, 5), {});
});
