<?php

namespace Tests\Unit\Helpers\SettingsHelper;

use App\Helpers\SettingsHelper;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Tests\TestCase;

class SettingsLoadIntegrationsConfigTest extends TestCase
{
    use RefreshDatabase;

    /**
     * A basic feature test example.
     *
     * @return void
     */
    public function testLoadIntegrationsConfig()
    {
        $preLoad = [
            "healthchecks_enabled" => false,
            "healthchecks_uuid" => null,
            "slack_webhook" => null,
            "telegram_bot_token" => null,
            "telegram_chat_id" => null,
        ];

        SettingsHelper::set('slack_webhook', 'test');

        SettingsHelper::loadIntegrationConfig();

        $this->assertEquals(config('integrations.slack_webhook'), 'test');
    }
}
