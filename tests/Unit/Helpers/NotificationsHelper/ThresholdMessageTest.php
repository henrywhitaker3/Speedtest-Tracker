<?php

namespace Tests\Unit\Helpers\NotificationsHelper;

use App\Helpers\NotificationsHelper;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Tests\TestCase;

class ThresholdMessageTest extends TestCase
{
    /**
     * Test absolute message
     *
     * @return void
     */
    public function testAbsoluteMessageMultiField()
    {
        $msg = NotificationsHelper::formatAbsoluteThresholdMessage([ 'ping', 'upload' ]);

        $this->assertEquals('For the latest speedtest, the ping, upload values exceeded the absolute threshold', $msg);
    }

    /**
     * Test absolute message
     *
     * @return void
     */
    public function testAbsoluteMessageSingleField()
    {
        $msg = NotificationsHelper::formatAbsoluteThresholdMessage([ 'ping' ]);

        $this->assertEquals('For the latest speedtest, the ping value exceeded the absolute threshold', $msg);
    }

    /**
     * Test absolute message
     *
     * @return void
     */
    public function testPercentageMessageMultiField()
    {
        $msg = NotificationsHelper::formatPercentageThresholdMessage([ 'ping', 'upload' ]);

        $this->assertEquals('For the latest speedtest, the ping, upload values exceeded the percentage threshold', $msg);
    }

    /**
     * Test absolute message
     *
     * @return void
     */
    public function testPercentageMessageSingleField()
    {
        $msg = NotificationsHelper::formatPercentageThresholdMessage([ 'ping' ]);

        $this->assertEquals('For the latest speedtest, the ping value exceeded the percentage threshold', $msg);
    }
}
