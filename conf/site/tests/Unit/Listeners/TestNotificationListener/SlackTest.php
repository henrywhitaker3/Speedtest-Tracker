<?php

namespace Tests\Unit\Listeners\TestNotificationListener;

use App\Events\TestNotificationEvent;
use App\Helpers\SettingsHelper;
use App\Listeners\TestNotificationListener;
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
    public function testSlackTest()
    {
        SettingsHelper::set('slack_webhook', env('SLACK_WEBHOOK'));

        $l = new TestNotificationListener();

        $event = new stdClass();
        $event->agents = [ 'slack' ];

        try {
            $l->handle($event);
        } catch(Exception $e) {
            $this->assertTrue(false);
            return false;
        }

        $this->assertTrue(true);
    }
}
