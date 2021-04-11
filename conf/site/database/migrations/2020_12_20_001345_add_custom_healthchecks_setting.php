<?php

use App\Helpers\SettingsHelper;
use App\Models\Setting;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class AddCustomHealthchecksSetting extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!SettingsHelper::get('healthchecks_server_url')) {
            Setting::create([
                'name' => 'healthchecks_server_url',
                'value' => 'https://hc-ping.com/',
                'description' => 'The URL of the healthchecks.io server. Change this to use a self-hosted server.'
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
            'healthchecks_server_url',
        ])->delete();
    }
}
