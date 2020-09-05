import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Container } from 'react-bootstrap';
import { Row } from 'react-bootstrap';
import { Col } from 'react-bootstrap';
import Backup from './Backup';
import Restore from './Restore';

export default class DataRow extends Component {
    render() {
        if( (window.config.auth == true && window.authenticated == true) || window.config.auth == false) {
            return (
                <Container className="mb-4">
                    <Row>
                        <Col sm={{ span: 12 }} className="text-center">
                            <p>Use these buttons to backup/restore your data</p>
                        </Col>
                    </Row>
                    <Row>
                        <Col sm={{ span: 12 }} className="text-center">
                            <Backup />
                            <Restore />
                        </Col>
                    </Row>
                </Container>
            );
        } else {
            return (
                <></>
            )
        }
    }
}

if (document.getElementById('DataRow')) {
    ReactDOM.render(<DataRow />, document.getElementById('DataRow'));
}
