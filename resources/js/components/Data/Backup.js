import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Button } from 'react-bootstrap';
import { toast } from 'react-toastify';
import Axios from 'axios';

export default class Backup extends Component {
    backup = () => {
        var url = '/api/backup';

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
            <Button
                variant="primary"
                className="mx-2"
                onClick={this.backup}
            >Backup</Button>
        );
    }
}

if (document.getElementById('Backup')) {
    ReactDOM.render(<Backup />, document.getElementById('Backup'));
}
