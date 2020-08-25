<?php

namespace Tests\Unit\Listeners\SpeedtestFailedListener;

use App\Helpers\SettingsHelper;
use App\Listeners\SpeedtestFailedListener;
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
    public function testSlackFailedNotification()
    {
        SettingsHelper::set('speedtest_notifications', true);
        SettingsHelper::set('slack_webhook', env('SLACK_WEBHOOK'));

        $l = new SpeedtestFailedListener();
        $test = Speedtest::create([ 'download' => 0, 'upload' => 0, 'ping' => 0, 'failed' => true ]);

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
