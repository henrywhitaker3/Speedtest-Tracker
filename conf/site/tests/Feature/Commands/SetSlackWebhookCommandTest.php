<?php

namespace Tests\Feature\Commands;

use App\Helpers\SettingsHelper;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Tests\TestCase;

class SetSlackWebhookCommandTest extends TestCase
{
    use RefreshDatabase;

    /**
     * A basic feature test example.
     *
     * @return void
     */
    public function testSetSlackWebhook()
    {
        SettingsHelper::set('slack_webhook', 'pre-test');

        $this->artisan('speedtest:slack', [ 'webhook' => 'test' ])
             ->expectsOutput('Slack webhook updated')
             ->assertExitCode(0);

        $this->assertEquals(SettingsHelper::get('slack_webhook')->value, 'test');
    }
}
