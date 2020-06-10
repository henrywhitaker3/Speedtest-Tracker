import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import Axios from 'axios';
import { toast } from 'react-toastify';
import { Modal, ProgressBar } from 'react-bootstrap';
import { Button } from 'react-bootstrap';
import Changelog from '../Data/Changelog';

export default class Version extends Component {
    constructor(props) {
        super(props)

        this.state = {
            version: document.querySelector('meta[name="version"]').content,
            update: false,
            modalShow: false,
            changelog: [],
            showProgress: false,
            updateProgress: 0,
        };
    }

    componentDidMount() {
        // this.checkForUpdates();
    }

    checkForUpdates = () => {
        var url = 'api/update/check';

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

    updateApp = () => {
        this.setState({
            showProgress: true,
            updateProgress: 0,
        });
        toast.info('Downloading update');
        Axios.get('api/update/download')
        .then((resp) => {
            this.setState({
                updateProgress: 50,
            });
            toast.info('Extracting update');
            Axios.get('api/speedtest/extract')
            .then((resp) => {
                this.setState({
                    updateProgress: 75,
                });
                toast.info('Applying update');
                Axios.get('api/update/move')
                .then((resp) => {
                    this.setState({
                        updateProgress: 100,
                    });
                    toast.success('Update successful. Refreshing your page...');
                    setTimeout(function() {
                        location.reload(true);
                    }, 5000);
                })
            })
        })
        .catch((err) => {
            toast.error('Something went wrong...');
        })
    }

    render() {
        var version = this.state.version;
        var update = this.state.update;
        var modalShow = this.state.modalShow;
        var changelog = this.state.changelog;
        var showProgress = this.state.showProgress;
        var updateProgress = this.state.updateProgress;

        if(update === false) {
            return (
                <div>
                    <p className="text-muted mb-0 d-inline-block">Speedtest Tracker Version: {version}</p>
                    <Changelog />
                </div>
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
                                        return (
                                            <li key={i}><a href={e.link} target="_blank" rel="noopener noreferer">{e.description}</a></li>
                                        )
                                    }
                                })}
                            </ul>
                            {showProgress &&
                                <div>
                                    <p>Update progress:</p>
                                    <ProgressBar animated now={updateProgress} />
                                </div>
                            }
                            {!showProgress &&
                                <Button variant="primary" onClick={this.updateApp}>Update</Button>
                            }
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
