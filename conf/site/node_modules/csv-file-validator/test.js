import test from 'ava';
import CSVFileValidator from './src/csv-file-validator';

const CSVInvalidFile = `
    Vasyl;Stokolosa;v.stokol@gmail.com;123;admin,manager\n
    Vasyl_2;"";v.stokol@gmail.com;123123123;user
`;

const CSVValidFile = `
    Vasyl;Stokolosa;v.stokol@gmail.com;123123;admin,manager\n
    Vasyl;Stokolosa;fake@test.com;123123123;user;Ukraine
`;

const requiredError = (headerName, rowNumber, columnNumber) => (
    `<div class="red">${headerName} is required in the <strong>${rowNumber} row</strong> / <strong>${columnNumber} column</strong></div>`
)

const validateError = (headerName, rowNumber, columnNumber) => (
    `<div class="red">${headerName} is not valid in the <strong>${rowNumber} row</strong> / <strong>${columnNumber} column</strong></div>`
)

const isEmailValid = (email) => {
    const reqExp = /[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,3}$/
    return reqExp.test(email)
}

const isPasswordValid = (password) => (password.length >= 4)
const uniqueError = (headerName) => (`<div class="red">${headerName} is not unique</div>`)

const CSVConfig = {
    headers: [
        { name: 'First Name', inputName: 'firstName', required: true, requiredError },
        { name: 'Last Name', inputName: 'lastName', required: true, requiredError },
        { name: 'Email', inputName: 'email', required: true, requiredError, unique: true, uniqueError, validate: isEmailValid, validateError },
        { name: 'Password', inputName: 'password', required: true, requiredError, validate: isPasswordValid, validateError },
        { name: 'Roles', inputName: 'roles', required: true, requiredError, isArray: true },
        { name: 'Country', inputName: 'country', optional: true }
    ]
}

test('module should be a function', t => {
    t.is(typeof CSVFileValidator, 'function');
});

test('should return an object with empty inValidMessages/data keys', async t => {
    const csvData = await CSVFileValidator('', {});

    t.is(typeof csvData, 'object');
    t.deepEqual(csvData.inValidMessages, []);
    t.deepEqual(csvData.data, []);
});

test('should validate .csv file and return invalid messages with data', async t => {
    const csvData = await CSVFileValidator(CSVInvalidFile, CSVConfig);

    t.is(csvData.inValidMessages.length, 3);
    t.is(csvData.data.length, 2);
});

test('should validate .csv file and return data, file is valid', async t => {
    const csvData = await CSVFileValidator(CSVValidFile, CSVConfig);

    t.is(csvData.inValidMessages.length, 0);
    t.is(csvData.data.length, 2);
});

test('should return optional column', async t => {
    const csvData = await CSVFileValidator(CSVValidFile, CSVConfig);
    
    t.is(csvData.data[1].country, 'Ukraine');
});
