<?php

use App\Helpers\SettingsHelper;
use App\Setting;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class AddShowFailedTestsSetting extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if(!SettingsHelper::get('show_failed_tests_on_graph')) {
            Setting::create([
                'name' => 'show_failed_tests_on_graph',
                'value' => true,
                'description' => 'If enabled, failed tests will appear on the graphs as 0.'
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
            'show_failed_tests_on_graph',
        ])->delete();
    }
}
