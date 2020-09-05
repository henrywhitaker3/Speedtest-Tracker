<?php

namespace Tests\Unit\Helpers\SpeedtestHelper;

use App\Helpers\SettingsHelper;
use App\Helpers\SpeedtestHelper;
use App\Speedtest;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class PercentageThresholdTest extends TestCase
{
    use RefreshDatabase;

    public function setUp() : void
    {
        parent::setUp();

        $this->createAverageOf5();
    }

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testPercentageThresholdDownloadExceeded()
    {
        SettingsHelper::set('threshold_alert_percentage', 15);

        $test = Speedtest::create([
            'download' => 5 * 0.8,
            'upload' => 5,
            'ping' => 5
        ]);

        $result = SpeedtestHelper::testIsLowerThanThreshold('percentage', $test);

        $this->assertEquals([ 'download' ], $result);
    }

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testThresholdNotSet()
    {
        SettingsHelper::set('threshold_alert_percentage', '');

        $test = Speedtest::create([
            'download' => 5 * 0.9,
            'upload' => 5,
            'ping' => 5
        ]);

        $result = SpeedtestHelper::testIsLowerThanThreshold('percentage', $test);

        $this->assertEquals([], $result);
    }

    private function createAverageOf5()
    {
        for($i = 0; $i < 5; $i++) {
            Speedtest::create([
                'download' => 5,
                'upload' => 5,
                'ping' => 5
            ]);
        }
    }
}
