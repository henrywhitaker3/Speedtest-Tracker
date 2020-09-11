import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Card } from 'react-bootstrap';

export default class Widget extends Component {
    constructor(props) {
        super(props)

        this.state = {
            title: this.props.title,
            unit: this.props.unit,
            icon: this.props.icon,
            failed: this.props.failed,
            data: this.props.data
        }
    }

    parseData(title, data) {
        var returnData = {};


        if(title == 'Ping') {
            returnData.value = parseFloat(data.data.ping).toFixed(1);

            if(window.config.widgets.show_average) {
                returnData.avg = parseFloat(data.average.ping).toFixed(1);
            }

            if(window.config.widgets.show_max) {
                returnData.max = parseFloat(data.maximum.ping).toFixed(1);
            }

            if(window.config.widgets.show_min) {
                returnData.min = parseFloat(data.minimum.ping).toFixed(1);
            }
        }

        if(title == 'Upload') {
            returnData.value = parseFloat(data.data.upload).toFixed(1);

            if(window.config.widgets.show_average) {
                returnData.avg = parseFloat(data.average.upload).toFixed(1);
            }

            if(window.config.widgets.show_max) {
                returnData.max = parseFloat(data.maximum.upload).toFixed(1);
            }

            if(window.config.widgets.show_min) {
                returnData.min = parseFloat(data.minimum.upload).toFixed(1);
            }
        }

        if(title == 'Download') {
            returnData.value = parseFloat(data.data.download).toFixed(1);

            if(window.config.widgets.show_average) {
                returnData.avg = parseFloat(data.average.download).toFixed(1);
            }

            if(window.config.widgets.show_max) {
                returnData.max = parseFloat(data.maximum.download).toFixed(1);
            }

            if(window.config.widgets.show_min) {
                returnData.min = parseFloat(data.minimum.download).toFixed(1);
            }
        }

        return returnData;
    }

    componentDidUpdate = () => {
        if(this.props.title != this.state.title || this.props.data != this.state.data || this.props.unit != this.state.unit || this.props.icon != this.state.icon || this.props.failed != this.state.failed) {
            this.setState({
                title: this.props.title,
                unit: this.props.unit,
                icon: this.props.icon,
                failed: this.props.failed,
                data: this.props.data
            });
        }
    }

    render() {
        var title = this.state.title;
        var unit = this.state.unit;
        var icon = this.state.icon;
        var failed = Boolean(Number(this.state.failed));

        var data = this.parseData(title, this.state.data);

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
                                <h3 className="d-inline">{(!failed) ? data.value : <span className="ti-close text-danger"></span> }</h3>
                                <p className="d-inline ml-2">{unit} (current)</p>
                            </div>

                            {window.config.widgets.show_average &&
                                <div className="text-muted text-truncate">
                                    <h5 className="d-inline">{data.avg}</h5>
                                    <p className="d-inline ml-2">{unit} (average)</p>
                                </div>
                            }

                            {window.config.widgets.show_max &&
                                <div className="text-muted text-truncate">
                                    <h5 className="d-inline">{data.max}</h5>
                                    <p className="d-inline ml-2">{unit} (maximum)</p>
                                </div>
                            }

                            {window.config.widgets.show_min &&
                                <div className="text-muted text-truncate">
                                    <h5 className="d-inline">{data.min}</h5>
                                    <p className="d-inline ml-2">{unit} (minimum)</p>
                                </div>
                            }
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
