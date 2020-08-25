<?php

namespace Tests\Unit\Helpers\SettingsHelper;

use App\Helpers\SettingsHelper;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Tests\TestCase;

class SettingsSetTest extends TestCase
{
    use RefreshDatabase;

    /**
     * A basic feature test example.
     *
     * @return void
     */
    public function testUpdatingSetting()
    {
        $response = SettingsHelper::set('auth', 'hello');

        $this->assertEquals('hello', $response->value);
    }
}
