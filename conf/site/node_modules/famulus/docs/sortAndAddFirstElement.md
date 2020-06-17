## sortAndAddFirstElement(array, sortBy, element)
[Source](../sortAndAddFirstElement.js)

Sort an array by the name of existing property and add a first element into array

#### Custom Needs
Sort array by name and add first element<br>
For e.g. user names - [ {name: 'All'}, {name: 'Aron'}, {name: 'Bob'} ]

#### Since
1.1.0

#### Category 
Array

#### Arguments
{Array} array    - The array to sort and add<br>
{String} sortBy  - Name of the property from an array by which array will be sorted<br>
{*Any} element   - Element which is added into an array
 
#### Returns
{Array} Returns the new array

#### Example
```javascript
sortAndAddFirstElement([{name:'Bob'}, {name:'Aron'}], 'name', {name:'All'});
// => [ {name:'All'}, {name:'Aron'}, {name:'Bob'} ]
```
