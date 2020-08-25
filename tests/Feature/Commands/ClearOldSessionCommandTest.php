<?php

namespace Tests\Feature\Commands;

use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Tests\TestCase;

class ClearOldSessionCommandTest extends TestCase
{
    use RefreshDatabase;

    /**
     * A basic feature test example.
     *
     * @return void
     */
    public function testClearSessionsCommand()
    {
        $this->artisan('speedtest:clear-sessions')
             ->expectsOutput('Invalidated expired sessions')
             ->assertExitCode(0);
    }
}
