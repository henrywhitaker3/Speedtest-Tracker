import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import HistoryGraph from '../Graphics/HistoryGraph';
import LatestResults from '../Graphics/LatestResults';
import Footer from './Footer';
import DataRow from '../Data/DataRow';
import TestsTable from '../Graphics/TestsTable';
import Login from '../Login';
import Authentication from '../Authentication/Authentication';
import Navbar from '../Navbar';
import axios from 'axios';

export default class HomePage extends Component {
    constructor(props) {
        super(props)

        this.state = {
            latest: null,
            time: null,
            fail: null,
            config: null,
            days: 7,
            interval: null,
        }
    }

    componentDidMount = () => {
        this.getData();
        var interval = setInterval(this.getData, 10000);
        this.setState({
            interval: interval,
        });
    }

    componentWillUnmount() {
        clearInterval(this.state.interval);
    }

    updateDays = (days) => {
        this.setState({ days: days });
        this.getData();
    }


    getData = () => {
        axios.get('api/speedtest/home/' + this.state.days)
        .then((resp) => {
            this.setState({
                latest: resp.data.latest,
                time: resp.data.time,
                fail: resp.data.fail,
                config: resp.data.config
            });
        })
        .catch((err) => {
            console.log(err);
        })
    }

    render() {
        let latest = this.state.latest;
        let time = this.state.time;
        let fail = this.state.fail;
        let config = this.state.config;
        let days = this.state.days;

        return (
            <div>
                <Navbar />
                <div className="my-4">
                    {(window.config.auth == true && window.authenticated == false) &&
                        <Login />
                    }
                    <LatestResults data={latest} />
                    <HistoryGraph updateDays={this.updateDays} dlUl={time} fail={fail} config={config} days={days} />
                </div>
                <Footer />
            </div>
        );
    }
}

if (document.getElementById('homePage')) {
    ReactDOM.render(<HomePage />, document.getElementById('homePage'));
}
