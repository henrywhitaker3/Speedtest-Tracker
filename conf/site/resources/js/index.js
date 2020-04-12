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
        }
    }

    render() {
        var loading = this.state.loading;

        return (
            <BrowserRouter>
                <Route render={(props) => (<ToastContainer />)} />
                <Switch>
                        <Route exact path="/" render={(props) => (
                            <div>
                                <HomePage />
                            </div>
                        )} />
                        <Route exact path="/error/:code" render={(props) => ( <ErrorPage code={props.match.params.code} /> )} />
                        <Route render={(props) => (<ErrorPage code="404" />)} />
                </Switch>
            </BrowserRouter>
        );
    }
}

if (document.getElementById('main')) {
    ReactDOM.render(<Index />, document.getElementById('main'));
}
