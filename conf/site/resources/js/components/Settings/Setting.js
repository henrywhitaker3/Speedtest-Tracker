import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Card, Form, Button } from 'react-bootstrap';
import Axios from 'axios';
import { toast } from 'react-toastify';

export default class Setting extends Component {
    constructor(props) {
        super(props)

        this.state = {
            name: this.props.name,
            value: this.props.value,
            description: this.props.description,
        }
    }

    ucfirst(string) {
        return string.charAt(0).toUpperCase() + string.slice(1);
    }

    update = () => {
        var url = 'api/settings?token=' + window.token;
        var data = {
            name: this.state.name,
            value: this.state.value
        };

        Axios.post(url, data)
        .then((resp) => {
            toast.success(this.ucfirst(this.state.name) + ' updated');
        })
        .catch((err) => {
            if(err.response.status == 422) {
                var errors = err.response.data.error;
                for(var key in errors) {
                    var error = errors[key];
                    toast.error(error[0])
                }
            } else {
                toast.error('Something went wrong')
            }
        })
    }

    updateValue = (e) => {
        this.setState({
            value: e.target.value
        });
    }

    render() {
        var name = this.state.name;
        var value = this.state.value;
        var description = this.state.description;

        return (
            <Card className="m-2 setting-card">
                <Card.Body className="d-flex align-items-center">
                    <div>
                        <h4>{this.ucfirst(name)}</h4>
                        <div dangerouslySetInnerHTML={{ __html: description}} />
                        <Form.Group controlId={name}>
                            <Form.Label>{this.ucfirst(name)}</Form.Label>
                            <Form.Control type="text" label={name} defaultValue={value} onInput={this.updateValue} />
                        </Form.Group>
                        <Button variant="primary" onClick={this.update}>Save</Button>
                    </div>
                </Card.Body>
            </Card>
        );
    }
}

if (document.getElementById('Setting')) {
    ReactDOM.render(<Setting />, document.getElementById('Setting'));
}
