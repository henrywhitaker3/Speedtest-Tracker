## dateDifference(date1, date2, differenceType)
[Source](../dateDifference.js)

Difference between dates which are passed, in formats 'milliseconds', 'days', 'hours', 'minutes'

#### Custom Needs

#### Since
2.0.0

#### Category
Date

#### Arguments
{Date} date1             - The Date for compare<br>
{Date} date2             - The Date for compare<br>
{String} differenceType - [ 'days', 'hours', 'minutes', 'milliseconds', 'all' ]
 
#### Returns
{Number|Object} Returns the numeric value or object depends on passed differenceType param

#### Example

```javascript
dateDifference(new Date('06-20-2018'), new Date('06-26-2018'), 'days')
// => 6
```

```javascript
dateDifference(new Date('06-20-2018'), new Date('06-26-2018'), 'hours')
// => 144
```

```javascript
dateDifference(new Date('06-20-2018'), new Date('06-26-2018'), 'minutes')
// => 8640
```

```javascript
dateDifference(new Date('06-26-2018'), new Date('06-20-2018'), 'milliseconds')
// => 518400000
```

```javascript
dateDifference(new Date('06-26-2018 10:10'), new Date('06-20-2018 08:00'), 'all')
// => {days: 6, hours: 2, minutes: 10, milliseconds: 526200000}
```
