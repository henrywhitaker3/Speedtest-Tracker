import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import TestsTable from './Graphics/TestsTable';
import Footer from './Home/Footer';
import Navbar from './Navbar';

export default class SpeedtestsPage extends Component {
    constructor(props) {
        super(props)

        this.state = {

        }
    }

    render() {
        return (
            <div>
                <Navbar />
                <TestsTable />
                <Footer />
            </div>
        );
    }
}

if (document.getElementById('SpeedtestsPage')) {
    ReactDOM.render(<SpeedtestsPage />, document.getElementById('SpeedtestsPage'));
}
