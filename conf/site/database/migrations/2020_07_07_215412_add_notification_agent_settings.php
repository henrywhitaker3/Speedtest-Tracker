<?php

use App\Helpers\SettingsHelper;
use App\Setting;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class AddNotificationAgentSettings extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if(!SettingsHelper::get('slack_webhook')) {
            Setting::create([
                'name' => 'slack_webhook',
                'value' => '',
                'description' => ''
            ]);
        }

        if(!SettingsHelper::get('telegram_chat_id')) {
            Setting::create([
                'name' => 'telegram_chat_id',
                'value' => '',
                'description' => ''
            ]);
        }

        if(!SettingsHelper::get('telegram_bot_token')) {
            Setting::create([
                'name' => 'telegram_bot_token',
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
            'slack_webhook',
            'telegram_chat_id',
            'telegram_bot_token',
        ])->delete();
    }
}
