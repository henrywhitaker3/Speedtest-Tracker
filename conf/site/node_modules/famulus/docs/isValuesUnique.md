## isValuesUnique(array, keyName)
[Source](../isValuesUnique.js)

Checking if values are unique

#### Custom Needs
Having an array of objects and checking if values are unqiue by object key

#### Since
1.3.0

#### Category 
Array

#### Arguments
{Array} array   -  The array of objects<br>
{String} keyName    -  Name of the object property from an array in which unique will be checking<br>

#### Returns
{Boolean} Returns true if values are unique and false if not

#### Example
Unique emails
```javascript
isValuesUnique([{email: 'api@test.com'}, {email: 'api@test.com'}], 'email');
// => false
```

Emails are not unique
```javascript
isValuesUnique([{email: 'api@test.com'}, {email: 'api_1@test.com'}], 'email');
// => true
```
