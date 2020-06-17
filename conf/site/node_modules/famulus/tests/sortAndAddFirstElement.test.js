import test from 'ava';
import sortAndAddFirstElement from './../sortAndAddFirstElement';

test('module should be a function', t => {
    t.is(typeof sortAndAddFirstElement, 'function');
});

test('should sort and add first element', t => {
    t.deepEqual(
        sortAndAddFirstElement([{name:'Bob'}, {name:'Aron'}], 'name', {name:'All'}),
        [ {name:'All'}, {name:'Aron'}, {name:'Bob'} ]
    );
});
