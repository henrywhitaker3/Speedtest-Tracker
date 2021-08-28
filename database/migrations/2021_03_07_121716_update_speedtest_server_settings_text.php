<?php

use App\Models\Setting;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class UpdateSpeedtestServerSettingsText extends Migration
{
    private Setting $setting;

    public function __construct()
    {
        $this->setting = Setting::where('name', 'server')->first();
    }
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        $this->setting->description = '<p class="d-inline">Comma-separated list of speedtest.net server IDs picked randomly. Leave blank to use default settings.</p>;';
        $this->setting->save();
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        $this->setting->description = '<p class="d-inline">Comma-separated list of speedtest.net servers picked randomly. Leave blank to use default settings.</p>';
        $this->setting->save();
    }
}
