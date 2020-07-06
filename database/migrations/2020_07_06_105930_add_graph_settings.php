<?php

use App\Setting;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class AddGraphSettings extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Setting::create([
            'name' => 'download_upload_graph_enabled',
            'value' => true,
            'description' => 'Enable the download/upload graph'
        ]);

        Setting::create([
            'name' => 'download_upload_graph_width',
            'value' => 6,
            'description' => 'Set the width of the download/upload graph'
        ]);

        Setting::create([
            'name' => 'ping_graph_enabled',
            'value' => true,
            'description' => 'Enable the ping graph'
        ]);

        Setting::create([
            'name' => 'ping_graph_width',
            'value' => 6,
            'description' => 'Set the width of the ping graph'
        ]);

        Setting::create([
            'name' => 'failure_graph_enabled',
            'value' => true,
            'description' => 'Enable the failure rate graph'
        ]);

        Setting::create([
            'name' => 'failure_graph_width',
            'value' => 6,
            'description' => 'Set the width of the failure rate graph'
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
            'download_upload_graph_enabled',
            'download_upload_graph_width',
            'ping_graph_enabled',
            'ping_graph_width',
            'failure_graph_enabled',
            'failure_graph_width'
        ])->delete();
    }
}
