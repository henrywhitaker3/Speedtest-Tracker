<?php

use App\Helpers\SettingsHelper;
use App\Models\Setting;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class AddVisibleColumnsSetting extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!SettingsHelper::get('visible_columns')) {
            Setting::create([
                'name' => 'visible_columns',
                'value' => [
                    'id', 'created_at', 'download', 'upload', 'ping'
                ],
                'description' => 'Choose and order the columns shown in the "All Tests" table.'
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
            'visible_columns',
        ])->delete();
    }
}
