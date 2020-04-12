import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Container, Row, Form, Toast } from 'react-bootstrap';
import { Col } from 'react-bootstrap';
import { Button } from 'react-bootstrap';
import Axios from 'axios';
import { toast } from 'react-toastify';

export default class Login extends Component {
    constructor(props) {
        super(props)

        this.state = {
            loginEmailInput: '',
            loginPasswordInput: '',
            toast: null,
        }
    }

    login = (e) => {
        e.preventDefault();

        var data = {
            email: this.state.loginEmailInput,
            password: this.state.loginPasswordInput
        };
        var url = '/api/auth/login';

        Axios.post(url, data)
        .then((resp) => {
            var token = resp.data;
            this.props.setToken(token);
        })
        .catch((err) => {
            toast.error('Something went wrong. Please try again.')
        });
    }

    updateTextField = (e) => {
        this.setState({
            [e.target.id]: e.target.value
        });
    }

    render() {
        var error = this.state.error;

        return (
            <Container>
                <Row className="fullscreen align-items-center">
                    <Col
                        lg={{ span: 4, offset: 4 }}
                        md={{ span: 6, offset: 3 }}
                        sm={{ span: 10, offset: 1 }}
                        xs={{ span: 12 }}
                    >
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
                    </Col>
                </Row>
            </Container>
        );
    }
}

if (document.getElementById('login')) {
    ReactDOM.render(<Login />, document.getElementById('login'));
}
