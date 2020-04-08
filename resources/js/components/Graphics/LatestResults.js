import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import Axios from 'axios';
import Widget from './Widget';
import { Container, Row, Spinner } from 'react-bootstrap';
import { Col } from 'react-bootstrap';

export default class LatestResults extends Component {
    constructor(props) {
        super(props)

        this.state = {
            token: this.props.token,
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
        var url = '/api/speedtest/latest?token=' + this.state.token.access_token;

        Axios.get(url)
        .then((resp) => {
            this.setState({
                data: resp.data,
                loading: false
            });
        })
        .catch((err) => {
            console.log(err);
        })
    }

    render() {
        var loading = this.state.loading;
        var data = this.state.data;

        if(loading) {
            return (
                <Container fluid>
                    <Row>
                        <Col sm={{ span: 12 }}>
                            <Spinner animation="grow" />
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
                        <Col sm={{ span: 12 }} className="text-center">
                            <p className="text-muted mb-0">Last scan performed at: {new Date(data.data.created_at).toLocaleString()}</p>
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
