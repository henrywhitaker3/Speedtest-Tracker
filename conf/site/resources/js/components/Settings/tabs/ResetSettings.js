import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Button } from 'react-bootstrap';
import Axios from 'axios';
import { toast } from 'react-toastify';

export default class ResetSettings extends Component {
    constructor(props) {
        super(props)

        this.state = {
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
                <h4>Clear all speedtests</h4>
                <p className="text-muted">If using SQLite, a backup of the database will be stored in the location of the current database.</p>
                <Button onClick={this.deleteAll} variant="danger">Delete all</Button>
            </>
        );
    }
}

if (document.getElementById('ResetSettings')) {
    ReactDOM.render(<ResetSettings />, document.getElementById('ResetSettings'));
}
