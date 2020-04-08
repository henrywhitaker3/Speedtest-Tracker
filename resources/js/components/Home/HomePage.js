import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import HistoryGraph from '../Graphics/HistoryGraph';

export default class HomePage extends Component {
    constructor(props) {
        super(props)

        this.state = {
            token: this.props.token,
        }
    }

    render() {
        var token = this.state.token;

        return (
            <div>
                <HistoryGraph token={token} />
            </div>
        );
    }
}

if (document.getElementById('homePage')) {
    ReactDOM.render(<HomePage />, document.getElementById('homePage'));
}
