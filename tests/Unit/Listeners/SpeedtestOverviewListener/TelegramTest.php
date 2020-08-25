<?php

namespace Tests\Unit\Listeners\SpeedtestOverviewListener;

use App\Helpers\SettingsHelper;
use App\Listeners\SpeedtestOverviewListener;
use App\Speedtest;
use Exception;
use Illuminate\Foundation\Testing\RefreshDatabase;
use stdClass;
use Tests\TestCase;

class TelegramTest extends TestCase
{
    use RefreshDatabase;

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testTelegramOverviewNotification()
    {
        SettingsHelper::set('speedtest_overview_notification', true);
        SettingsHelper::set('telegram_bot_token', env('TELEGRAM_BOT_TOKEN'));
        SettingsHelper::set('telegram_chat_id', env('TELEGRAM_CHAT_ID'));
        SettingsHelper::set('slack_webhook', false);

        $l = new SpeedtestOverviewListener();
        $test = Speedtest::create([ 'download' => 5, 'upload' => 5, 'ping' => 5, 'failed' => true ]);

        $event = new stdClass();
        $event->speedtest = $test;

        try {
            $l->handle($event);
        } catch(Exception $e) {
            $this->assertTrue(false);
            return false;
        }

        $this->assertTrue(true);
    }
}
