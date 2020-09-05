<?php

use App\Helpers\SettingsHelper;
use App\Setting;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class AddConditionalNotificationsSettings extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if(!SettingsHelper::get('threshold_alert_percentage_notifications')) {
            Setting::create([
                'name' => 'threshold_alert_percentage_notifications',
                'value' => false,
                'description' => 'Enable/disable theshold percentage notifications'
            ]);
        }

        if(!SettingsHelper::get('threshold_alert_percentage')) {
            Setting::create([
                'name' => 'threshold_alert_percentage',
                'value' => 15,
                'description' => 'When any value of a speedtest is x percent lower than the average, a notification will be sent.'
            ]);
        }

        if(!SettingsHelper::get('threshold_alert_absolute_notifications')) {
            Setting::create([
                'name' => 'threshold_alert_absolute_notifications',
                'value' => false,
                'description' => 'Enable/disable absolute theshold notifications'
            ]);
        }

        if(!SettingsHelper::get('threshold_alert_absolute_download')) {
            Setting::create([
                'name' => 'threshold_alert_absolute_download',
                'value' => '',
                'description' => 'When the download is lower than this value, a notification will be sent. Leave blank to disable'
            ]);
        }

        if(!SettingsHelper::get('threshold_alert_absolute_upload')) {
            Setting::create([
                'name' => 'threshold_alert_absolute_upload',
                'value' => '',
                'description' => 'When the upload is lower than this value, a notification will be sent. Leave blank to disable'
            ]);
        }

        if(!SettingsHelper::get('threshold_alert_absolute_ping')) {
            Setting::create([
                'name' => 'threshold_alert_absolute_ping',
                'value' => '',
                'description' => 'When the ping is higher than this value, a notification will be sent. Leave blank to disable'
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
            'threshold_alert_percentage',
            'threshold_alert_absolute_download',
            'threshold_alert_absolute_upload',
            'threshold_alert_absolute_ping',
            'threshold_alert_percentage_notifications',
            'threshold_alert_absolute_notifications'
        ])->delete();
    }
}
