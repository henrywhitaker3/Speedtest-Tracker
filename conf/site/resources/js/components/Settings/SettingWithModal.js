import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Card, Form, Button, Modal, Row, Col } from 'react-bootstrap';
import Axios from 'axios';
import { toast } from 'react-toastify';
import SettingsModalCard from '../Settings/SettingsModalCard';

export default class SettingWithModal extends Component {
    constructor(props) {
        super(props)

        this.state = {
            title: this.props.title,
            description: this.props.description,
            settings: this.props.settings,
            show: false,
            autoClose: this.props.autoClose
        }
    }

    ucfirst(string) {
        return string.charAt(0).toUpperCase() + string.slice(1);
    }

    update = () => {
        var url = 'api/settings/bulk?token=' + window.token;
        var data = [];
        var settings = this.state.settings;

        settings.forEach(e => {
            if(e.type !== 'button-get') {
                var res = {
                    name: e.obj.name,
                    value: e.obj.value
                };
                data.push(res);
            }
        });

        data = {
            data: data
        };

        Axios.post(url, data)
        .then((resp) => {
            toast.success(this.state.title + ' updated');
            if(this.state.autoClose) {
                this.toggleShow();
            }
            Axios.get('api/settings/config')
            .then((resp) => {
                window.config = resp.data;
            })
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
                <SettingsModalCard title={title} description={description} toggleShow={this.toggleShow} />
                <Modal show={show} onHide={this.toggleShow}>
                    <Modal.Header closeButton>
                        <Modal.Title>{title}</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        {settings.map((e,i) => {
                            var name = e.obj.name.split('_');
                            name[0] = this.ucfirst(name[0]);
                            name = name.join(' ');

                            if(e.obj.description == null || e.obj.description == '') {
                                var sm = { span: 12 };
                                var md = { span: 12 };
                            } else {
                                var sm = { span: 12 };
                                var md = { span: 6 };
                            }

                            var readonly = false;
                            if(window.config.editable[e.obj.name] == false) {
                                readonly = true;
                            }

                            if(e.type == 'info') {
                                return (
                                    <Row key={e.obj.id} className="d-flex align-items-center">
                                        <Col md={md} sm={sm}>
                                            <p>{e.obj.content}</p>
                                        </Col>
                                    </Row>
                                )
                            } else if(e.type == 'checkbox') {
                                return (
                                    <Row key={e.obj.id} className="d-flex align-items-center">
                                        <Col  md={md} sm={sm}>
                                            <Form.Group controlId={e.obj.name}>
                                                {readonly ?
                                                    <>
                                                        <Form.Check type="checkbox" disabled label={name} defaultChecked={Boolean(Number(e.obj.value))} onInput={this.updateValue} />
                                                        <Form.Text className="text-muted">This setting is defined as an env variable and is not editable.</Form.Text>
                                                    </>
                                                :
                                                    <Form.Check type="checkbox" label={name} defaultChecked={Boolean(Number(e.obj.value))} onInput={this.updateValue} />
                                                }
                                            </Form.Group>
                                        </Col>
                                        {e.description == null &&
                                            <Col md={md} sm={sm}>
                                                <p>{e.obj.description}</p>
                                            </Col>
                                        }
                                    </Row>
                                );
                            } else if(e.type == 'number') {
                                return (
                                    <Row key={e.obj.id}>
                                        <Col md={md} sm={sm}>
                                            <Form.Group controlId={e.obj.name}>
                                                <Form.Label>{name}</Form.Label>
                                                {readonly ?
                                                    <>
                                                        <Form.Control type="number" disabled min={e.min} max={e.max} defaultValue={e.obj.value} onInput={this.updateValue} />
                                                        <Form.Text className="text-muted">This setting is defined as an env variable and is not editable.</Form.Text>
                                                    </>
                                                :
                                                    <Form.Control type="number" min={e.min} max={e.max} defaultValue={e.obj.value} onInput={this.updateValue} />
                                                }
                                            </Form.Group>
                                        </Col>
                                        {e.description == null &&
                                            <Col md={md} sm={sm}>
                                                <p>{e.obj.description}</p>
                                            </Col>
                                        }
                                    </Row>
                                );
                            } else if(e.type == 'text') {
                                return (
                                    <Row key={e.obj.id}>
                                        <Col md={md} sm={sm}>
                                            <Form.Group controlId={e.obj.name}>
                                                <Form.Label>{name}</Form.Label>
                                                {readonly ?
                                                    <>
                                                        <Form.Control type="text" disabled defaultValue={e.obj.value} onInput={this.updateValue} />
                                                        <Form.Text className="text-muted">This setting is defined as an env variable and is not editable.</Form.Text>
                                                    </>
                                                :
                                                    <Form.Control type="text" defaultValue={e.obj.value} onInput={this.updateValue} />
                                                }
                                            </Form.Group>
                                        </Col>
                                        {e.description == null &&
                                            <Col md={md} sm={sm}>
                                                <p dangerouslySetInnerHTML={{ __html: e.obj.description}}></p>
                                            </Col>
                                        }
                                    </Row>
                                );
                            } else if(e.type == 'select') {
                                return (
                                    <Row key={e.obj.id}>
                                        <Col md={md} sm={sm}>
                                            <Form.Group controlId={e.obj.name}>
                                                <Form.Label>{name}</Form.Label>
                                                {readonly ?
                                                    <>
                                                        <Form.Control as="select" disabled defaultValue={e.obj.value} onInput={this.updateValue}>
                                                            {e.options.map((e,i) => {
                                                                return (
                                                                    <option key={i} value={e.value}>{e.name}</option>
                                                                )
                                                            })}
                                                        </Form.Control>
                                                        <Form.Text className="text-muted">This setting is defined as an env variable and is not editable.</Form.Text>
                                                    </>
                                                :
                                                    <Form.Control as="select" defaultValue={e.obj.value} onInput={this.updateValue}>
                                                        {e.options.map((e,i) => {
                                                            return (
                                                                <option key={i} value={e.value}>{e.name}</option>
                                                            )
                                                        })}
                                                    </Form.Control>
                                                }
                                            </Form.Group>
                                        </Col>
                                        {e.description == null &&
                                            <Col md={md} sm={sm}>
                                                <p>{e.obj.description}</p>
                                            </Col>
                                        }
                                    </Row>
                                )
                            } else if(e.type == 'button-get') {
                                return (
                                    <Row key={e.obj.id}>
                                        <Col md={md} sm={sm}>
                                            <p>{name}</p>
                                            <Button onClick={() => { Axios.get(e.url) }} >{name}</Button>
                                        </Col>
                                        {e.description == null &&
                                            <Col md={md} sm={sm}>
                                                <p>{e.obj.description}</p>
                                            </Col>
                                        }
                                    </Row>
                                )
                            } else if(e.type == 'group') {
                                return (
                                    <div key={e.obj.id}>
                                        <Row>
                                            <Col md={md} sm={sm}>
                                                <p className="mb-0">{name}</p>
                                            </Col>
                                            {e.description == null &&
                                                <Col md={md} sm={sm}>
                                                    <p>{e.obj.description}</p>
                                                </Col>
                                            }
                                        </Row>
                                        <Row>
                                            <Col sm={{ span: 12 }}>
                                                {e.children.map((ee,ii) => {
                                                    if(ee.type == 'button-get') {
                                                        return (
                                                            <Button key={ii} variant={ee.btnType} className={'mr-2 mb-3'} onClick={() => { Axios.get(ee.url)
                                                                                                                                                .then((resp) => { toast.success('Healthcheck sent') })
                                                                                                                                                .catch((resp) => { resp = resp.response; toast.error(resp.data.error) })
                                                            }} >{ee.text}</Button>
                                                        )
                                                    }
                                                })}
                                            </Col>
                                        </Row>
                                    </div>
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
