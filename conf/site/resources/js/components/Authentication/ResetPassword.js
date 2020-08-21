import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Container, Row, Col, Collapse, Button, Modal, Form } from 'react-bootstrap';
import SessionsTable from './SessionsTable';
import Axios from 'axios';
import { toast } from 'react-toastify';

export default class ResetPassword extends Component {
    constructor(props) {
        super(props)

        this.state = {
            showModal: false,
            currentPassword: '',
            newPassword: '',
            newPasswordConfirmation: '',
            logoutDevices: false
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

    updateTextField = (e) => {
        this.setState({
            [e.target.id]: e.target.value
        })
    }

    updateCheckbox = (e) => {
        this.setState({
            [e.target.id]: e.target.checked
        });
    }

    changePassword = (e) => {
        e.preventDefault();

        var data = {
            currentPassword: this.state.currentPassword,
            newPassword: this.state.newPassword,
            newPassword_confirmation: this.state.newPasswordConfirmation,
            logoutDevices: this.state.logoutDevices
        }

        var url = 'api/auth/change-password?token=' + window.token;
        Axios.post(url, data)
        .then((resp) => {
            toast.success('Password updated');
            this.toggleModal();
            if(this.state.logoutDevices == true) {
                location.reload(true);
            }
        })
        .catch((err) => {
            if(err.response) {
                for(var key in err.response.data.error) {
                    toast.error(err.response.data.error[key][0]);
                }
            }
        })
    }

    render() {
        var showModal = this.state.showModal;

        return (
            <div>
                <Button variant="primary" onClick={this.toggleModal} className="mb-3">Change password</Button>
                <Modal show={showModal} onHide={this.toggleModal}>
                    <Modal.Header closeButton>
                        <Modal.Title>Change password</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        <Form onSubmit={this.changePassword}>
                            <Form.Group controlId="currentPassword">
                                <Form.Label>Current password</Form.Label>
                                <Form.Control type="password" onInput={this.updateTextField} required />
                            </Form.Group>
                            <Form.Group controlId="newPassword">
                                <Form.Label>New Password</Form.Label>
                                <Form.Control type="password" onInput={this.updateTextField} required />
                            </Form.Group>
                            <Form.Group controlId="newPasswordConfirmation">
                                <Form.Label>Confirm New Password</Form.Label>
                                <Form.Control type="password" onInput={this.updateTextField} required />
                            </Form.Group>
                            <Button variant="primary" type="submit" className="d-inline-block">Change password</Button>
                            <Form.Group controlId="logoutDevices" className="d-inline-block ml-2">
                                <Form.Check type="checkbox" label="Log everywhere out" onInput={this.updateCheckbox} />
                            </Form.Group>
                        </Form>
                    </Modal.Body>
                </Modal>
            </div>
        );
    }
}

if (document.getElementById('ResetPassword')) {
    ReactDOM.render(<ResetPassword />, document.getElementById('ResetPassword'));
}
