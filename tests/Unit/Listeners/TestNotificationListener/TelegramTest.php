<?php

namespace Tests\Unit\Listeners\TestNotificationListener;

use App\Helpers\SettingsHelper;
use App\Listeners\TestNotificationListener;
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
    public function testTelegramTest()
    {
        SettingsHelper::set('telegram_bot_token', env('TELEGRAM_BOT_TOKEN'));
        SettingsHelper::set('telegram_chat_id', env('TELEGRAM_CHAT_ID'));
        SettingsHelper::set('slack_webhook', false);

        $l = new TestNotificationListener();

        $event = new stdClass();
        $event->agents = [ 'telegram' ];

        try {
            $l->handle($event);
        } catch(Exception $e) {
            $this->assertTrue(false);
            return false;
        }

        $this->assertTrue(true);
    }
}
