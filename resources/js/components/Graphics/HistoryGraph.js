import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import Axios from 'axios';
import { Spinner, Container, Row, Form, Card } from 'react-bootstrap';
import { Line, Bar } from 'react-chartjs-2';
import { Col } from 'react-bootstrap';
import { toast } from 'react-toastify';

export default class HistoryGraph extends Component {
    constructor(props) {
        super(props)

        this.state = {
            days: props.days,
            time: props.dlUl,
            fail: props.fail,
            config: props.config,
            duData: {},
            duOptions: {},
            pingData: {},
            pingOptions: {},
            failData: {},
            failOptions: {},
            loading: true,
            interval: null,
            graph_ul_dl_enabled: true,
            graph_ul_dl_width: 6,
            graph_failure_enabled: true,
            graph_failure_width: 6,
            graph_ping_enabled: true,
            graph_ping_width: 6,
            firstUpdate: false,
        }
    }

    componentDidMount = () => {
    }

    componentDidUpdate() {
        if(
            this.state.time != this.props.dlUl ||
            this.state.fail != this.props.fail ||
            this.state.config != this.props.config ||
            this.state.days != this.props.days
        ) {
            this.setState({
                time: this.props.dlUl,
                fail: this.props.fail,
                config: this.props.config,
                days: this.props.days
            });

            if(this.state.config !== null) {
                this.processData();
            }
        }

        if(
            !this.state.firstUpdate &&
            this.state.config !== null
        ) {
            this.processData();
            this.setState({
                firstUpdate: true,
            });
        }
    }

    processData() {
        this.processConfig();
        this.processDlUlPing();
        this.processFailure();
    }

    processConfig() {
        this.setState({
            graph_ul_dl_enabled: Boolean(Number(this.state.config.graphs.download_upload_graph_enabled.value)),
            graph_ul_dl_width: this.state.config.graphs.download_upload_graph_width.value,
            graph_ping_enabled: Boolean(Number(this.state.config.graphs.ping_graph_enabled.value)),
            graph_ping_width: this.state.config.graphs.ping_graph_width.value,
            graph_failure_enabled: Boolean(Number(this.state.config.graphs.failure_graph_enabled.value)),
            graph_failure_width: this.state.config.graphs.failure_graph_width.value,
        });
    }

    processDlUlPing() {
        let days = this.state.days;

        var duData = {
            labels: [],
            datasets:[
                {
                    data: [],
                    label: 'Download',
                    borderColor: "#fca503",
                    fill: false,
                    yAxisID: 'y-axis-1',
                },
                {
                    data: [],
                    label: 'Upload',
                    borderColor: "#3e95cd",
                    fill: false,
                    yAxisID: 'y-axis-2',
                }
            ],
        };
        var duOptions = {
            maintainAspectRatio: false,
            responsive: true,
            tooltips: {
                callbacks: {
                  label: (item) => `${item.yLabel} Mbit/s`,
                },
            },
            title: {
                display: false,
                text: 'Speedtests results for the last ' + days + ' days',
            },
            scales: {
                xAxes: [{
                    display: false,
                    scaleLabel: {
                        display: true,
                        labelString: 'DateTime'
                    }
                }],
                yAxes: [
                    {
                      type: 'linear',
                      display: true,
                      position: 'left',
                      id: 'y-axis-1',
                    },
                    {
                      type: 'linear',
                      display: true,
                      position: 'right',
                      id: 'y-axis-2',
                      gridLines: {
                        drawOnArea: false,
                      },
                    },
                  ],
            },
            elements: {
                point:{
                    radius: 0,
                    hitRadius: 8
                }
            }
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
            maintainAspectRatio: false,
            responsive: true,
            tooltips: {
                callbacks: {
                  label: (item) => `${item.yLabel} ms`,
                },
            },
            title: {
                display: false,
                text: 'Ping results for the last ' + days + ' days',
            },
            scales: {
                xAxes: [{
                    display: false,
                    scaleLabel: {
                        display: true,
                        labelString: 'DateTime'
                    }
                }],
            },
            elements: {
                point:{
                    radius: 0,
                    hitRadius: 8
                }
            }
        }

        this.state.time.forEach(e => {
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
            duData.labels.push(new Date(e.created_at).toLocaleString());
            pingData.labels.push(new Date(e.created_at).toLocaleString());
        });

        this.setState({
            duData: duData,
            duOptions: duOptions,
            pingData: pingData,
            pingOptions: pingOptions,
            loading: false,
        });
    }

