import React, { Component } from 'react';
import ReactDOM from 'react-dom';

export default class Navbar extends Component {
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

if (document.getElementById('navbar')) {
    ReactDOM.render(<Navbar />, document.getElementById('navbar'));
}
