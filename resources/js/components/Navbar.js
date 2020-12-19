import React, { Component } from 'react';
import {Nav, Navbar as BootstrapNavbar, NavLink as BootstrapNavLink} from 'react-bootstrap';
import ReactDOM from 'react-dom';
import { Link, NavLink } from 'react-router-dom';

export default class Navbar extends Component {
    constructor(props) {
        super(props)

        this.state = {
            brand: {
                name: "Speedtest Tracker",
                url: window.config.base
            },
            pages: [
                {
                    name: 'Home',
                    url: window.config.base,
                    authRequired: false
                },
                {
                    name: 'Settings',
                    url: window.config.base + 'settings',
                    authRequired: true
                }
            ]
        }
    }

    generateLinks = () => {
        return this.state.pages.map(page => {
            if(
                page.authRequired === false ||
                (
                    page.authRequired === true &&
                    window.config.auth &&
                    window.authenticated
                ) ||
                (
                    page.authRequired === true &&
                    window.config.auth === false
                )
            ) {
                return <BootstrapNavLink key={page.url} as={NavLink} to={page.url}>{page.name}</BootstrapNavLink>;
            }
        });
    }

    render() {
        var brand = this.state.brand;
        var pages = this.generateLinks();

        return (
            <BootstrapNavbar variant="dark" bg="dark">
                <BootstrapNavbar.Brand as={Link} to={brand.url}>{brand.name}</BootstrapNavbar.Brand>
                <BootstrapNavbar.Toggle  aria-controls="basic-navbar-nav" />
                <BootstrapNavbar.Collapse id="basic-navbar-nav">
                    <Nav className="ml-auto">
                        {pages}
                    </Nav>
                </BootstrapNavbar.Collapse>
            </BootstrapNavbar>
        );
    }
}

if (document.getElementById('navbar')) {
    ReactDOM.render(<Navbar />, document.getElementById('navbar'));
}
