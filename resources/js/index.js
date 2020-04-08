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
            redirectLogin: true,
            redirectHome: false,
            token: null,
            user: null,
        }
    }

    componentDidMount() {
        this.lookForToken();
    }

    lookForToken = () => {
        var token = JSON.parse(localStorage.getItem('token'));
        if(token == null) {
            this.setState({
                loading: false,
            });
        } else {
            this.tryToken(token)
        }
    }

    setToken = (token) => {
        localStorage.setItem('token', JSON.stringify(token));
        this.setState({
            loading: true,
            token: token,
        });
        this.tryToken(token);
    }

    tryToken = (token, reload = true) => {
        Axios.get('/api/auth/me?token=' + token.access_token)
        .then((resp) => {
            this.setState({
                user: resp.data,
                loading: false,
                redirectLogin: false,
            });
            if(reload) {
                this.setState({
                    redirectHome: true
                });
            }
        })
        .catch((err) => {
            console.log('Invalid token');
            console.log(err);
            this.setState({
                redirectLogin: true
            });
        })
    }

    render() {
        var loading = this.state.loading;
        var redirectLogin = this.state.redirectLogin;
        var redirectHome = this.state.redirectHome;
        var user = this.state.user;
        var token = this.state.token;

        return (
            <div>
                {loading ?
                <div>
                    <Loader />
                </div>
                :
                    <div>
                        <BrowserRouter>
                            <Route render={(props) => (<ToastContainer />)} />
                            <Switch>
                                    <Route exact path="/" render={(props) => (
                                        <div>
                                            <HomePage user={user} token={token} />
                                        </div>
                                    )} />
                                    <Route exact path="/login" render={(props) => (
                                        <div>
                                            <Login setToken={this.setToken} />
                                        </div>
                                    )} />
                                    <Route exact path="/error/:code" render={(props) => ( <ErrorPage code={props.match.params.code} /> )} />
                                    <Route render={(props) => (<ErrorPage code="404" />)} />
                            </Switch>
                            {redirectLogin &&
                                <Redirect to="/login"></Redirect>
                            }
                            {redirectHome &&
                                <Redirect to="/"></Redirect>
                            }
                        </BrowserRouter>
                    </div>
                }
            </div>
        );
    }
}

if (document.getElementById('main')) {
    ReactDOM.render(<Index />, document.getElementById('main'));
}
