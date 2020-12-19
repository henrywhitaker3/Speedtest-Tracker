import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Modal, Button, Tab } from 'react-bootstrap';
import Axios from 'axios';
import { toast } from 'react-toastify';
import SettingsInput from '../SettingsInput';

export default class GeneralSettings extends Component {
    constructor(props) {
        super(props)

        this.state = {
            data: this.props.data
        }
    }

    inputHandler = () => {
        var settings = this.state.data;
        var i = 0;
        settings.forEach(ele => {
            if(ele.obj.name == name) {
                ele.obj.value = val;
            }
            settings[i] = ele;
            i++;
        });
        this.setState({
            data: settings
        });
    }

    render() {
        var settings = this.state.data[0];

        return (
            <Tab.Content>
                General
                <SettingsInput
                    name={settings.obj.name}
                    id={settings.obj.id}
                    type={settings.type}
                    value={settings.obj.value}
                    description={settings.obj.description}
                    handler={this.inputHandler}
                    label={settings.obj.name.split('_').join(' ')}
                    description={settings.obj.description}
                />
            </Tab.Content>
        );
    }
}

if (document.getElementById('GeneralSettings')) {
    ReactDOM.render(<GeneralSettings />, document.getElementById('GeneralSettings'));
}
