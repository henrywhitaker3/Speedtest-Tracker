import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import Navbar from '../Navbar';

export default class SettingsIndex extends Component {
    constructor(props) {
        super(props)

        this.state = {

        }
    }

    render() {
        return (
            <div>
                <Navbar />
            </div>
        );
    }
}

if (document.getElementById('settingsIndex')) {
    ReactDOM.render(<SettingsIndex />, document.getElementById('settingsIndex'));
}
