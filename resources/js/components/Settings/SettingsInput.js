import React, { Component } from 'react';
import { Form } from 'react-bootstrap';
import ReactDOM from 'react-dom';

export default class SettingsInput extends Component {
    constructor(props) {
        super(props)

        this.state = {
            type: this.props.type,
            name: this.props.name,
            displayName: (this.props.name) ? this.formatName(this.props.name) : '',
            value: (this.props.value) ? this.props.value : '',
            classes: this.props.classes,
            id: this.props.id,
            label: (this.props.label) ? this.props.label : false,
            readonly: true,
            description: (this.props.description) ? this.props.description : false,
            options: this.props.options ? this.props.options : [],
            hideDescription: this.props.hideDescription ? true : false,
            min: this.props.min ? this.props.min : null,
            max: this.props.max ? this.props.max : null,
            url: this.props.url,
            inline: this.props.inline ? 'd-inline-block' : 'd-block',
            btnType: this.props.btnType,
            earlyReturn: this.props.earlyReturn ? true : false,
            autoComplete: String(this.props.autoComplete ? true : Math.random().toString(36).replace(/[^a-z]+/g, '').substr(0, 7)),
        }
    }

    componentDidMount() {
        this.setState({
            readonly: this.isReadOnly()
        });
    }

    formatName(name) {
        name = name.split('_').join(' ');

        return name.charAt(0).toUpperCase() + name.slice(1);
    }

    handleInput = (evt) => {
        var val = evt.target.value;

        if(this.state.type === 'checkbox') {
            val = evt.target.checked;
        }

        this.props.handler(
            this.state.name,
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

    generateNumberInput(disabled) {
        return <Form.Control
            name={this.state.name}
            type={this.state.type}
            defaultValue={this.state.value}
            disabled={disabled}
            min={this.state.min}
            max={this.state.max}
            onInput={this.handleInput}
            autoComplete={this.state.autoComplete} />
    }

    generateSelectInput(disabled) {
        return (
            <Form.Control
                as="select"
                name={this.state.name}
                type={this.state.type}
                defaultValue={this.state.value}
                disabled={disabled}
                onInput={this.handleInput}
            >
                {this.state.options.map((option,i) => {
                    return <option key={i} value={option.value}>{option.name}</option>
                })}
            </Form.Control>
        );
    }

    generateCheckboxInput(disabled) {
        return <Form.Control
            custom
            className="ml-2"
            name={this.state.name}
            type={this.state.type}
            defaultChecked={this.state.value}
            disabled={disabled}
            onInput={this.handleInput} />
    }

    generateTextInput(disabled) {
        return <Form.Control
            name={this.state.name}
            type={this.state.type}
            defaultValue={this.state.value}
            disabled={disabled}
            onInput={this.handleInput}
            autoComplete={this.state.autoComplete} />
    }

    generatePasswordInput(disabled) {
        return <Form.Control
            name={this.state.name}
            type={this.state.type}
            defaultValue={this.state.value}
            disabled={disabled}
            onInput={this.handleInput}
            autoComplete={this.state.autoComplete} />
    }

    generateButtonGetInput() {
        var url = this.state.url;

        return (
            <button
                type="button"
                className={"btn btn-" + this.state.btnType + ' ' + this.state.inline + ' ' + this.state.classes}
                onClick={() => {
                    window.axios.get(url)
                }}
            >{this.state.displayName}</button>
        );
    }

    generateInput = () => {
        var disabled = (this.state.readonly) ? true : false;
        var input = null;

        if(this.state.type === 'number') {
            input = this.generateNumberInput(disabled);
        }

        if(this.state.type === 'select') {
            input = this.generateSelectInput(disabled);
        }

        if(this.state.type === 'checkbox') {
            input = this.generateCheckboxInput(disabled);
        }

        if(this.state.type === 'text') {
            input = this.generateTextInput(disabled);
        }

        if(this.state.type === 'password') {
            input = this.generatePasswordInput(disabled);
        }

        if(this.state.type === 'btn-get') {
            input = this.generateButtonGetInput();
        }

        if(this.state.earlyReturn) {
            return input;
        }

        return (
            <Form.Group controlId={this.state.id}>
                {this.state.label &&
                    <Form.Label style={{fontSize: '1.25rem'}}>{this.formatName(this.state.name)}</Form.Label>
                }

                {input}

                {this.state.description && !this.state.hideDescription &&
                    <p className="mt-1 text-muted" dangerouslySetInnerHTML={{ __html: this.state.description }}></p>
                }

                {this.state.readonly &&
                    <Form.Text className="text-muted">This setting is defined as an env variable and is not editable.</Form.Text>
                }
            </Form.Group>
        );
    }

    render() {
        var input = this.generateInput();

        return input;
    }
}

if (document.getElementById('SettingsInput')) {
    ReactDOM.render(<SettingsInput />, document.getElementById('SettingsInput'));
}
