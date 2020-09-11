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
        var url = 'api/speedtest/latest';

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
        var url = 'api/speedtest/run?token=' + window.token;

        Axios.get(url)
        .then((resp) => {
            toast.info('A test has been queued. This page will refresh when the test has finished.');
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
            if( (window.config.auth == true && window.authenticated == true) || window.config.auth == false) {
                return (
                    <Container fluid>
                        <Row>
                            <Col sm={{ span: 12 }} className="text-center">
                                <div>
                                    <Button variant="primary" onClick={this.newScan}>Start your first test!</Button>
                                </div>
                            </Col>
                        </Row>
                    </Container>
                );
            } else if(window.config.auth == true && window.authenticated == false) {
                return (
                    <Container fluid>
                        <Row>
                            <Col sm={{ span: 12 }} className="text-center">
                                <div>
                                    <p>Please login to run the first test</p>
                                </div>
                            </Col>
                        </Row>
                    </Container>
                );
            }
        } else {
            return (
                <Container fluid>
                    <Row>
                        <Col sm={{ span: 12 }} className="text-center mb-2">
                            <div>
                                {(window.config.auth == true && window.authenticated == true) || window.config.auth == false ?
                                    <div>
                                        <Button className="d-inline-block mx-3 mb-2" variant="primary" onClick={this.newScan}>Test again</Button>
                                        <p className="text-muted mb-0 d-inline-block">Last test performed at: {new Date(data.data.created_at).toLocaleString()}</p>
                                    </div>
                                :
                                    <div>
                                        <p className="text-muted mb-0 d-inline-block">Last test performed at: {new Date(data.data.created_at).toLocaleString()}</p>
                                    </div>
                                }
                            </div>
                        </Col>
                    </Row>
                    <Row>
                        <Col
                            lg={{ span: 4 }}
                            md={{ span: 4 }}
                            sm={{ span: 12 }}
                            className="my-2"
                        >
                            <Widget
                                title="Ping"
                                data={data}
                                failed={data.data.failed}
                                unit="ms"
                                icon="ping"
                            />
                        </Col>
                        <Col
                            lg={{ span: 4 }}
                            md={{ span: 4 }}
                            sm={{ span: 12 }}
                            className="my-2"
                        >
                            <Widget
                                title="Download"
                                data={data}
                                failed={data.data.failed}
                                unit="Mbit/s"
                                icon="dl"
                            />
                        </Col>
                        <Col
                            lg={{ span: 4 }}
                            md={{ span: 4 }}
                            sm={{ span: 12 }}
                            className="my-2"
                        >
                            <Widget
                                title="Upload"
                                data={data}
                                failed={data.data.failed}
                                unit="Mbit/s"
                                icon="ul"
                            />
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
