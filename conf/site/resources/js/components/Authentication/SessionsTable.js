import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Container, Row, Col, Table } from 'react-bootstrap';
import Axios from 'axios';

export default class SessionsTable extends Component {
    constructor(props) {
        super(props)

        this.state = {
            sessions: []
        }
    }

    componentDidMount() {
        this.getSessions();
    }

    getSessions = () => {
        var url = 'api/auth/sessions?token=' + window.token;

        Axios.get(url)
        .then((resp) => {
            this.setState({
                sessions: resp.data.response
            })
        })
    }

    render() {
        var sessions = this.state.sessions;

        return (
            <Container className="mb-4">
                <Row>
                    <Col sm={{ span: 12 }} className="mb-3 text-center">
                        <h5>Login Sessions</h5>
                        <Table responsive>
                            <thead>
                                <tr>
                                    <th>IP</th>
                                    <th>Expires</th>
                                    <th>Created at</th>
                                </tr>
                            </thead>
                            <tbody>
                                {sessions.map((e,i) => {
                                    return(
                                        <tr key={i}>
                                            <td>{e.ip}</td>
                                            <td>{new Date(e.expires * 1000).toLocaleDateString() + ' ' + new Date(e.expires * 1000).toLocaleTimeString()}</td>
                                            <td>{e.created_at}</td>
                                        </tr>
                                    )
                                })}
                            </tbody>
                        </Table>
                    </Col>
                </Row>
            </Container>
        );
    }
}

if (document.getElementById('SessionsTable')) {
    ReactDOM.render(<SessionsTable />, document.getElementById('SessionsTable'));
}
