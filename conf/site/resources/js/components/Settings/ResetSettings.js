import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Modal, Button } from 'react-bootstrap';
import SettingsModalCard from './SettingsModalCard';
import Axios from 'axios';
import { toast } from 'react-toastify';

export default class ResetSettings extends Component {
    constructor(props) {
        super(props)

        this.state = {
            show: false,
        }
    }

    toggleShow = () => {
        if(this.state.show) {
            this.setState({ show: false });
        } else {
            this.setState({ show:true });
        }
    }

    deleteAll = () => {
        var url = 'api/speedtest/delete/all';

        Axios.delete(url)
        .then((resp) => {
            toast.success('All speedtests have been deleted.');
            this.toggleShow();
        })
        .catch((err) => {
            if(err.response.data.error == undefined) {
                toast.error('Something went wrong.');
            }

            toast.error(err.response.data.error);
        })
    }

    render() {
        var show = this.state.show;
        const title = 'Reset Speedtests';

        return (
            <>
                <SettingsModalCard title={title} description="Bulk delete speedtests from the database." toggleShow={this.toggleShow} />
                <Modal show={show} onHide={this.toggleShow}>
                    <Modal.Header>
                        <Modal.Title>{title}</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        <h4>Clear all speedtests</h4>
                        <p className="text-muted">If using SQLite, a backup of the database will be stored in the location of the current database.</p>
                        <Button onClick={this.deleteAll} variant="danger">Delete all</Button>
                    </Modal.Body>
                </Modal>
            </>
        );
    }
}

if (document.getElementById('ResetSettings')) {
    ReactDOM.render(<ResetSettings />, document.getElementById('ResetSettings'));
}
