import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Container, Row, Col, Collapse, Button, Modal } from 'react-bootstrap';
import SessionsTable from './SessionsTable';
import ResetPassword from './ResetPassword';

export default class Authentication extends Component {
    constructor(props) {
        super(props)

        this.state = {
            showCollapse: false,
            showModal: false
        }
    }

    toggleCollapse = () => {
        if(this.state.showCollapse) {
            this.setState({
                showCollapse: false
            });
        } else {
            this.setState({
                showCollapse: true
            });
        }
    }

    toggleModal = () => {
        if(this.state.showModal) {
            this.setState({
                showModal: false
            });
        } else {
            this.setState({
                showModal: true
            });
        }
    }

    render() {
        var showCollapse = this.state.showCollapse;
        var showModal = this.state.showModal;

        if( (window.config.auth == true && window.authenticated == true)) {
            return (
                <Container className="mb-4">
                    <Row>
                        <Col sm={{ span: 12 }} className="mb-3 text-center">
                            <div className="mouse"  aria-controls="testsTable" onClick={this.toggleCollapse} aria-expanded={showCollapse}>
                                <h4 className="d-inline mr-2">Authentication</h4>
                                {(showCollapse) ?
                                    <span className="ti-angle-up"></span>
                                :
                                    <span className="ti-angle-down"></span>
                                }
                            </div>
                        </Col>
                    </Row>
                    <Collapse in={showCollapse}>
                        <div>
                            <Row>
                                <Col sm={{ span: 12 }} className="text-center">
                                    <ResetPassword />
                                </Col>
                            </Row>
                            <Row>
                                <Col sm={{ span: 12 }} className="text-center">
                                    <SessionsTable />
                                </Col>
                            </Row>
                        </div>
                    </Collapse>
                </Container>
            );
        } else {
            return (
                <></>
            );
        }
    }
}

if (document.getElementById('Authentication')) {
    ReactDOM.render(<Authentication />, document.getElementById('Authentication'));
}
