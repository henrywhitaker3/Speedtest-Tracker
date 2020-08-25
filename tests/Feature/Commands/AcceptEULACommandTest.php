<?php

namespace Tests\Feature\Commands;

use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Tests\TestCase;

class AcceptEULACommandTest extends TestCase
{
    /**
     * A basic feature test example.
     *
     * @return void
     */
    public function testAcceptEULA()
    {
        $response = $this->artisan('speedtest:eula')
                         ->assertExitCode(0);
    }
}
