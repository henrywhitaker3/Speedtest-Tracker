<?php

use App\Helpers\SettingsHelper;
use App\Models\Setting;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class AddScheduleEnabledSetting extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!SettingsHelper::get('schedule_enabled')) {
            Setting::create([
                'name' => 'schedule_enabled',
                'value' => true,
                'description' => 'Enable/disable the schedule worker'
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
            'schedule_enabled',
        ])->delete();
    }
}
