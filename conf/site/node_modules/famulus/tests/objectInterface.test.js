import test from 'ava';
import objectInterface from './../objectInterface';

const result = {body: "Hello world!", count: 1, sender: "Vasyl Stokolosa", isRead: false, created: 1531119600000};

test('module should be a function', t => {
    t.is(typeof objectInterface, 'function');
});

test('should return object based on interface configuration', t => {
    let email = objectInterface([
        'body', 'count/1', 'sender|this.firstName + " " + this.lastName', 'created: Date.parse("Tue Jul 09 2018 17 GMT-0700")', 'isRead: false'
    ]);

    t.deepEqual(email({
            body: 'Hello world!', count: '', firstName: 'Vasyl', lastName: 'Stokolosa', another: ''
        }), result
    );
});
