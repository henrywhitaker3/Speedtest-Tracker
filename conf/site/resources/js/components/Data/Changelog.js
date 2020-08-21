import React, { Component, version } from 'react';
import ReactDOM from 'react-dom';
import Axios from 'axios';
import { Modal, Collapse, Button } from 'react-bootstrap';

export default class Changelog extends Component {
    constructor(props) {
        super(props)

        this.state = {
            changelog: {},
            modal: false,
            loading: true,
            hidden: false,
        }
    }

    componentDidMount = () => {
        if( (window.config.auth == true && window.authenticated == true) || window.config.auth == false) {
            this.getChangelog();
        }
    }

    getChangelog = () => {
        Axios.get('api/update/changelog?token=' + window.token)
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

    toggleHidden = () => {
        var hidden = this.state.hidden;
        if(hidden) {
            this.setState({
                hidden: false
            });
        } else {
            this.setState({
                hidden: true
            });
        }
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
        var versionsVis = [];
        var versionsHid = [];

        var i = 0;
        for(var key in changelog) {
            if(i <= 5) {
                versionsVis.push(this.versionList(key, changelog[key]));
            } else {
                versionsHid.push(this.versionList(key, changelog[key]));
            }
            i++;
        }

        return {
            visible: versionsVis,
            hidden: versionsHid
        };
    }

    render() {
        var show = this.state.modal;
        var loading = this.state.loading;
        var showHidden = this.state.hidden;

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
                            {changelog.visible}
                            {changelog.hidden.length > 5 &&
                                <>
                                    <Collapse in={showHidden}>
                                        <div>
                                            {changelog.hidden}
                                        </div>
                                    </Collapse>
                                    <div className="w-100 text-center">
                                    {showHidden ?
                                        <Button variant="primary" className="mx-auto mouse" onClick={this.toggleHidden}>Show less</Button>
                                    :
                                        <Button variant="primary" className="mx-auto mouse" onClick={this.toggleHidden}>Show more</Button>
                                    }
                                    </div>
                                </>
                            }
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
