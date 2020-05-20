import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import Spinner from 'react-bootstrap/Spinner';
import { Container } from 'react-bootstrap';
import { Row } from 'react-bootstrap';
import { Col } from 'react-bootstrap';

export default class Loader extends Component {
    constructor(props) {
        super(props);

        this.state = {

        };
    }

    render() {
        if(this.props.small) {
            return (
                <Container fluid>
                    <Row className="text-center align-items-center">
                        <Col
                            lg={{ span: 2, offset: 5}}
                            md={{ span: 4, offset: 4}}
                            sm={{ span: 4, offset: 4}}
                            xs={{ span: 12}}
                        >
                            <Spinner animation="grow" size="lg"/>
                        </Col>
                    </Row>
                </Container>
            )
        }

        return (
            <Container fluid>
                <Row className="fullscreen text-center align-items-center">
                    <Col
                        lg={{ span: 2, offset: 5}}
                        md={{ span: 4, offset: 4}}
                        sm={{ span: 4, offset: 4}}
                        xs={{ span: 12}}
                    >
                        <Spinner animation="grow" size="lg"/>
                    </Col>
                </Row>
            </Container>
        );
    }
}

if (document.getElementById('loader')) {
    ReactDOM.render(<Loader />, document.getElementById('loader'));
}
