<?php

use App\Helpers\SettingsHelper;
use App\Models\Setting;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class AddSpeedtestProviderSetting extends Migration
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
                'name' => 'speedtest_provider',
                'value' => 'ookla',
                'description' => 'The provider/package used to run speedtests.'
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
            'speedtest_provider',
        ])->delete();
    }
}
