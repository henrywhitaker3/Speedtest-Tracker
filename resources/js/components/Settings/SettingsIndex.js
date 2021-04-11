import Axios from 'axios';
import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import Footer from '../Home/Footer';
import Loader from '../Loader';
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
                    obj: data.app_name,
                    type: 'text',
                },
                {
                    obj: data.schedule_enabled,
                    type: 'checkbox',
                },
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
            Tables: [
                {
                    obj: data.visible_columns,
                    type: 'list'
                },
                {
                    obj: data.hidden_columns,
                    type: 'list'
                }
            ],
            Graphs: [
                {
                    obj: data.download_upload_graph_enabled,
                    type: 'checkbox',
                    hideDescription: true
                },
                {
                    obj: data.download_upload_graph_width,
                    type: 'select',
                    options: [
                        {
                            name: 'Full-width',
                            'value': 12
                        },
                        {
                            name: 'Half-width',
                            'value': 6
                        }
                    ],
                },
                {
                    obj: data.ping_graph_enabled,
                    type: 'checkbox',
                    hideDescription: true
                },
                {
                    obj: data.ping_graph_width,
                    type: 'select',
                    options: [
                        {
                            name: 'Full-width',
                            'value': 12
                        },
                        {
                            name: 'Half-width',
                            'value': 6
                        }
                    ],
                },
                {
                    obj: data.failure_graph_enabled,
                    type: 'checkbox',
                    hideDescription: true
                },
                {
                    obj: data.failure_graph_width,
                    type: 'select',
                    options: [
                        {
                            name: 'Full-width',
                            'value': 12
                        },
                        {
                            name: 'Half-width',
                            'value': 6
                        }
                    ],
                },
                {
                    obj: data.show_failed_tests_on_graph,
                    type: 'checkbox',
                },
            ],
            Notifications: [
                {
                    obj: data.slack_webhook,
                    type: 'text'
                },
                {
                    obj: data.telegram_bot_token,
                    type: 'text'
                },
                {
                    obj: data.telegram_chat_id,
                    type: 'text'
                },
                {
                    type: 'btn-get',
                    url: 'api/settings/test-notification?token=' + window.token,
                    btnType: 'primary',
                    obj: {
                        id: (Math.floor(Math.random() * 10000) + 1),
                        name: 'Test notifications',
                        description: 'After saving your updated notification settings, use this to check your settings are correct.'
                    }
                },
                {
                    obj: data.speedtest_notifications,
                    type: 'checkbox'
                },
                {
                    obj: data.speedtest_overview_notification,
                    type: 'checkbox'
                },
                {
                    obj: data.speedtest_overview_time,
                    type: 'number',
                    min: 0,
                    max: 23,
                },
                // Add handling for title stuff
                {
                    obj: data.threshold_alert_percentage,
                    type: 'number',
                    min: 0,
                    max: 100
                },
                {
                    obj: data.threshold_alert_absolute_notifications,
                    type: 'checkbox'
                },
                {
                    obj: data.threshold_alert_absolute_download,
                    type: 'number'
                },
                {
                    obj: data.threshold_alert_absolute_upload,
                    type: 'number'
                },
                {
                    obj: data.threshold_alert_absolute_ping,
                    type: 'number'
                },
            ],
            healthchecks: [
                {
                    obj: data.healthchecks_enabled,
                    type: 'checkbox'
                },
                {
                    obj: data.healthchecks_server_url,
                    type: 'text'
                },
                {
                    obj: data.healthchecks_uuid,
                    type: 'text'
                },
                {
                    obj: {
                        id: (Math.floor(Math.random() * 10000) + 1),
                        name: "Test healthchecks.io integration",
                        description: ""
                    },
                },
                {
                    obj: {
                        id: (Math.floor(Math.random() * 10000) + 1),
                        name: "Start",
                        description: ""
                    },
                    type: 'btn-get',
                    url: 'api/settings/test-healthchecks/start?token=' + window.token,
                    btnType: 'outline-success',
                    inline: true,
                    earlyReturn: true,
                    classes: 'mr-2'
                },
                {
                    obj: {
                        id: (Math.floor(Math.random() * 10000) + 1),
                        name: "Success",
                        description: ""
                    },
                    type: 'btn-get',
                    url: 'api/settings/test-healthchecks/success?token=' + window.token,
                    btnType: 'success',
                    text: 'Success',
                    inline: true,
                    earlyReturn: true,
                    classes: 'mr-2'
                },
                {
                    obj: {
                        id: (Math.floor(Math.random() * 10000) + 1),
                        name: "Fail",
                        description: ""
                    },
                    type: 'btn-get',
                    url: 'api/settings/test-healthchecks/fail?token=' + window.token,
                    btnType: 'danger',
                    text: 'Fail',
                    inline: true,
                    earlyReturn: true,
                    classes: 'mr-2'
                },

            ],
            influxdb: [
                {
                    obj: data.influx_db_enabled,
                    type: 'checkbox'
                },
                {
                    obj: data.influx_db_host,
                    type: 'text'
                },
                {
                    obj: data.influx_db_port,
                    type: 'number'
                },
                {
                    obj: data.influx_db_database,
                    type: 'text'
                },
                {
                    obj: data.influx_db_username,
                    type: 'text',
                    autoComplete: false,
                },
                {
                    obj: data.influx_db_password,
                    type: 'password',
                    autoComplete: false,
                }
            ],
        };
    }

    componentDidMount() {
        this.getData();
    }

    render() {
        var data = this.state.data;
        var loading = this.state.loading;

        return (
            <div>
                <Navbar />
                <div className="container my-5">
                    {loading ?
                        <Loader />
                        :
                        <SettingsTabs data={data} refreshConfig={this.props.refreshConfig} />
                    }
                </div>
                <Footer />
            </div>
        );
    }
}

if (document.getElementById('settingsIndex')) {
    ReactDOM.render(<SettingsIndex />, document.getElementById('settingsIndex'));
}
