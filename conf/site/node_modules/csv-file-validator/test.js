import test from 'ava';
import CSVFileValidator from './src/csv-file-validator';

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

const CSVHeader = CSVConfig.headers.map(i => i.name).join(';');

const CSVInvalidFile = [
	CSVHeader,
	'Vasyl;Stokolosa;v.stokol@gmail.com;123;admin,manager',
	'Vasyl_2;"";v.stokol@gmail.com;123123123;user',
].join('\n');

const CSVValidFile = [
	CSVHeader,
	'Vasyl;Stokolosa;v.stokol@gmail.com;123123;admin,manager',
	'Vasyl;Stokolosa;fake@test.com;123123123;user;Ukraine',
].join('\n');

const CSVValidFileWithoutHeaders = [
	'Vasyl;Stokolosa;v.stokol@gmail.com;123123;admin,manager',
	'Vasyl;Stokolosa;fake@test.com;123123123;user;Ukraine',
].join('\n');

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

test('should validate .csv file without headers and return data, file is valid and headers are optional', async t => {
	const csvData = await CSVFileValidator(CSVValidFileWithoutHeaders, { ...CSVConfig, isHeaderNameOptional: true });

	t.is(csvData.inValidMessages.length, 0);
	t.is(csvData.data.length, 2);
});

test('should validate .csv file with headers and return data, file is valid and headers are optional', async t => {
	const csvData = await CSVFileValidator(CSVValidFile, { ...CSVConfig, isHeaderNameOptional: true });

	t.is(csvData.inValidMessages.length, 0);
	t.is(csvData.data.length, 2);
});

test('should validate .csv file without headers and return invalid messages with data, file is valid and headers are missed', async t => {
	const csvData = await CSVFileValidator(CSVValidFileWithoutHeaders, CSVConfig);

	t.is(csvData.inValidMessages.length, 5);
	t.is(csvData.data.length, 1);
});

test('should return optional column', async t => {
	const csvData = await CSVFileValidator(CSVValidFile, CSVConfig);

	t.is(csvData.data[1].country, 'Ukraine');
});
