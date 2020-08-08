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
            icon: this.props.icon,
            avg: this.props.avg,
            max: this.props.max,
            failed: this.props.failed,
        }
    }

    componentDidUpdate = () => {
        if(this.props.title != this.state.title || this.props.value != this.state.value || this.props.unit != this.state.unit || this.props.icon != this.state.icon || this.props.avg != this.state.avg || this.props.max != this.state.max || this.props.failed != this.state.failed) {
            this.setState({
                title: this.props.title,
                value: this.props.value,
                unit: this.props.unit,
                icon: this.props.icon,
                avg: this.props.avg,
                max: this.props.max,
                failed: this.props.failed,
            });
        }
    }

    render() {
        var title = this.state.title;
        var value = this.state.value;
        var unit = this.state.unit;
        var icon = this.state.icon;
        var max = this.state.max;
        var avg = this.state.avg;
        var failed = Boolean(Number(this.state.failed));

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
                    <div>
                        <div>
                            <div className="d-flex align-items-center justify-content-between">
                                <h4>{title}</h4>
                                {icon}
                            </div>

                            <div className="text-truncate">
                                <h3 className="d-inline">{(!failed) ? value : <span className="ti-close text-danger"></span> }</h3>
                                <p className="d-inline ml-2">{unit} (current)</p>
                            </div>
                            <div className="text-muted text-truncate">
                                <h5 className="d-inline">{avg}</h5>
                                <p className="d-inline ml-2">{unit} (average)</p>
                            </div>
                            <div className="text-muted text-truncate">
                                <h5 className="d-inline">{max}</h5>
                                <p className="d-inline ml-2">{unit} (maximum)</p>
                            </div>
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
