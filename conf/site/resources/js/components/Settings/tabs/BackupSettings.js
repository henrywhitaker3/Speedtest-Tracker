import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Modal, Button, Tab } from 'react-bootstrap';
import Axios from 'axios';
import DataRow from '../../Data/DataRow';

export default class BackupSettings extends Component {
    constructor(props) {
        super(props)

        this.state = {
        }
    }

    render() {
        return (
            <Tab.Content>
                <DataRow />
            </Tab.Content>
        );
    }
}

if (document.getElementById('BackupSettings')) {
    ReactDOM.render(<BackupSettings />, document.getElementById('BackupSettings'));
}
