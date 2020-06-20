import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import Axios from 'axios';
import { Container, Row, Table, Col, Collapse, Button } from 'react-bootstrap';
import TableRow from './TableRow';

export default class TestsTable extends Component {
    constructor(props) {
        super(props)

        this.state = {
            page: 1,
            lastPage: 1,
            data: [],
            showTable: false,
        }
    }

    componentDidMount() {
        this.getData();
    }

    getData = (page = this.state.page) => {
        var url = 'api/speedtest/?page=' + page;

        Axios.get(url)
        .then((resp) => {
            var data = resp.data.data.data;
            data = this.state.data.concat(data);
            var page = resp.data.data.current_page;
            var lastPage = resp.data.data.last_page;
            this.setState({
                data: data,
                page: page,
                lastPage: lastPage,
            });
        })
        .catch((err) => {
            console.log(err);
        })
    }

    getMoreData = () => {
        var page = this.state.page;
        page = page + 1;

        this.getData(page);
    }

    toggleCollapse = () => {
        var show = this.state.showTable;

        if(show) {
            this.setState({
                showTable: false
            });
        } else {
            this.setState({
                showTable: true
            });
        }
    }

    render() {
        var page = this.state.page;
        var lastPage = this.state.lastPage;
        var data = this.state.data;
        var show = this.state.showTable;

        if(data.length > 0) {
            return (
                <Container className="mb-4 mt-4" fluid>
                    <Row>
                        <Col sm={{ span: 12 }} className="mb-3 text-center">
                            <div className="mouse"  aria-controls="testsTable" onClick={this.toggleCollapse} aria-expanded={show}>
                                <h4 className="d-inline mr-2">All tests</h4>
                                {(show) ?
                                    <span className="ti-angle-up"></span>
                                :
                                    <span className="ti-angle-down"></span>
                                }
                            </div>
                        </Col>
                    </Row>
                    <Collapse in={show}>
                        <div>
                            <Row>
                                <Col sm={{ span: 12 }} id="testsTable">
                                    <Table responsive>
                                        <thead>
                                            <tr>
                                                <th>ID</th>
                                                <th>Time</th>
                                                <th>Download (Mbit/s)</th>
                                                <th>Upload (Mbit/s)</th>
                                                <th>Ping (ms)</th>
                                                <th>More</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {data.map((e,i) => {
                                                return (
                                                    <TableRow key={e.id} data={e} />
                                                );
                                            })}
                                        </tbody>
                                    </Table>
                                </Col>
                            </Row>
                            {page < lastPage &&
                                <Row>
                                    <Col sm={{ span: 12 }} className="text-center">
                                        <Button variant="primary" onClick={this.getMoreData}>Show more</Button>
                                    </Col>
                                </Row>
                            }
                        </div>
                    </Collapse>
                </Container>
            );
        } else {
            return (
                <>
                </>
            )
        }
    }
}

if (document.getElementById('TestsTable')) {
    ReactDOM.render(<TestsTable />, document.getElementById('TestsTable'));
}
