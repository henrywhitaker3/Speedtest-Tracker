<?php

namespace Tests\Feature;

use App\Helpers\SettingsHelper;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Tests\TestCase;

class SetSlackWebhookTest extends TestCase
{
    use RefreshDatabase;

    /**
     * Test settings slack webhook via API
     *
     * @return void
     */
    public function testSetSlackWebhookAPI()
    {
        $response = $this->json('PUT', 'api/settings', [
            'name' => 'slack_webhook',
            'value' => 'PHPUnitAPI'
        ]);

        $response->assertStatus(200);
        $this->assertEquals('PHPUnitAPI', SettingsHelper::get('slack_webhook')->value);
    }
}
