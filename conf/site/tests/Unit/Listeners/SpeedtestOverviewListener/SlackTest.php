<?php

namespace Tests\Unit\Listeners\SpeedtestOverviewListener;

use App\Helpers\SettingsHelper;
use App\Listeners\SpeedtestOverviewListener;
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
    public function testSlackOverviewNotification()
    {
        SettingsHelper::set('speedtest_overview_notification', true);
        SettingsHelper::set('slack_webhook', env('SLACK_WEBHOOK'));

        $l = new SpeedtestOverviewListener();

        $event = new stdClass();

        try {
            $l->handle($event);
        } catch(Exception $e) {
            $this->assertTrue(false);
            return false;
        }

        $this->assertTrue(true);
    }
}
