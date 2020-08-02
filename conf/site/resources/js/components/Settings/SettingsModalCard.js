import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Card, Button } from 'react-bootstrap';

export default class SettingsModalCard extends Component {
    constructor(props) {
        super(props)

        this.state = {
            title: this.props.title,
            description: this.props.description,
            toggleShow: this.props.toggleShow,
        }
    }

    render() {
        var title = this.state.title;
        var description = this.state.description;
        var toggleShow = this.state.toggleShow;

        return (
            <Card className="m-2 setting-card">
                <Card.Body className="d-flex align-items-center">
                    <div>
                        <h4>{title}</h4>
                        <p>{description}</p>
                        <Button variant="primary" onClick={toggleShow}>Edit</Button>
                    </div>
                </Card.Body>
            </Card>
        );
    }
}

if (document.getElementById('SettingModalCard')) {
    ReactDOM.render(<SettingsModalCard />, document.getElementById('SettingModalCard'));
}
