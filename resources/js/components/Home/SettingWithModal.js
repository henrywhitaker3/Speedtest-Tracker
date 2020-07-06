import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Card, Form, Button, Modal, Row, Col } from 'react-bootstrap';
import Axios from 'axios';
import { toast } from 'react-toastify';

export default class SettingWithModal extends Component {
    constructor(props) {
        super(props)

        this.state = {
            title: this.props.title,
            description: this.props.description,
            settings: this.props.settings,
            show: false
        }
    }

    ucfirst(string) {
        return string.charAt(0).toUpperCase() + string.slice(1);
    }

    update = () => {
        var url = 'api/settings/bulk';
        var data = [];
        var settings = this.state.settings;

        settings.forEach(e => {
            var res = {
                name: e.obj.name,
                value: e.obj.value
            };
            data.push(res);
        });

        data = {
            data: data
        };

        Axios.post(url, data)
        .then((resp) => {
            toast.success(this.state.title + ' updated');
            this.toggleShow();
        })
        .catch((err) => {
            if(err.response.status == 422) {
                toast.error('Your input was invalid');
            } else {
                toast.error('Something went wrong')
            }
        })
    }

    updateValue = (e) => {
        var name = e.target.id;
        if(e.target.type == 'checkbox') {
            var val = e.target.checked;
        } else {
            var val = e.target.value;
        }
        var settings = this.state.settings;
        var i = 0;
        settings.forEach(ele => {
            if(ele.obj.name == name) {
                ele.obj.value = val;
            }
            settings[i] = ele;
            i++;
        });
        this.setState({
            settings: settings
        });
    }

    toggleShow = () => {
        var show = this.state.show;
        if(show) {
            this.setState({
                show: false
            });
        } else {
            this.setState({
                show: true
            });
        }
    }

    render() {
        var title = this.state.title;
        var description = this.state.description;
        var show = this.state.show;
        var settings = this.state.settings;

        return (
            <>
                <Card className="m-2 setting-card">
                    <Card.Body className="d-flex align-items-center">
                        <div>
                            <h4>{title}</h4>
                            <p>{description}</p>
                            <Button variant="primary" onClick={this.toggleShow}>Edit</Button>
                        </div>
                    </Card.Body>
                </Card>
                <Modal show={show} onHide={this.toggleShow}>
                    <Modal.Header>
                        <Modal.Title>{title}</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        {settings.map((e,i) => {
                            var name = e.obj.name.split('_');
                            name[0] = this.ucfirst(name[0]);
                            name = name.join(' ');
                            if(e.type == 'checkbox') {
                                return (
                                    <Row key={e.obj.id} className="d-flex align-items-center">
                                        <Col  md={{ span: 6 }} sm={{ span: 12 }}>
                                            <Form.Group controlId={e.obj.name}>
                                                <Form.Check type="checkbox" label={name} defaultChecked={Boolean(Number(e.obj.value))} onInput={this.updateValue} />
                                            </Form.Group>
                                        </Col>
                                        <Col md={{ span: 6 }} sm={{ span: 12 }}>
                                            <p>{e.obj.description}</p>
                                        </Col>
                                    </Row>
                                );
                            } else if(e.type == 'number') {
                                return (
                                    <Row key={e.obj.id}>
                                        <Col md={{ span: 6 }} sm={{ span: 12 }}>
                                            <Form.Group controlId={e.obj.name}>
                                                <Form.Label>{name}</Form.Label>
                                                <Form.Control type="number" min={e.min} max={e.max} defaultValue={e.obj.value} onInput={this.updateValue} />
                                            </Form.Group>
                                        </Col>
                                        <Col md={{ span: 6 }} sm={{ span: 12 }}>
                                            <p>{e.obj.description}</p>
                                        </Col>
                                    </Row>
                                );
                            } else if(e.type == 'select') {
                                return (
                                    <Row key={e.obj.id}>
                                        <Col md={{ span: 6 }} sm={{ span: 12 }}>
                                            <Form.Group controlId={e.obj.name}>
                                                <Form.Label>{name}</Form.Label>
                                                <Form.Control as="select" defaultValue={e.obj.value} onInput={this.updateValue}>
                                                    {e.options.map((e,i) => {
                                                        return (
                                                            <option key={i} value={e.value}>{e.name}</option>
                                                        )
                                                    })}
                                                </Form.Control>
                                            </Form.Group>
                                        </Col>
                                        <Col md={{ span: 6 }} sm={{ span: 12 }}>
                                            <p>{e.obj.description}</p>
                                        </Col>
                                    </Row>
                                )
                            }
                        })}
                        <Button variant="primary" type="submit" onClick={this.update} >Save</Button>
                    </Modal.Body>
                </Modal>
            </>
        );
    }
}

if (document.getElementById('Setting')) {
    ReactDOM.render(<Setting />, document.getElementById('Setting'));
}
