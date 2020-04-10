import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Container, Row } from 'react-bootstrap';
import { Col } from 'react-bootstrap';
import Version from './Version';

export default class Footer extends Component {
    render() {
        return (
            <Container>
                <Row>
                    <Col sm={{ span: 12 }} className="text-center">
                        <Version />
                        <p className="text-muted">See the code on <a href="https://github.com/henrywhitaker3/Speedtest-Tracker" target="_blank" rel="noopener noreferrer">GitHub</a></p>
                    </Col>
                </Row>
            </Container>
        );
    }
}

if (document.getElementById('Footer')) {
    ReactDOM.render(<Footer />, document.getElementById('Footer'));
}
