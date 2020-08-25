<?php

namespace Tests\Feature\Commands;

use App\Helpers\SettingsHelper;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Tests\TestCase;

class SetTelegramOptionsCommandTest extends TestCase
{
    use RefreshDatabase;

    /**
     * A basic feature test example.
     *
     * @return void
     */
    public function testSetTelegramOptions()
    {
        SettingsHelper::set('telegram_bot_token', 'pre-test-bot');
        SettingsHelper::set('telegram_chat_id', 'pre-test-bot');

        $this->artisan('speedtest:telegram', [
            '--bot' => 'test-bot',
            '--chat' => 'test-chat'
        ])->expectsOutput('Telegram options updated')
          ->assertExitCode(0);

        $this->assertEquals(SettingsHelper::get('telegram_bot_token')->value, 'test-bot');
        $this->assertEquals(SettingsHelper::get('telegram_chat_id')->value, 'test-chat');
    }
}
