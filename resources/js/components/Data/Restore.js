import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Button, Modal, Form, Tooltip, OverlayTrigger } from 'react-bootstrap';
import { toast } from 'react-toastify';
import Axios from 'axios';

export default class Restore extends Component {
    constructor(props) {
        super(props);

        this.state = {
            show: false,
            data: null,
            uploadReady: false,
            filename: 'Upload your backup JSON'
        };
    }

    showModal = () => {
        this.setState({
            show: true
        });
    }

    hideModal = () => {
        this.setState({
            show: false
        });
    }

    readFile = (e) => {
        var file = e.target.files[0];
        var reader = new FileReader();
        reader.readAsText(file, 'UTF-8');
        reader.onload = function(evt) {
            try {
                var data = evt.target.result.trim();
                var data = JSON.parse(data);
                this.setState({
                    data: data,
                    uploadReady: true,
                    filename: file.name
                });
            } catch(e) {
                console.log(e);
                toast.error('Your upload file is not valid JSON');
            }
        }.bind(this)
        reader.onerror = function (evt) {
            toast.error('Something went wrong parsing your backup file.');
        }
    }

    uploadFile = () => {
        var data = { data: this.state.data };
        var url = '/api/restore';

        Axios.post(url, data)
        .then((resp) => {
            toast.success('Your is being restored...');
            this.setState({
                show: false,
                data: null,
                uploadReady: false,
                filename: 'Upload your backup JSON'
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
                <Button variant="secondary" className="mx-2" onClick={this.showModal}>Restore</Button>

                <Modal show={show} onHide={this.hideModal} animation={true}>
                    <Modal.Header closeButton>
                        <Modal.Title>Restore from a backup</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        <p>Upload your JSON backup file here:</p>
                        <Form.File
                            id="restoreFileInput"
                            label="Upload JSON file"
                            className="mb-3"
                            custom
                        >
                            <Form.File.Input onChange={this.readFile} />
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
