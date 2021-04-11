<?php

namespace Tests\Unit\Helpers\SpeedtestHelper;

use App\Utils\OoklaTester;
use PHPUnit\Framework\TestCase;
use Tests\Mocks\OoklaTesterMocker;

class CheckOutputTest extends TestCase
{
    private OoklaTester $speedtestProvider;

    public function setUp(): void
    {
        $this->speedtestProvider = new OoklaTester();
        $this->mocker = new OoklaTesterMocker();
    }

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testGoodOutput()
    {
        $expected = [
            'type' => 'result',
            'download' => ['bandwidth' => '*'],
            'upload' => ['bandwidth' => '*'],
            'ping' => ['latency' => '*'],
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

        $this->assertTrue($this->speedtestProvider->isOutputComplete($expected));
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
            'download' => ['bandwidth' => '*'],
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

        $this->assertFalse($this->speedtestProvider->isOutputComplete($expected));
    }
}
