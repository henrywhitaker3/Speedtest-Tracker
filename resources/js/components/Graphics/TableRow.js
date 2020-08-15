import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Modal, Button } from 'react-bootstrap';
import Axios from 'axios';
import { toast } from 'react-toastify';

export default class TableRow extends Component {
    constructor(props) {
        super(props)

        this.state = {
            data: this.props.data,
            show: false,
        }
    }

    toggleShow = () => {
        var show = this.state.show;
        if(show) {
            this.setState({
                show: false
            });
        } else {
            this.setState({
                show: true
            });
        }
    }

    delete = (id) => {
        var url = 'api/speedtest/delete/' + id;

        Axios.delete(url)
        .then((resp) => {
            console.log(resp);
            toast.success('Speedtest deleted');
        })
        .catch((err) => {
            if(err.response.status == 404) {
                toast.warning('Speedtest not found');
            } else {
                toast.error('Something went wrong');
            }
        })

        this.toggleShow();
    }

    render() {
        var e = this.state.data;
        var show = this.state.show;

        if(e.failed != true) {
            return (
                <tr>
                    <td>{e.id}</td>
                    <td>{new Date(e.created_at).toLocaleString()}</td>
                    <td>{e.download}</td>
                    <td>{e.upload}</td>
                    <td>{e.ping}</td>
                    {e.server_host != null ?
                        <td>
                            <span onClick={this.toggleShow} className="ti-arrow-top-right mouse"></span>
                            <Modal show={show} onHide={this.toggleShow}>
                                <Modal.Header closeButton>
                                    <Modal.Title>More info</Modal.Title>
                                </Modal.Header>
                                <Modal.Body className="text-center">
                                    <p>Server ID: {e.server_id}</p>
                                    <p>Name: {e.server_name}</p>
                                    <p>Host: {e.server_host}</p>
                                    <p>URL: <a href={e.url} target="_blank" rel="noopener noreferer">Speedtest.net</a></p>
                                    {e.scheduled != undefined &&
                                        <p>Type: {e.scheduled == true ? 'scheduled' : 'manual'}</p>
                                    }
                                    <Button variant="danger" onClick={() => { this.delete(e.id) }}>Delete</Button>
                                </Modal.Body>
                            </Modal>
                        </td>
                    :
                        <td></td>
                    }
                </tr>
            );
        } else {
            return (
                <tr>
                    <td>{e.id}</td>
                    <td>{new Date(e.created_at).toLocaleString()}</td>
                    <td><span className="ti-close text-danger"></span></td>
                    <td><span className="ti-close text-danger"></span></td>
                    <td><span className="ti-close text-danger"></span></td>
                    <td></td>
                </tr>
            );
        }
    }
}

if (document.getElementById('tableRow')) {
    ReactDOM.render(<TableRow />, document.getElementById('tableRow'));
}
