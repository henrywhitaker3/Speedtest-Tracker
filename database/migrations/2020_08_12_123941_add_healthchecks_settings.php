<?php

use App\Helpers\SettingsHelper;
use App\Setting;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class AddHealthchecksSettings extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if(!SettingsHelper::get('healthchecks_enabled')) {
            Setting::create([
                'name' => 'healthchecks_enabled',
                'value' => false,
                'description' => 'Enable the healthchecks.io integration for speedtests.'
            ]);
        }

        if(!SettingsHelper::get('healthchecks_uuid')) {
            Setting::create([
                'name' => 'healthchecks_uuid',
                'value' => '',
                'description' => ''
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
            'healthchecks_enabled',
            'healthchecks_uuid',
        ])->delete();
    }
}
