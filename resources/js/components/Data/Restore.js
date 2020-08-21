import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Button, Modal, Form, Tooltip, OverlayTrigger, Dropdown, DropdownButton } from 'react-bootstrap';
import { toast } from 'react-toastify';
import Axios from 'axios';
import CSVFileValidator from 'csv-file-validator';

export default class Restore extends Component {
    constructor(props) {
        super(props);

        this.state = {
            show: false,
            data: null,
            uploadReady: false,
            filename: 'Upload your backup',
            format: 'json'
        };
    }

    showModal = (format) => {
        this.setState({
            show: true,
            format: format
        });
    }

    hideModal = () => {
        this.setState({
            show: false
        });
    }

    readFile = (e, format) => {
        var file = e.target.files[0];
        var reader = new FileReader();
        reader.readAsText(file, 'UTF-8');
        reader.onload = function(evt) {
            var data = evt.target.result.trim();
            if(format == 'csv') {
                var csv = data.substr(45);
                var config = {
                    headers: [
                        {
                            name: "id",
                            inputName: 'id',
                            required: false,
                        },
                        {
                            name: "ping",
                            inputName: 'ping',
                            required: true,
                            requiredError: function (headerName, rowNumber, columnNumber) {
                                return `${headerName} is required in the ${rowNumber} row / ${columnNumber} column`
                            }
                        },
                        {
                            name: "upload",
                            inputName: 'upload',
                            required: true,
                            requiredError: function (headerName, rowNumber, columnNumber) {
                                return `${headerName} is required in the ${rowNumber} row / ${columnNumber} column`
                            }
                        },
                        {
                            name: "download",
                            inputName: 'download',
                            required: true,
                            requiredError: function (headerName, rowNumber, columnNumber) {
                                return `${headerName} is required in the ${rowNumber} row / ${columnNumber} column`
                            }
                        },
                        {
                            name: "created_at",
                            inputName: 'created_at',
                            required: false,
                        },
                        {
                            name: "server_id",
                            inputName: 'server_id',
                            required: false,
                        },
                        {
                            name: "server_name",
                            inputName: 'server_name',
                            required: false,
                        },
                        {
                            name: "server_host",
                            inputName: 'server_host',
                            required: false,
                        },
                        {
                            name: "url",
                            inputName: 'url',
                            required: false,
                        },
                        {
                            name: "scheduled",
                            inputName: 'scheduled',
                            required: false,
                        },
                        {
                            name: "failed",
                            inputName: 'failed',
                            required: false,
                        },
                        {
                            name: "updated_at",
                            inputName: 'updated_at',
                            required: false,
                        }
                    ]
                };
                CSVFileValidator(csv, config)
                .then((e) => {
                    if(e.inValidMessages.length > 0) {
                        toast.error('Your upload file is not valid ' + format.toUpperCase());
                    } else {
                        this.setState({
                            data: data,
                            uploadReady: true,
                            filename: file.name
                        });
                    }
                })
                .catch((e) => {
                    toast.error('Your upload file is not valid ' + format.toUpperCase());
                })
            } else {
                try {
                    var data = JSON.parse(data);
                    this.setState({
                        data: data,
                        uploadReady: true,
                        filename: file.name
                    });
                } catch(e) {
                    console.log(e);
                    toast.error('Your upload file is not valid ' + format.toUpperCase());
                }
            }
        }.bind(this)
        reader.onerror = function (evt) {
            toast.error('Something went wrong parsing your backup file.');
        }
    }

    uploadFile = () => {
        var data = { data: this.state.data, format: this.state.format };
        var url = 'api/restore?token=' + window.token;

        Axios.post(url, data)
        .then((resp) => {
            toast.success('Your data is being restored...');
            this.setState({
                show: false,
                data: null,
                uploadReady: false,
                filename: 'Upload your backup'
            });
        })
        .catch((err) => {
            console.log(err);
        })
    }

    render() {
        var show = this.state.show;
        var uploadReady = this.state.uploadReady;
        var filename = this.state.filename;

        return (
            <>
                <DropdownButton variant="secondary" title="Restore" className="m-2 d-inline-block">
                    <Dropdown.Item href="#" onClick={() => { this.showModal('json') }}>JSON</Dropdown.Item>
                    <Dropdown.Item href="#" onClick={() => { this.showModal('csv') }}>CSV</Dropdown.Item>
                </DropdownButton>

                <Modal show={show} onHide={this.hideModal} animation={true}>
                    <Modal.Header closeButton>
                        <Modal.Title>Restore from a backup</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        <p>Upload your {this.state.format.toUpperCase()} backup file here:</p>
                        <Form.File
                            id="restoreFileInput"
                            label={"Upload " + this.state.format.toUpperCase() + " file"}
                            className="mb-3"
                            custom
                        >
                            <Form.File.Input onChange={(e) => { this.readFile(e, this.state.format) }} />
                            <Form.File.Label data-browse="Choose file">
                                {filename}
                            </Form.File.Label>
                        </Form.File>
                        {uploadReady === true &&
                            <Button variant="secondary" onClick={this.uploadFile}>Restore</Button>
                        }
                    </Modal.Body>
                </Modal>
            </>
        );
    }
}

if (document.getElementById('Restore')) {
    ReactDOM.render(<Restore />, document.getElementById('Restore'));
}
