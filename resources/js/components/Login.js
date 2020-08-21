import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Container, Row, Form, Toast, Modal } from 'react-bootstrap';
import { Col } from 'react-bootstrap';
import { Button } from 'react-bootstrap';
import Axios from 'axios';
import { toast } from 'react-toastify';
import Cookies from 'js-cookie';

export default class Login extends Component {
    constructor(props) {
        super(props)

        this.state = {
            loginEmailInput: '',
            loginPasswordInput: ''
        }
    }

    updateTextField = (e) => {
        this.setState({
            [e.target.id]: e.target.value
        })
    }

    login = (e) => {
        e.preventDefault();

        var data = {
            email: this.state.loginEmailInput,
            password: this.state.loginPasswordInput
        }

        var url = 'api/auth/login';
        Axios.post(url, data)
        .then((resp) => {
            var token = resp.data.access_token;
            var expires = (resp.data.expires_in / 60) / 24;
            Cookies.set('auth', token, { expires: expires })
            window.location.reload(true);
        })
    }

    toggleShow = () => {
        if(this.state.show) {
            this.setState({
                show: false
            })
        } else {
            this.setState({
                show: true
            })
        }
    }

    render() {
        var show = this.state.show;

        return (
            <Container>
                <Row>
                    <Col
                        xs={{ span: 12 }}
                        className="pb-2 text-center"
                    >
                        <Button variant="primary" onClick={this.toggleShow}>Login</Button>
                        <Modal show={show} onHide={this.toggleShow}>
                            <Modal.Header closeButton>
                                <Modal.Title>Login</Modal.Title>
                            </Modal.Header>
                            <Modal.Body>
                                <Form onSubmit={this.login}>
                                    <Form.Group controlId="loginEmailInput">
                                        <Form.Label>Email address</Form.Label>
                                        <Form.Control type="email" placeholder="admin@admin.com" onInput={this.updateTextField} required />
                                    </Form.Group>
                                    <Form.Group controlId="loginPasswordInput">
                                        <Form.Label>Password</Form.Label>
                                        <Form.Control type="password" onInput={this.updateTextField} required />
                                    </Form.Group>
                                    <Button variant="primary" type="submit">Login</Button>
                                </Form>
                            </Modal.Body>
                        </Modal>
                    </Col>
                </Row>
            </Container>
        );
    }
}

if (document.getElementById('login')) {
    ReactDOM.render(<Login />, document.getElementById('login'));
}
