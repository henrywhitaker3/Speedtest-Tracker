<?php

use App\Helpers\SettingsHelper;
use App\Setting;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class AddWidgetCardSettings extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!SettingsHelper::get('show_average')) {
            Setting::create([
                'name' => 'show_average',
                'value' => true,
                'description' => 'If enabled, the average value for speedtests will be shown in the widgets.'
            ]);
        }

        if (!SettingsHelper::get('show_max')) {
            Setting::create([
                'name' => 'show_max',
                'value' => true,
                'description' => 'If enabled, the maximum value for speedtests will be shown in the widgets.'
            ]);
        }

        if (!SettingsHelper::get('show_min')) {
            Setting::create([
                'name' => 'show_min',
                'value' => true,
                'description' => 'If enabled, the minimum value for speedtests will be shown in the widgets.'
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
            'show_average',
            'show_max',
            'show_min',
        ])->delete();
    }
}
