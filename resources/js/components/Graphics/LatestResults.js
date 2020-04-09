import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import Axios from 'axios';
import Widget from './Widget';
import { Container, Row, Spinner } from 'react-bootstrap';
import { Col } from 'react-bootstrap';
import { Button } from 'react-bootstrap';
import { toast } from 'react-toastify';

export default class LatestResults extends Component {
    constructor(props) {
        super(props)

        this.state = {
            data: {},
            interval: null,
            loading: true,
        }
    }

    componentDidMount = () => {
        this.getData();
        var int = setInterval(this.getData, 10000);
        this.setState({
            interval: int,
        });
    }

    getData = () => {
        var url = '/api/speedtest/latest';

        Axios.get(url)
        .then((resp) => {
            this.setState({
                data: resp.data,
                loading: false
            });
        })
        .catch((err) => {
            this.setState({
                data: false
            });
            console.log(err);
        })
    }

    newScan = () => {
        var url = '/api/speedtest/run';

        Axios.get(url)
        .then((resp) => {
            toast.info('A scan has been queued. This page will refresh when the scan has finished.');
        })
        .catch((err) => {
            if(err.response) {
                if(err.response.status == 429) {
                    toast.error('You are doing that too much. Try again later.');
                }
                console.log(err.response);
            } else {
                console.log(err.data);
            }
        })
    }

    render() {
        var loading = this.state.loading;
        var data = this.state.data;

        if(loading && data !== false) {
            return (
                <Container fluid>
                    <Row>
                        <Col sm={{ span: 12 }}>
                            <Spinner animation="grow" />
                        </Col>
                    </Row>
                </Container>
            );
        } else if(data === false) {
            return (
                <Container fluid>
                    <Row>
                        <Col sm={{ span: 12 }} className="text-center">
                            <div>
                                <Button variant="primary" onClick={this.newScan}>Start your first scan!</Button>
                            </div>
                        </Col>
                    </Row>
                </Container>
            );
        } else {
            return (
                <Container fluid>
                    <Row>
                        <Col sm={{ span: 12 }} className="text-center">
                            <h4>Latest test results:</h4>
                        </Col>
                    </Row>
                    <Row>
                        <Col
                            lg={{ span: 2, offset: 3 }}
                            md={{ span: 4 }}
                            sm={{ span: 4 }}
                            className="my-2"
                        >
                            <Widget
                                title="Ping"
                                value={parseFloat(data.data.ping).toFixed(1)}
                                unit="ms"
                                icon="ping"
                            />
                        </Col>
                        <Col
                            lg={{ span: 2 }}
                            md={{ span: 4 }}
                            sm={{ span: 4 }}
                            className="my-2"
                        >
                            <Widget
                                title="Download"
                                value={parseFloat(data.data.download).toFixed(1)}
                                unit="Mbit/s"
                                icon="dl"
                            />
                        </Col>
                        <Col
                            lg={{ span: 2 }}
                            md={{ span: 4 }}
                            sm={{ span: 4 }}
                            className="my-2"
                        >
                            <Widget
                                title="Upload"
                                value={parseFloat(data.data.upload).toFixed(1)}
                                unit="Mbit/s"
                                icon="ul"
                            />
                        </Col>
                    </Row>
                    <Row>
                        <Col sm={{ span: 12 }} className="text-center mb-2">
                            <p className="text-muted mb-0">Last scan performed at: {new Date(data.data.created_at).toLocaleString()}</p>
                        </Col>
                    </Row>
                    <Row>
                        <Col sm={{ span: 12 }} className="text-center">
                            <div>
                                <Button variant="primary" onClick={this.newScan}>Scan again</Button>
                            </div>
                        </Col>
                    </Row>
                </Container>
            );
        }
    }
}

if (document.getElementById('LatestResults')) {
    ReactDOM.render(<LatestResults />, document.getElementById('LatestResults'));
}
