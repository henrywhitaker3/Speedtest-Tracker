<?php

namespace Tests\Unit\Helpers\SpeedtestHelper;

use App\Helpers\SpeedtestHelper;
use PHPUnit\Framework\TestCase;

class CheckOutputTest extends TestCase
{
    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testGoodOutput()
    {
        $expected = [
            'type' => 'result',
            'download' => [ 'bandwidth' => '*' ],
            'upload' => [ 'bandwidth' => '*' ],
            'ping' => [ 'latency' => '*' ],
            'server' => [
                'id' => '*',
                'name' => '*',
                'host' => '*',
                'port' => '*',
            ],
            'result' => [
                'url' => '*',
            ]
        ];

        $this->assertTrue(SpeedtestHelper::checkOutputIsComplete($expected));
    }

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testBadOutput()
    {
        $expected = [
            'type' => 'result',
            'download' => [ 'bandwidth' => '*' ],
            'server' => [
                'id' => '*',
                'name' => '*',
                'host' => '*',
                'port' => '*',
            ],
            'result' => [
                'url' => '*',
            ]
        ];

        $this->assertFalse(SpeedtestHelper::checkOutputIsComplete($expected));
    }
}
