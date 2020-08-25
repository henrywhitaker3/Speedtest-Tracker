<?php

namespace Tests\Unit\Helpers\SettingsHelper;

use App\Helpers\SettingsHelper;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Tests\TestCase;

class SettingsGetTest extends TestCase
{
    use RefreshDatabase;

    /**
     * A basic feature test example.
     *
     * @return void
     */
    public function testValidSetting()
    {
        $response = SettingsHelper::get('auth');

        $this->assertNotFalse($response);
        $this->assertFalse((bool)$response->value);
    }

    /**
     * A basic feature test example.
     *
     * @return void
     */
    public function testInvalidSetting()
    {
        $response = SettingsHelper::get('test');

        $this->assertFalse($response);
    }
}
