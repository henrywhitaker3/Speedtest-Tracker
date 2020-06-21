<?php

use App\Setting;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class AddNotificationsSettings extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Setting::create([
            'name' => 'speedtest_notifications',
            'value' => true,
            'description' => 'Enable notifications for every speedtest that runs'
        ]);

        Setting::create([
            'name' => 'speedtest_overview_notification',
            'value' => true,
            'description' => 'Enable a daily notification with average values for the last 24 hours.'
        ]);

        Setting::create([
            'name' => 'speedtest_overview_time',
            'value' => '12',
            'description' => 'The hour (24-hour format) that the daily overview notification will be sent.'
        ]);
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Setting::whereIn('name', [
            'speedtest_notifications',
            'speedtest_overview_notification',
            'speedtest_overview_time',
        ])->delete();
    }
}
