<?php

namespace Tests\Unit\Helpers\SpeedtestHelper;

use App\Helpers\SpeedtestHelper;
use App\Speedtest;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class FailureRateTest extends TestCase
{
    use RefreshDatabase;

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testFailureRate()
    {
        $success = rand(1, 15);
        $failed = rand(1, 15);

        for($i = 0; $i < $success; $i++) {
            Speedtest::create([
                'ping' => 5,
                'download' => 5,
                'upload' => 5
            ]);
        }

        for($i = 0; $i < $failed; $i++) {
            Speedtest::create([
                'ping' => 5,
                'download' => 5,
                'upload' => 5,
                'failed' => true
            ]);
        }

        $output = SpeedtestHelper::failureRate(1);

        $this->assertEquals($output[0]['success'], $success);
        $this->assertEquals($output[0]['failure'], $failed);
    }
}
