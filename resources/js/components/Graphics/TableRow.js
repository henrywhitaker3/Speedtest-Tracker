import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Modal } from 'react-bootstrap';

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

    render() {
        var e = this.state.data;
        var show = this.state.show;

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
                            <Modal.Header>
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
                            </Modal.Body>
                        </Modal>
                    </td>
                :
                    <td></td>
                }
            </tr>
        );
    }
}

if (document.getElementById('tableRow')) {
    ReactDOM.render(<TableRow />, document.getElementById('tableRow'));
}
