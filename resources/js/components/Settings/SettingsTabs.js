import React, { Component } from 'react';
import { Nav, Tab, Tabs } from 'react-bootstrap';
import ReactDOM from 'react-dom';
import GeneralSettings from './tabs/GeneralSettings';
import GraphsSettings from './tabs/GraphsSettings';

export default class SettingsTabs extends Component {
    constructor(props) {
        super(props)

        this.state = {
            tab: "General",
            data: this.props.data
        }
    }

    generateTabs = () => {
        var tabs = [
            'General',
            'Graphs',
            'Notifications',
            'healthchecks.io',
            'Reset',
        ];

        return tabs.map((tab) => {
            return <Tab key={tab} eventKey={tab} title={tab} />
        });
    }

    switchTab = (tab) => {
        this.setState({
            tab: tab
        });
    }

    getTabContent = () => {
        var data = this.state.data;
        console.log(data);

        switch(this.state.tab) {
            case 'General':
                return <GeneralSettings data={data.General} />
                break;
            case 'Graphs':
                return <GraphsSettings />
                break;
        }
    }

    render() {
        var tabs = this.generateTabs();
        var activeTab = this.state.tab;
        var tabContent = this.getTabContent();

        return (
            <div>
                <Tabs
                    variant="tabs"
                    onSelect={(tab) => { this.switchTab(tab) }}
                    activeKey={activeTab}
                >
                    {tabs}
                </Tabs>

                <div className="mt-3">
                    {tabContent}
                </div>
            </div>
        );
    }
}

if (document.getElementById('settingsTabs')) {
    ReactDOM.render(<SettingsTabs />, document.getElementById('settingsTabs'));
}
