import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Row, Col} from 'react-bootstrap';
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
