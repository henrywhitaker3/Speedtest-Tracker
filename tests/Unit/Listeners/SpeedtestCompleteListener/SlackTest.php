<?php

namespace Tests\Unit\Listeners\SpeedtestCompleteListener;

use App\Helpers\SettingsHelper;
use App\Listeners\SpeedtestCompleteListener;
use App\Speedtest;
use Exception;
use Illuminate\Foundation\Testing\RefreshDatabase;
use stdClass;
use Tests\TestCase;

class SlackTest extends TestCase
{
    use RefreshDatabase;

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testSlackSpeedtestNotification()
    {
        SettingsHelper::set('speedtest_notifications', true);
        SettingsHelper::set('slack_webhook', env('SLACK_WEBHOOK'));

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
    public function testInvalidSlackWebhook()
    {
        SettingsHelper::set('speedtest_notifications', true);
        SettingsHelper::set('slack_webhook', 'invalid');

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
    public function testSlackPercentageThresholdNotification()
    {
        SettingsHelper::set('speedtest_notifications', false);
        SettingsHelper::set('slack_webhook', env('SLACK_WEBHOOK'));
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
    public function testSlackAbsoluteThresholdNotification()
    {
        SettingsHelper::set('speedtest_notifications', false);
        SettingsHelper::set('slack_webhook', env('SLACK_WEBHOOK'));
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
