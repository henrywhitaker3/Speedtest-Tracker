import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter, Switch, Route, Redirect, useHistory } from 'react-router-dom';
import Axios from 'axios';
import ErrorPage from './components/ErrorPage';
import Loader from './components/Loader';
import Login from './components/Login';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import HomePage from './components/Home/HomePage';

export default class Index extends Component {
    constructor(props) {
        super(props)

        this.state = {
            loading: true,
            redirect: false
        }
    }

    componentDidMount = () => {
        this.getConfig();
    }

    getConfig() {
        var url = 'api/settings/config';

        Axios.get(url)
        .then((resp) => {
            window.config = resp.data;
            this.setState({
                loading: false,
                redirect: true,
            });
            if(window.location.pathname != window.config.base) {
                window.location.href = window.config.base;
            }
        })
    }

    render() {
        var loading = this.state.loading;
        var redirect = this.state.redirect;
        var baseSet = this.isset(window.config);

        if(loading) {
            return (
                <Loader />
            );
        } else {
            if(baseSet) {
                return (
                    <BrowserRouter>
                        <Route render={(props) => (<ToastContainer />)} />
                        <Switch>
                                <Route exact path={window.config.base} render={(props) => (
                                    <div>
                                        <HomePage />
                                    </div>
                                )} />
                                <Route exact path={window.config.base + "error/:code"} render={(props) => ( <ErrorPage code={props.match.params.code} /> )} />
                                <Route render={(props) => (<ErrorPage code="404" />)} />
                        </Switch>
                    </BrowserRouter>
                );
            } else {
                return (
                    <Loader />
                );
            }
        }
    }

    isset(v) {
        if(typeof v !== "undefined" || v !== null) {
            return true;
        } else {
            return false;
        }
    }
}

if (document.getElementById('main')) {
    ReactDOM.render(<Index />, document.getElementById('main'));
}
