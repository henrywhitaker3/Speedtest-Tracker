<?php

use App\Helpers\SettingsHelper;
use App\Models\Setting;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class AddHiddenColumnsSetting extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!SettingsHelper::get('hidden_columns')) {
            Setting::create([
                'name' => 'hidden_columns',
                'value' => [
                    'server_id', 'server_name', 'server_host', 'url', 'scheduled',
                ],
                'description' => 'Columns hidden from the "All Tests" table.'
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
            'hidden_columns',
        ])->delete();
    }
}
