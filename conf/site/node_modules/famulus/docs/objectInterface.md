## objectInterface(config)(object)
[Source](../objectInterface.js)

Interface for building an object by configuration

#### Custom Needs
Have an interface for building an object based on configuration

#### Since
2.1.0

#### Category 
Object

#### Arguments
{Array} config   -  Keys with configuration<br>
```
[
    'key/value' - "OR" if no value, set value after "/",
    'key|this.firstName + " " + this.lastName' - set value from the expression after "|" which is bind to the passed object,
    'key:new Date()' - set value from the expression after ":"
 ]
 ```

#### Returns
{Object}

#### Example
```javascript
var email = objectInterface(['body', 'count/1', 'sender|this.firstName + " " + this.lastName', 'isRead: false', 'created: new Date()'])
// => function

email({body: 'Hello world!', count: '', firstName: 'Vasyl', lastName: 'Stokolosa', another: ''})
// => {body: "Hello world!", count: 1, created: Mon Jul 09 2018 10:31:08, isRead: false, sender: "Vasyl Stokolosa"}
```
