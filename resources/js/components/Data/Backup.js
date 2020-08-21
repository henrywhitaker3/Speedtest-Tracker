import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Button, Dropdown, DropdownButton } from 'react-bootstrap';
import { toast } from 'react-toastify';
import Axios from 'axios';

export default class Backup extends Component {
    backup = (format) => {
        var url = 'api/backup?format=' + format + '&token=' + window.token;

        toast.info('Your backup has started downloading...');

        Axios.get(url, {
            responseType: 'blob'
        })
        .then((resp) => {
            var a = document.createElement('a');
            a.href = url;
            a.download = "";
            document.body.appendChild(a);
            a.click();
            a.remove();
            toast.success('Backup downloaded');
        })
        .catch((err) => {
            console.log(err);
        })
    }

    render() {
        return (
            <DropdownButton title="Backup" variant="primary" className="m-2 d-inline-block">
                <Dropdown.Item href="#" onClick={() => { this.backup('json') }}>JSON</Dropdown.Item>
                <Dropdown.Item href="#" onClick={() => { this.backup('csv') }}>CSV</Dropdown.Item>
            </DropdownButton>
        );
    }
}

if (document.getElementById('Backup')) {
    ReactDOM.render(<Backup />, document.getElementById('Backup'));
}
