import React, { Component } from 'react';
import ReactDOM from 'react-dom';

export default class Index extends Component {
    constructor(props) {
        super(props)

        this.state = {

        }
    }

    render() {
        return (
            <div>
            </div>
        );
    }
}

if (document.getElementById('main')) {
    ReactDOM.render(<Index />, document.getElementById('main'));
}
