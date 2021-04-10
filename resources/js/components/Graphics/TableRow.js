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
        var url = 'api/speedtest/delete/' + id + '?token=' + window.token;

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

        this.props.refresh();
        this.toggleShow();
    }

    getDataFields = () => {
        let allFields = this.props.allFields;
        let data = this.state.data;
        let processedFields = [];

        for(var key in allFields) {
            let field = allFields[key];

            let value = data[key];

            if(field.type === 'date') {
                value = new Date(value).toLocaleString();
            } else if(field.type === 'bool') {
                value = Boolean(value) ? field.if_true : field.if_false
            }

            let final = {
                name: key,
                key: field.alias,
                value: value,
                type: field.type
            };

            processedFields.push(final);
        }

        let visible = [];
        let inModal = [];

        window.config.tables.visible_columns.forEach(column => {
            visible.push(processedFields.find(x => x.name == column));
        });

        inModal = processedFields.filter(el => {
            return !visible.includes(el);
        });

        return {
            visible: visible,
            modal: inModal
        };
    }

    render() {
        var e = this.state.data;
        var show = this.state.show;
        var fields = this.getDataFields();

        if(e.failed != true) {
            return (
                <tr>
                    {fields.visible.map((e, i) => {
                        return (
                            <td key={i}>{e.value}</td>
                        );
                    })}
                    {e.server_host != null ?
                        <td>
                            <span onClick={this.toggleShow} className="ti-arrow-top-right mouse"></span>
                            <Modal show={show} onHide={this.toggleShow}>
                                <Modal.Header closeButton>
                                    <Modal.Title>More info</Modal.Title>
                                </Modal.Header>
                                <Modal.Body className="text-center">
                                    {fields.modal.map((e, i) => {
                                        if(e.type === 'url') {
                                            return (
                                                <p key={i}>{e.key}: <a href={e.value} target="_blank" rel="noopener noreferer">Speedtest.net</a></p>
                                            );
                                        } else {
                                            return (
                                                <p key={i}>{e.key}: {e.value}</p>
                                            );
                                        }
                                    })}
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
                    {fields.visible.map((e, i) => {
                        console.log(e);
                        if(e.name === 'created_at') {
                            return <td key={i}>{new Date(e.value).toLocaleString()}</td>
                        } else if (e.name === 'id') {
                            return <td key={i}>{e.value}</td>
                        }

                        return (
                            <td key={i}><span className="ti-close text-danger"></span></td>
                        );
                    })}
                    {(window.config.auth && window.authenticated) || !window.config.auth ?
                        <td><Button variant="danger" onClick={() => { this.delete(e.id) }}>Delete</Button></td>
                    :
                        <td></td>
                    }
                </tr>
            );
        }
    }
}

if (document.getElementById('tableRow')) {
    ReactDOM.render(<TableRow />, document.getElementById('tableRow'));
}
