import React, { Component } from 'react';
import { Form } from 'react-bootstrap';
import ReactDOM from 'react-dom';

export default class SettingsInput extends Component {
    constructor(props) {
        super(props)

        this.state = {
            type: this.props.type,
            name: this.props.name,
            value: (this.props.value) ? this.props.value : '',
            classes: this.props.classes,
            id: this.props.id,
            label: (this.props.label) ? this.props.label : false,
            readonly: true,
            description: (this.props.description) ? this.props.description : false,
        }
    }

    componentDidMount() {
        this.setState({
            readonly: this.isReadOnly()
        });
    }

    handleInput = (evt) => {
        var val = evt.target.value;

        if(this.state.type === 'checkbox') {
            val = e.target.checked;
        }

        this.props.handleInput(
            this.state.name.split(' ').join('_'),
            val
        );

        this.setState({
            value: val
        });
    }

    isReadOnly = () => {
        if(window.config.editable[this.state.name] == false) {
            return true;
        }

        return false;
    }

    generateInput = () => {
        var disabled = (this.state.readonly) ? true : false;

        return <Form.Control
                    name={this.state.name}
                    type={this.state.type}
                    defaultValue={this.state.value}
                    disabled={disabled}
                    onInput={this.handleInput} />
    }

    render() {
        var input = this.generateInput();
        var id = this.state.id;
        var readonly = this.state.readonly;
        var label = this.state.label;
        var description = this.state.description;

        return (
            <Form.Group controlId={id}>
                {label &&
                    <Form.Label dangerouslySetInnerHTML={{ __html: label }} />
                }

                {input}

                {description &&
                    <p dangerouslySetInnerHTML={{ __html: description }}></p>
                }

                {readonly &&
                    <Form.Text className="text-muted">This setting is defined as an env variable and is not editable.</Form.Text>
                }
            </Form.Group>
        );
    }
}

if (document.getElementById('SettingsInput')) {
    ReactDOM.render(<SettingsInput />, document.getElementById('SettingsInput'));
}
