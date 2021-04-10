<?php

use App\Helpers\SettingsHelper;
use App\Models\Setting;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class AddAppNameSetting extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!SettingsHelper::get('app_name')) {
            Setting::create([
                'name' => 'app_name',
                'value' => 'Speedtest Tracker',
                'description' => 'Set a custom app name'
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
            'app_name',
        ])->delete();
    }
}
