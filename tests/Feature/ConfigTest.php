<?php

namespace Tests\Feature;

use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Tests\TestCase;

class ConfigTest extends TestCase
{
    use RefreshDatabase;

    private $configStructure = [
        'base',
        'graphs' => [
            'download_upload_graph_enabled' => [],
            'download_upload_graph_width' => [],
            'ping_graph_enabled' => [],
            'ping_graph_width' => [],
            'failure_graph_enabled' => [],
            'failure_graph_width' => [],
        ],
        'editable' => [
            'slack_webhook',
            'telegram_bot_token',
            'telegram_chat_id'
        ],
    ];

    /**
     * Test config returned by API
     *
     * @return void
     */
    public function testAPIConfig()
    {
        $response = $this->get('api/settings/config');

        $response->assertJsonStructure($this->configStructure);
    }
}
