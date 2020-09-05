<?php

namespace Tests\Unit\Listeners\SpeedtestFailedListener;

use App\Helpers\SettingsHelper;
use App\Listeners\SpeedtestFailedListener;
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
    public function testTelegramFailedNotification()
    {
        SettingsHelper::set('speedtest_notifications', true);
        SettingsHelper::set('telegram_bot_token', env('TELEGRAM_BOT_TOKEN'));
        SettingsHelper::set('telegram_chat_id', env('TELEGRAM_CHAT_ID'));
        SettingsHelper::set('slack_webhook', false);

        $l = new SpeedtestFailedListener();
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
