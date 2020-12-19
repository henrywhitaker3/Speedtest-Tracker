import Axios from 'axios';
import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import Navbar from '../Navbar';
import SettingsTabs from './SettingsTabs';

export default class SettingsIndex extends Component {
    constructor(props) {
        super(props)

        this.state = {
            data: null,
            loading: true,
        }
    }

    getData = () => {
        var url = 'api/settings/?token=' + window.token;

        Axios.get(url)
        .then((resp) => {
            this.setState({
                data: this.sortSettings(resp.data),
                loading: false,
            });
        })
        .catch((err) => {
            //
        })
    }

    sortSettings = (data) => {
        return {
            General: [
                {
                    obj: data.schedule,
                    type: 'text',
                },
                {
                    obj: data.server,
                    type: 'text',
                },
                {
                    obj: data.show_average,
                    type: 'checkbox',
                },
                {
                    obj: data.show_max,
                    type: 'checkbox',
                },
                {
                    obj: data.show_min,
                    type: 'checkbox',
                }
            ],
            Graphs: {}
        };
    }

    componentDidMount() {
        this.getData();
    }

    render() {
        var data = this.state.data;
        var loading = this.state.loading;

        if(loading) {
            return (
                <div>Loading</div>
            );
        }

        return (
            <div>
                <Navbar />
                <div className="container my-5">
                    <SettingsTabs data={data} />
                </div>
            </div>
        );
    }
}

if (document.getElementById('settingsIndex')) {
    ReactDOM.render(<SettingsIndex />, document.getElementById('settingsIndex'));
}
