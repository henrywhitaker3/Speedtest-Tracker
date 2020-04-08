import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import Axios from 'axios';
import { Spinner } from 'react-bootstrap';
import { Line } from 'react-chartjs-2';

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
        }
    }

    componentDidMount = () => {
        this.getData()
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
                    display: true,
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
                    display: true,
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

    render() {
        var loading = this.state.loading;
        var duData = this.state.duData;
        var duOptions = this.state.duOptions;
        var pingData = this.state.pingData;
        var pingOptions = this.state.pingOptions;

        if(loading) {
            return (
                <div>
                    <Spinner animation="grow" />
                </div>
            )
        } else {
            return (
                <div>
                    <Line data={duData} options={duOptions} />
                    <Line data={pingData} options={pingOptions} />
                </div>
            );
        }
    }
}

if (document.getElementById('HistoryGraph')) {
    ReactDOM.render(<HistoryGraph />, document.getElementById('HistoryGraph'));
}
