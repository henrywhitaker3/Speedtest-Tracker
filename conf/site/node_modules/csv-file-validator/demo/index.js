import CSVFileValidator from '../src/csv-file-validator'

const requiredError = (headerName, rowNumber, columnNumber) => {
    return `<div class="red">${headerName} is required in the <strong>${rowNumber} row</strong> / <strong>${columnNumber} column</strong></div>`
}
const validateError = (headerName, rowNumber, columnNumber) => {
    return `<div class="red">${headerName} is not valid in the <strong>${rowNumber} row</strong> / <strong>${columnNumber} column</strong></div>`
}
const uniqueError = (headerName) => {
    return `<div class="red">${headerName} is not unique</div>`
}
const isEmailValid = function (email) {
    const reqExp = /[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,3}$/
    return reqExp.test(email)
}
const isPasswordValid = function (password) {
    return password.length >= 4
}

const CSVConfig = {
    headers: [
        { name: 'First Name', inputName: 'firstName', required: true, requiredError },
        { name: 'Last Name', inputName: 'lastName', required: true, requiredError },
        { name: 'Email', inputName: 'email', required: true, requiredError, unique: true, uniqueError, validate: isEmailValid, validateError },
        { name: 'Password', inputName: 'password', required: true, requiredError, validate: isPasswordValid, validateError },
        { name: 'Roles', inputName: 'roles', required: true, requiredError, isArray: true }
    ]
}

document.getElementById('file').onchange = function(event) {
    CSVFileValidator(event.target.files[0], CSVConfig)
        .then(csvData => {
            csvData.inValidMessages.forEach(message => {
                document.getElementById('invalidMessages').insertAdjacentHTML('beforeend', message)
            })
            console.log(csvData.inValidMessages)
            console.log(csvData.data)
        })
}
