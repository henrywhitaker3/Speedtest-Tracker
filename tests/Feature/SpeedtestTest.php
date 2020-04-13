<?php

namespace Tests\Feature;

use App\Helpers\SpeedtestHelper;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Tests\TestCase;

class SpeedtestTest extends TestCase
{
    use RefreshDatabase;

    /**
     * Runs a speedtest
     *
     * @test
     * @return void
     */
    public function runSpeedtest()
    {
        $data = SpeedtestHelper::runSpeedtest();

        $this->assertArrayHasKey('ping', $data);
        $this->assertArrayHasKey('download', $data);
        $this->assertArrayHasKey('upload', $data);
    }
}
