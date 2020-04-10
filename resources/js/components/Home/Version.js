import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import Axios from 'axios';
import { toast } from 'react-toastify';
import { Modal } from 'react-bootstrap';

export default class Version extends Component {
    constructor(props) {
        super(props)

        this.state = {
            version: document.querySelector('meta[name="version"]').content,
            update: false,
            modalShow: false,
            changelog: [],
        };
    }

    componentDidMount() {
        this.checkForUpdates();
    }

    checkForUpdates = () => {
        var url = '/api/update/check';

        Axios.get(url)
        .then((resp) => {
            var update = resp.data.update;
            if(update !== false) {
                toast.info('A new version of Speedtest Tracker is available (v' + update.version + '). Go to the bottom of the page to update.');
                this.setState({
                    update: update.version,
                    changelog: update.changelog,
                });
            }
        })
        .catch((err) => {
            console.log(err);
        });
    }

    showModal = () => {
        this.setState({
            modalShow: true
        });
    }

    hideModal = () => {
        this.setState({
            modalShow: false
        });
    }

    render() {
        var version = this.state.version;
        var update = this.state.update;
        var modalShow = this.state.modalShow;
        var changelog = this.state.changelog;

        if(update === false) {
            return (
                <p className="text-muted mb-0">Speedtest Tracker Version: {version}</p>
            );
        } else {
            return (
                <div>
                    <p className="text-muted mb-0 d-inline">Speedtest Tracker Version: {version} - </p>
                    <a href="#!" className="mb-0 d-inline" onClick={this.showModal}>New version available - v{update}</a>

                    <Modal show={modalShow} onHide={this.hideModal} animation={true}>
                        <Modal.Header closeButton>
                            <Modal.Title>Update to v{update}</Modal.Title>
                        </Modal.Header>
                        <Modal.Body>
                            <h5>Changelog:</h5>
                            <ul>
                                {changelog.map((e, i) => {
                                    if(e.link == '') {
                                        return (
                                            <li key={i}>{e.description}</li>
                                        );
                                    } else {
                                        <li key={i}><a href={e.link} target="_blank" rel="noopener noreferer">{e.description}</a></li>
                                    }
                                })}
                            </ul>
                        </Modal.Body>
                    </Modal>
                </div>
            );
        }
    }
}

if (document.getElementById('Version')) {
    ReactDOM.render(<Version />, document.getElementById('Version'));
}
