import React, { Component } from 'react';
import {Nav, Navbar as BootstrapNavbar, NavLink as BootstrapNavLink} from 'react-bootstrap';
import ReactDOM from 'react-dom';
import { Link, NavLink } from 'react-router-dom';

export default class Navbar extends Component {
    constructor(props) {
        super(props)

        this.state = {
            brand: {
                name: window.config.name,
                url: window.config.base
            },
        }
    }

    generatePagesArray() {
        var pages = [
            {
                name: 'Home',
                url: window.config.base,
                authRequired: false
            },
            {
                name: 'All Tests',
                url: window.config.base + 'speedtests',
                authRequired: false
            },
            {
                name: 'Settings',
                url: window.config.base + 'settings',
                authRequired: true
            },
        ]

        return pages;
    }

    generateLinks = () => {
        var pages = this.generatePagesArray();

        return pages.map(page => {
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
            <BootstrapNavbar variant="dark" bg="dark" expand="sm">
                <BootstrapNavbar.Brand as={Link} to={brand.url}><img style={{width: '15%'}} src={window.config.base + 'files/icons/fav/android-icon-192x192.png'} /> {brand.name}</BootstrapNavbar.Brand>
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
