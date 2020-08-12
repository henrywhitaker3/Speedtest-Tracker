import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Modal, Container, Row, Col, Collapse } from 'react-bootstrap';
import Loader from '../Loader';
import Axios from 'axios';
import Setting from './Setting';
import SettingWithModal from './SettingWithModal';
import ResetSettings from './ResetSettings';

export default class Settings extends Component {
    constructor(props) {
        super(props)

        this.state = {
            show: false,
            loading: true,
            data: [],
        }
    }

    componentDidMount = () => {
        this.getData();
    }

    toggleShow = () => {
        if(this.state.show) {
            var show = false;
        } else {
            var show = true;
        }

        this.setState({
            show: show
        });
    }

    getData = () => {
        var url = 'api/settings/';

        Axios.get(url)
        .then((resp) => {
            this.setState({
                loading: false,
                data: resp.data
            });
        })
        .catch((err) => {
            if(err.response) {

            }
        })
    }

    buildSettingsCards = () => {
        var e = this.state.data;
        return (
            <Row>
                <Col lg={{ span: 4 }} md={{ span: 6 }} sm={{ span: 12 }}>
                    <Setting name={e.schedule.name} value={e.schedule.value} description={e.schedule.description} />
                </Col>
                <Col lg={{ span: 4 }} md={{ span: 6 }} sm={{ span: 12 }}>
                    <Setting name={e.server.name} value={e.server.value} description={e.server.description} />
                </Col>
                <Col lg={{ span: 4 }} md={{ span: 6 }} sm={{ span: 12 }}>
                    <SettingWithModal title="Graph settings" description="Control settings for the graphs." autoClose={true} settings={[
                        {
                            obj: e.download_upload_graph_enabled,
                            type: 'checkbox'
                        },
                        {
                            obj: e.download_upload_graph_width,
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
                            obj: e.ping_graph_enabled,
                            type: 'checkbox'
                        },
                        {
                            obj: e.ping_graph_width,
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
                            obj: e.failure_graph_enabled,
                            type: 'checkbox'
                        },
                        {
                            obj: e.failure_graph_width,
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
                        }
                    ]} />
                </Col>
                <Col lg={{ span: 4 }} md={{ span: 6 }} sm={{ span: 12 }}>
                    <SettingWithModal title="Notification settings" description="Control which types of notifications the server sends." autoClose={false} settings={[
                        {
                            obj: e.slack_webhook,
                            type: 'text'
                        },
                        {
                            obj: e.telegram_bot_token,
                            type: 'text'
                        },
                        {
                            obj: e.telegram_chat_id,
                            type: 'text'
                        },
                        {
                            obj: {
                                id: (Math.floor(Math.random() * 10000) + 1),
                                name: "Test notifications",
                                description: "After saving your updated notification settings, use this to check your settings are correct."
                            },
                            type: 'button-get',
                            url: 'api/settings/test-notification'
                        },
                        {
                            obj: e.speedtest_notifications,
                            type: 'checkbox'
                        },
                        {
                            obj: e.speedtest_overview_notification,
                            type: 'checkbox'
                        },
                        {
                            obj: e.speedtest_overview_time,
                            type: 'number',
                            min: 0,
                            max: 23
                        }
                    ]} />
                </Col>
                <Col lg={{ span: 4 }} md={{ span: 6 }} sm={{ span: 12 }}>
                    <SettingWithModal title="healthchecks.io settings" description="Control settings for healthchecks.io" autoClose={false} settings={[
                        {
                            obj: e.healthchecks_uuid,
                            type: 'text'
                        },
                        {
                            obj: e.healthchecks_enabled,
                            type: 'checkbox'
                        }
                    ]} />
                </Col>
                <Col lg={{ span: 4 }} md={{ span: 6 }} sm={{ span: 12 }}>
                    <ResetSettings />
                </Col>
            </Row>
        )
    }

    render() {
        var show = this.state.show;
        var loading = this.state.loading;
        var data = this.state.data;
        if(!loading) {
            var cards = this.buildSettingsCards();
        }

        return (
            <div>
                <Container className="my-4">
                    <Row>
                        <Col sm={{ span: 12 }} className="mb-3 text-center">
                            <div className="mouse" onClick={this.toggleShow}>
                                <h4 className="mb-0 mr-2 d-inline">Settings</h4>
                                {(show) ?
                                    <span className="ti-angle-up"></span>
                                :
                                    <span className="ti-angle-down"></span>
                                }
                            </div>
                        </Col>
                    </Row>
                    <Collapse in={show}>
                        <div>
                            <Row>
                                <Col sm={{ span: 12 }}>
                                    {loading ?
                                        <Loader small />
                                    :
                                        cards
                                    }
                                </Col>
                            </Row>
                        </div>
                    </Collapse>
                </Container>

            </div>
        );
    }
}

if (document.getElementById('Settings')) {
    ReactDOM.render(<Settings />, document.getElementById('Settings'));
}
