import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import Axios from 'axios';
import { Spinner, Container, Row, Form, Card } from 'react-bootstrap';
import { Line } from 'react-chartjs-2';
import { Col } from 'react-bootstrap';
import { toast } from 'react-toastify';

export default class HistoryGraph extends Component {
    constructor(props) {
        super(props)

        this.state = {
            days: 30,
            token: this.props.token,
            duData: {},
            duOptions: {},
            pingData: {},
            pingOptions: {},
            loading: true,
            interval: null,
        }
    }

    componentDidMount = () => {
        this.getData();
        var int = setInterval(this.getData, 10000);
        this.setState({
            interval: int,
        });
    }

    getData = (days = this.state.days) => {
        var url = '/api/speedtest/time/' + days + '?token=' + this.state.token.access_token;

        Axios.get(url)
        .then((resp) => {
            var duData = {
                labels: [],
                datasets:[
                    {
                        data: [],
                        label: 'Download',
                        borderColor: "#fca503",
                        fill: false,
                    },
                    {
                        data: [],
                        label: 'Upload',
                        borderColor: "#3e95cd",
                        fill: false,
                    }
                ],
            };
            var duOptions = {
                tooltips: {
                    callbacks: {
                      label: (item) => `${item.yLabel} Mbit/s`,
                    },
                },
                title: {
                    display: false,
                    text: 'Speedtests results for the last ' + days + ' days',
                },
            };

            var pingData = {
                labels: [],
                datasets:[
                    {
                        data: [],
                        label: 'Ping',
                        borderColor: "#07db71",
                        fill: false,
                    },
                ],
            };
            var pingOptions = {
                tooltips: {
                    callbacks: {
                      label: (item) => `${item.yLabel} ms`,
                    },
                },
                title: {
                    display: false,
                    text: 'Ping results for the last ' + days + ' days',
                },
            }

            resp.data.data.forEach(e => {
                var download = {
                    t: new Date(e.created_at),
                    y: e.download,
                };
                var upload = {
                    t: new Date(e.created_at),
                    y: e.upload,
                };
                var ping = {
                    t: new Date(e.created_at),
                    y: parseFloat(e.ping)
                }
                duData.datasets[0].data.push(download);
                duData.datasets[1].data.push(upload);
                pingData.datasets[0].data.push(ping);
                duData.labels.push(new Date(e.created_at).toLocaleDateString());
                pingData.labels.push(new Date(e.created_at).toLocaleDateString());
            });

            this.setState({
                duData: duData,
                duOptions: duOptions,
                pingData: pingData,
                pingOptions: pingOptions,
                loading: false,
            });
        })
        .catch((err) => {
            console.log(err);
        })
    }

    updateDays = (e) => {
        var days = e.target.value;
        if(days) {
            this.getData(days);
            clearInterval(this.state.int);
            var int = setInterval(this.getData, 10000);
            toast.info('Showing results for the last ' + days + ' days');
            this.setState({
                days: days,
                interval: int
            });
        }
    }

    render() {
        var loading = this.state.loading;
        var duData = this.state.duData;
        var duOptions = this.state.duOptions;
        var pingData = this.state.pingData;
        var pingOptions = this.state.pingOptions;
        var days = this.state.days;

        if(loading) {
            return (
                <div>
                    <Spinner animation="grow" />
                </div>
            )
        } else {
            return (
                <Container className="my-4" fluid>
                    <Row>
                        <Col sm={{ span: 12 }}>
                            <div className="text-center">
                                <h4 className="d-inline">Show results for the last</h4>
                                <Form.Control id="duDaysInput" className="d-inline mx-2" defaultValue={days} onInput={this.updateDays}></Form.Control>
                                <h4 className="d-inline">days</h4>
                                <p className="text-muted">This data refreshes every 10 seconds</p>
                            </div>
                        </Col>
                    </Row>
                    <Row>
                        <Col
                            lg={{ span: 6 }}
                            md={{ span: 12 }}
                            sm={{ span: 12 }}
                            xs={{ span: 12 }}
                            className="my-2"
                        >
                            <Card className="shadow-sm">
                                <Card.Body>
                                    <Line data={duData} options={duOptions} />
                                </Card.Body>
                            </Card>
                        </Col>
                        <Col
                            lg={{ span: 6 }}
                            md={{ span: 12 }}
                            sm={{ span: 12 }}
                            xs={{ span: 12 }}
                            className="my-2"
                        >
                            <Card className="shadow-sm">
                                <Card.Body>
                                    <Line data={pingData} options={pingOptions} />
                                </Card.Body>
                            </Card>
                        </Col>
                    </Row>
                </Container>
            );
        }
    }
}

if (document.getElementById('HistoryGraph')) {
    ReactDOM.render(<HistoryGraph />, document.getElementById('HistoryGraph'));
}
