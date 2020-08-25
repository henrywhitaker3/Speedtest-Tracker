<?php

namespace Tests\Unit\Listeners\SpeedtestCompleteListener;

use App\Helpers\SettingsHelper;
use App\Listeners\SpeedtestCompleteListener;
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
    public function testTelegramSpeedtestNotification()
    {
        SettingsHelper::set('speedtest_notifications', true);
        SettingsHelper::set('telegram_bot_token', env('TELEGRAM_BOT_TOKEN'));
        SettingsHelper::set('telegram_chat_id', env('TELEGRAM_CHAT_ID'));
        SettingsHelper::set('slack_webhook', false);

        $l = new SpeedtestCompleteListener();
        $test = Speedtest::create([ 'download' => 5, 'upload' => 5, 'ping' => 5 ]);

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

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testTelegramPercentageThresholdNotification()
    {
        SettingsHelper::set('speedtest_notifications', false);
        SettingsHelper::set('telegram_bot_token', env('TELEGRAM_BOT_TOKEN'));
        SettingsHelper::set('telegram_chat_id', env('TELEGRAM_CHAT_ID'));
        SettingsHelper::set('slack_webhook', false);
        SettingsHelper::set('threshold_alert_percentage_notifications', 1);
        Speedtest::create([ 'download' => 50, 'upload' => 50, 'ping' => 1 ]);

        $l = new SpeedtestCompleteListener();
        $test = Speedtest::create([ 'download' => 5, 'upload' => 5, 'ping' => 5 ]);

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

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testTelegramAbsoluteThresholdNotification()
    {
        SettingsHelper::set('speedtest_notifications', false);
        SettingsHelper::set('telegram_bot_token', env('TELEGRAM_BOT_TOKEN'));
        SettingsHelper::set('telegram_chat_id', env('TELEGRAM_CHAT_ID'));
        SettingsHelper::set('slack_webhook', false);
        SettingsHelper::set('threshold_alert_absolute_notifications', 1);
        SettingsHelper::set('threshold_alert_absolute_download', 50);
        SettingsHelper::set('threshold_alert_absolute_upload', 50);
        SettingsHelper::set('threshold_alert_absolute_ping', 1);


        $l = new SpeedtestCompleteListener();
        $test = Speedtest::create([ 'download' => 5, 'upload' => 5, 'ping' => 5 ]);

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
