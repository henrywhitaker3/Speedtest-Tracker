import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter, Route, Link } from "react-router-dom";
import { Container } from 'react-bootstrap';
import { Row } from 'react-bootstrap';
import { Col } from 'react-bootstrap';
import { Button } from 'react-bootstrap';

export default class ErrorPage extends Component {
    constructor(props) {
        super(props);

        var colour = '';
        var message = false;
        switch(this.props.code.toString()[0]) {
            case 2:
            case '2':
                colour = 'success';
                break;
            case 4:
            case '4':
                colour = 'danger';
                break;
            case 5:
            case '5':
            default:
                colour = 'warning';
                break;
        }

        switch(this.props.code) {
            case '400':
                message = 'Bad request'
                break;
            case '401':
                message = 'You aren\'t authenticated';
                break;
            case '403':
                message = 'You aren\'t authorised to view this page';
                break;
            case '404':
                message = 'Page not found';
                break;
            case '405':
                message = 'Method not allowed'
                break;
            case '413':
                message = 'Request too large'
                break;
            case '422':
                message = 'Your request was unprocessable'
                break;
        }

        this.state = {
            code: this.props.code,
            colour: colour,
            message: message
        };
    }

    render() {
        const colour = this.state.colour;
        const code = this.state.code;
        const message = this.state.message;
        return (
            <Container fluid>
                <Row className="fullscreen text-center align-items-center">
                    <Col
                        lg={{ span: 2, offset: 5}}
                        md={{ span: 4, offset: 4}}
                        sm={{ span: 4, offset: 4}}
                        xs={{ span: 12}}
                    >
                        <h1 className={'text-' + colour + ' mb-0'}>{code}</h1>
                        {message &&
                            <p className={colour + '-text mt-0 mb-2'}>{message}</p>
                        }
                        <Link to="/" className={'waves-effect waves-' + colour + ' btn ' + colour}><Button variant={colour}>Go home</Button></Link>
                    </Col>
                </Row>
            </Container>
        );
    }
}

if (document.getElementById('errorpage')) {
    ReactDOM.render(<ErrorPage />, document.getElementById('errorpage'));
}
