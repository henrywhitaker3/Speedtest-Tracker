import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import HistoryGraph from '../Graphics/HistoryGraph';
import LatestResults from '../Graphics/LatestResults';
import Footer from './Footer';
import DataRow from '../Data/DataRow';

export default class HomePage extends Component {

    render() {
        return (
            <div>
                <div className="my-4">
                    <LatestResults />
                    <HistoryGraph />
                    <DataRow />
                </div>
                <Footer />
            </div>
        );
    }
}

if (document.getElementById('homePage')) {
    ReactDOM.render(<HomePage />, document.getElementById('homePage'));
}
