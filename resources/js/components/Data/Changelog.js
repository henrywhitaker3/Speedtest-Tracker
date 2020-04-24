import React, { Component, version } from 'react';
import ReactDOM from 'react-dom';
import Axios from 'axios';
import { Modal } from 'react-bootstrap';

export default class Changelog extends Component {
    constructor(props) {
        super(props)

        this.state = {
            changelog: {},
            modal: false,
            loading: true,
        }
    }

    componentDidMount = () => {
        this.getChangelog();
    }

    getChangelog = () => {
        Axios.get('/api/update/changelog')
        .then((resp) => {
            this.setState({
                changelog: resp.data.data,
                loading: false
            });
        })
    }

    showModal = () => {
        this.setState({
            modal: true,
        });
    }

    hideModal = () => {
        this.setState({
            modal: false,
        });
    }

    versionList = (key, data) => {
        return (
            <div key={key}>
                <h5>Version: {key}</h5>
                <ul>
                    {data.map((e,i) => {
                        if(e.link == '') {
                            return <li key={key.split('.').join() + i}>{e.description}</li>
                        } else {
                            return <li key={key + i}><a href={e.link} target="_blank" rel="noopener noreferer">{e.description}</a></li>
                        }
                    })}
                </ul>
            </div>
        );
    }

    makeChangelog() {
        var changelog = this.state.changelog;
        var versions = [];

        for(var key in changelog) {
            versions.push(this.versionList(key, changelog[key]));
        }

        return versions;
    }

    render() {
        var show = this.state.modal;
        var loading = this.state.loading;

        if(loading) {
            return <></>
        } else {
            var changelog = this.makeChangelog();
            return (
                <div className="text-muted ml-1 d-inline-block">
                    <i className="ti-link mouse" onClick={this.showModal} />

                    <Modal show={show} onHide={this.hideModal} animation={true}>
                        <Modal.Body>
                            <h3>Changelog:</h3>
                            {changelog}
                        </Modal.Body>
                    </Modal>
                </div>
            );
        }
    }
}

if (document.getElementById('Changelog')) {
    ReactDOM.render(<Changelog />, document.getElementById('Changelog'));
}
