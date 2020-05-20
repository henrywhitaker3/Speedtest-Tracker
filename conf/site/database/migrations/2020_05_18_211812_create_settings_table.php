<?php

use App\Helpers\SettingsHelper;
use App\Setting;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class CreateSettingsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('settings', function (Blueprint $table) {
            $table->id();
            $table->string('name')->unique();
            $table->string('description')->nullable();
            $table->string('value');
            $table->timestamps();
        });

        $settings = [
            [
                'name' => 'schedule',
                'value' => '0 * * * *',
                'description' => '<p class="d-inline">Set the schedule for speedtests to run using the CRON format. </p><a href="https://crontab.guru/" target="_blank" rel="noopener noreferer">This site</a> can help with formatting.'
            ],
            [
                'name' => 'server',
                'value' => '',
                'description' => '<p class="d-inline">Comma-separated list of speedtest.net servers picked randomly. Leave blank to use default settings.</p>'
            ]
        ];

        foreach($settings as $s) {
            Setting::create([
                'name' => $s['name'],
                'value' => $s['value'],
                'description' => $s['description']
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
        Schema::dropIfExists('settings');
    }
}
