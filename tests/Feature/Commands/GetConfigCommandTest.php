<?php

namespace Tests\Feature\Commands;

use App\Helpers\SettingsHelper;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Tests\TestCase;

class GetConfigCommandTest extends TestCase
{
    use RefreshDatabase;

    /**
     * A basic feature test example.
     *
     * @return void
     */
    public function testGetConfig()
    {
        $configJson = json_encode(SettingsHelper::getConfig(), JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES);

        $this->artisan('speedtest:config')
             ->expectsOutput($configJson)
             ->assertExitCode(0);
    }
}
