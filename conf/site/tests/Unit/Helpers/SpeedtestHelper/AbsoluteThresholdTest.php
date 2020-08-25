<?php

namespace Tests\Unit\Helpers\SpeedtestHelper;

use App\Helpers\SettingsHelper;
use App\Helpers\SpeedtestHelper;
use App\Speedtest;
use Illuminate\Foundation\Testing\RefreshDatabase;
use InvalidArgumentException;
use Tests\TestCase;

class AbsoluteThresholdTest extends TestCase
{
    use RefreshDatabase;

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testAbsoluteDownloadThresholdExceeded()
    {
        $threshold = 10;
        $dl = 5;

        SettingsHelper::set('threshold_alert_absolute_download', $threshold);

        $test = Speedtest::create([
            'download' => $dl,
            'upload' => 11,
            'ping' => 5
        ]);

        $result = SpeedtestHelper::testIsLowerThanThreshold('absolute', $test);

        $this->assertEquals([ 'download' ], $result);
    }

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testAbsoluteDownloadThresholdNotExceeded()
    {
        $threshold = 1;
        $dl = 5;

        SettingsHelper::set('threshold_alert_absolute_download', $threshold);

        $test = Speedtest::create([
            'download' => $dl,
            'upload' => 11,
            'ping' => 5
        ]);

        $result = SpeedtestHelper::testIsLowerThanThreshold('absolute', $test);

        $this->assertEquals([], $result);
    }

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testAbsoluteUploadThresholdExceeded()
    {
        $threshold = 10;
        $ul = 5;

        SettingsHelper::set('threshold_alert_absolute_upload', $threshold);

        $test = Speedtest::create([
            'download' => 11,
            'upload' => $ul,
            'ping' => 5
        ]);

        $result = SpeedtestHelper::testIsLowerThanThreshold('absolute', $test);

        $this->assertEquals([ 'upload' ], $result);
    }

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testAbsoluteUploadThresholdNotExceeded()
    {
        $threshold = 1;
        $ul = 5;

        SettingsHelper::set('threshold_alert_absolute_upload', $threshold);

        $test = Speedtest::create([
            'download' => 11,
            'upload' => $ul,
            'ping' => 5
        ]);

        $result = SpeedtestHelper::testIsLowerThanThreshold('absolute', $test);

        $this->assertEquals([], $result);
    }

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testAbsolutePingThresholdExceeded()
    {
        $threshold = 10;
        $ping = 11;

        SettingsHelper::set('threshold_alert_absolute_ping', $threshold);

        $test = Speedtest::create([
            'download' => 10,
            'upload' => 10,
            'ping' => $ping
        ]);

        $result = SpeedtestHelper::testIsLowerThanThreshold('absolute', $test);

        $this->assertEquals([ 'ping' ], $result);
    }

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testAbsolutePingThresholdNotExceeded()
    {
        $threshold = 10;
        $ping = 9;

        SettingsHelper::set('threshold_alert_absolute_ping', $threshold);

        $test = Speedtest::create([
            'download' => 10,
            'upload' => 10,
            'ping' => $ping
        ]);

        $result = SpeedtestHelper::testIsLowerThanThreshold('absolute', $test);

        $this->assertEquals([], $result);
    }

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testInvalidArgument()
    {
        try {
            SpeedtestHelper::testIsLowerThanThreshold('test', new Speedtest());
        } catch(InvalidArgumentException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }
}