    processFailure() {
        let days = this.state.days;

        var failData = {
            labels: [],
            datasets: [
                {
                    data: [],
                    label: 'Successful',
                    backgroundColor: '#07db71'
                },
                {
                    data: [],
                    label: 'Failed',
                    backgroundColor: '#E74C3C'
                },
            ],
        };
        var failOptions = {
            maintainAspectRatio: false,
            responsive: true,
            tooltips: {
                callbacks: {
                  label: (item) => `${item.yLabel} speedtests`,
                },
            },
            scales: {
                xAxes: [{
                    stacked: true
                }],
                yAxes: [{
                    stacked: true
                }]
            }
        };

        this.state.fail.forEach(e => {
            var success = {x: e.date, y: e.success};
            var fail = {x: e.date, y: e.failure};
            failData.datasets[0].data.push(success);
            failData.datasets[1].data.push(fail);
            failData.labels.push(new Date(e.date).toLocaleString([], {year: '2-digit', month:'2-digit', day:'2-digit'}));
        })

        this.setState({
            failData: failData,
            failOptions: failOptions
        });
    }

    updateDays = (e) => {
        var days = e.target.value;
        if(days) {
            toast.info('Showing results for the last ' + days + ' days');
            this.props.updateDays(days);
        }
    }

    render() {
        var loading = this.state.loading;
        var duData = this.state.duData;
        var duOptions = this.state.duOptions;
        var pingData = this.state.pingData;
        var pingOptions = this.state.pingOptions;
        var failData = this.state.failData;
        var failOptions = {
            maintainAspectRatio: false,
            responsive: true,
            tooltips: {
                callbacks: {
                  label: (item) => `${item.yLabel} speedtests`,
                },
            },
            scales: {
                xAxes: [{
                    stacked: true,
                    gridLines: {
                        display: false
                    }
                }],
                yAxes: [{
                    stacked: true,
                    ticks: {
                        stepSize: 1
                    }
                }]
            }
        };
        var days = this.state.days;

        var graph_ul_dl_enabled = this.state.graph_ul_dl_enabled;
        var graph_ul_dl_width = this.state.graph_ul_dl_width;
        var graph_ping_enabled = this.state.graph_ping_enabled;
        var graph_ping_width = this.state.graph_ping_width;
        var graph_failure_enabled = this.state.graph_failure_enabled;
        var graph_failure_width = this.state.graph_failure_width;

        var dlClasses = 'my-2 home-graph ';
        var pingClasses = 'my-2 home-graph ';
        var failureClasses = 'my-2 home-graph ';

        if(graph_ul_dl_enabled == true) {
            //
        } else {
            dlClasses += 'd-none ';
        }

        if(graph_ping_enabled == true) {
            //
        } else {
            pingClasses += 'd-none ';
        }

        if(graph_failure_enabled == true) {
            //
        } else {
            failureClasses += 'd-none ';
        }

        if(loading) {
            return (
                <div>
                    <Spinner animation="grow" />
                </div>
            )
        } else {
            return (
                <Container className="mb-4 mt-1" fluid>
                    <Row>
                        <Col
                            lg={{ span: graph_ul_dl_width }}
                            md={{ span: graph_ul_dl_width }}
                            sm={{ span: 12 }}
                            xs={{ span: 12 }}
                            className={dlClasses}
                        >
                            <Card className="shadow-sm">
                                <Card.Body>
                                    <Line data={duData} options={duOptions} height={440} />
                                </Card.Body>
                            </Card>
                        </Col>
                        <Col
                            lg={{ span: graph_ping_width }}
                            md={{ span: graph_ping_width }}
                            sm={{ span: 12 }}
                            xs={{ span: 12 }}
                            className={pingClasses}
                        >
                            <Card className="shadow-sm">
                                <Card.Body>
                                    <Line data={pingData} options={pingOptions} height={440} />
                                </Card.Body>
                            </Card>
                        </Col>
                        <Col
                            lg={{ span: graph_failure_width }}
                            md={{ span: graph_failure_width }}
                            sm={{ span: 12 }}
                            xs={{ span: 12 }}
                            className={failureClasses}
                        >
                            <Card className="shadow-sm h-100">
                                <Card.Body className="w-100 h-100">
                                    <Bar data={failData} options={failOptions} height={440} />
                            </Card.Body>
                            </Card>
                        </Col>
                    </Row>
                    <Row>
                        <Col sm={{ span: 12 }}>
                            <div className="text-center">
                                <div className="d-inline-flex align-items-center mb-2">
                                    <h4 className="d-inline mb-0">Show results for the last</h4>
                                    <Form.Control id="duDaysInput" className="d-inline-block mx-2" defaultValue={days} onInput={this.updateDays}></Form.Control>
                                    <h4 className="d-inline mb-0">days</h4>
                                </div>
                            </div>
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
