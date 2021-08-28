import Axios from 'axios';
import React, { Component } from 'react';
import { Nav, Tab, Tabs } from 'react-bootstrap';
import ReactDOM from 'react-dom';
import { toast } from 'react-toastify';
import SettingsInput from './SettingsInput';
import ResetSettings from './tabs/ResetSettings';
import BackupSettings from './tabs/BackupSettings';
import GeneralSettings from './tabs/GeneralSettings';
import GraphsSettings from './tabs/GraphsSettings';
import HealthchecksSettings from './tabs/HealthchecksSettings';
import NotificationsSettings from './tabs/NotificationsSettings';
import Authentication from '../Authentication/Authentication';
import TableSettings from './tabs/TableSettings';
import InfluxDBSettings from './tabs/InfluxDBSettings';

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
            'Tables',
            'Notifications',
            'healthchecks.io',
            'InfluxDB',
            'Reset',
            'Backup/Restore',
        ];

        if(window.config.auth) {
            tabs.push('Authentication');
        }

        return tabs.map((tab) => {
            return <Tab key={tab} eventKey={tab} title={tab} />
        });
    }

    switchTab = (tab) => {
        this.setState({
            tab: tab
        });
    }

    save = (settings, name) => {
        var url = 'api/settings/bulk?token=' + window.token;
        var data = [];

        settings.forEach(e => {
            if(e.type !== 'btn-get') {
                var res = {
                    name: e.obj.name,
                    value: e.obj.value
                };
                data.push(res);
            }
        });

        data = {
            data: data
        };

        Axios.post(url, data)
        .then((resp) => {
            toast.success(name + ' settings updated');
            Axios.get('api/settings/config')
            .then((resp) => {
                window.config = resp.data;
            })
        })
        .catch((err) => {
            if(err.response.status == 422) {
                toast.error('Your input was invalid');
            } else {
                toast.error('Something went wrong')
            }
        })
    }

    generateInputs = (settings, handler) => {
        return settings.map((setting) => {
            return <SettingsInput
                        key={setting.obj.id}
                        name={setting.obj.name}
                        id={setting.obj.id}
                        type={setting.type}
                        value={setting.obj.value}
                        description={setting.obj.description}
                        handler={handler}
                        label={setting.obj.name}
                        description={setting.obj.description}
                        options={setting.type == 'select' ? setting.options : []}
                        hideDescription={setting.hideDescription ? setting.hideDescription : false}
                        min={setting.min ? setting.min : false}
                        max={setting.max ? setting.max : false}
                        btnType={setting.btnType}
                        inline={setting.inline}
                        url={setting.url}
                        earlyReturn={setting.earlyReturn ? true : false}
                        classes={setting.classes ? setting.classes : ''}
                        autoComplete={setting.autoComplete ? true : false}
                    />
        })
    }

    getTabContent = () => {
        var data = this.state.data;

        switch(this.state.tab) {
            case 'General':
                return <GeneralSettings
                            data={data.General}
                            generateInputs={this.generateInputs}
                            save={this.save} />
            case 'Graphs':
                return <GraphsSettings
                            data={data.Graphs}
                            generateInputs={this.generateInputs}
                            save={this.save} />
            case 'Tables':
                return <TableSettings
                            data={data.Tables}
                            refreshConfig={this.props.refreshConfig}
                            save={this.save} />
            case 'Notifications':
                return <NotificationsSettings
                            data={data.Notifications}
                            generateInputs={this.generateInputs}
                            save={this.save} />
            case 'healthchecks.io':
                return <HealthchecksSettings
                            data={data.healthchecks}
                            generateInputs={this.generateInputs}
                            save={this.save} />
            case 'InfluxDB':
                return <InfluxDBSettings
                            data={data.influxdb}
                            generateInputs={this.generateInputs}
                            save={this.save} />
            case 'Reset':
                return <ResetSettings
                            data={data.healthchecks}
                            generateInputs={this.generateInputs}
                            save={this.save} />
            case 'Backup/Restore':
                return <BackupSettings
                            data={data.healthchecks}
                            generateInputs={this.generateInputs}
                            save={this.save} />
            case 'Authentication':
                return <Authentication
                            data={data.healthchecks}
                            generateInputs={this.generateInputs}
                            save={this.save} />
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
