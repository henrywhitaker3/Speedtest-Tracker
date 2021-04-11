<?php

use App\Helpers\SettingsHelper;
use App\Models\Setting;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class AddInfluxDbSettings extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!SettingsHelper::get('influx_db_enabled')) {
            Setting::create([
                'name' => 'influx_db_enabled',
                'value' => false,
                'description' => 'Enable the InfluxDB integration for speedtests.'
            ]);
        }

        if (!SettingsHelper::get('influx_db_host')) {
            Setting::create([
                'name' => 'influx_db_host',
                'value' => '',
                'description' => 'InfluxDB hostname, include the protocol (http:// or https://).'
            ]);
        }

        if (!SettingsHelper::get('influx_db_port')) {
            Setting::create([
                'name' => 'influx_db_port',
                'value' => '',
                'description' => 'InfluxDB port'
            ]);
        }

        if (!SettingsHelper::get('influx_db_database')) {
            Setting::create([
                'name' => 'influx_db_database',
                'value' => '',
                'description' => 'InfluxDB database'
            ]);
        }

        if (!SettingsHelper::get('influx_db_username')) {
            Setting::create([
                'name' => 'influx_db_username',
                'value' => '',
                'description' => 'InfluxDB username'
            ]);
        }

        if (!SettingsHelper::get('influx_db_password')) {
            Setting::create([
                'name' => 'influx_db_password',
                'value' => '',
                'description' => 'InfluxDB password'
            ]);
        }

        if (!SettingsHelper::get('influx_db_version')) {
            Setting::create([
                'name' => 'influx_db_version',
                'value' => 1,
                'description' => 'InfluxDB version'
            ]);
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Setting::whereIn('name', [
            'influx_db_enabled',
            'influx_db_host',
            'influx_db_port',
            'influx_db_database',
            'influx_db_username',
            'influx_db_password',
            'influx_db_version',
        ])->delete();
    }
}
