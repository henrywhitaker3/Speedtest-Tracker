import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import HistoryGraph from '../Graphics/HistoryGraph';
import LatestResults from '../Graphics/LatestResults';
import Footer from './Footer';

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
                <div className="my-4">
                    <LatestResults token={token} />
                    <HistoryGraph token={token} />
                </div>
                <Footer />
            </div>
        );
    }
}

if (document.getElementById('homePage')) {
    ReactDOM.render(<HomePage />, document.getElementById('homePage'));
}
