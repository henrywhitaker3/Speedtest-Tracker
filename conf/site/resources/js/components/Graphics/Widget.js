import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Card } from 'react-bootstrap';

export default class Widget extends Component {
    constructor(props) {
        super(props)

        this.state = {
            title: this.props.title,
            value: this.props.value,
            unit: this.props.unit,
            icon: this.props.icon
        }
    }

    componentDidUpdate = () => {
        if(this.props.title != this.state.title || this.props.value != this.state.value || this.props.unit != this.state.unit || this.props.icon != this.state.icon) {
            this.setState({
                title: this.props.title,
                value: this.props.value,
                unit: this.props.unit,
                icon: this.props.icon
            });
        }
    }

    render() {
        var title = this.props.title;
        var value = this.props.value;
        var unit = this.props.unit;
        var icon = this.props.icon;

        switch(icon) {
            case 'ping':
                icon = <span className="ti-pulse icon text-success"></span>;
                break;
            case 'dl':
                icon = <span className="ti-download icon text-warning"></span>;
                break;
            case 'ul':
                icon = <span className="ti-upload icon text-primary"></span>;
                break;
        }

        return (
            <Card className="widget-card shadow-sm">
                <Card.Body>
                    <div className="d-flex align-items-center justify-content-between">
                        <div>
                            <p>{title}</p>
                            <h3 className="d-inline">{value}</h3>
                            <p className="d-inline ml-2">{unit}</p>
                        </div>
                        <div>
                            {icon}
                        </div>
                    </div>
                </Card.Body>
            </Card>
        );
    }
}

if (document.getElementById('Widget')) {
    ReactDOM.render(<Widget />, document.getElementById('Widget'));
}
