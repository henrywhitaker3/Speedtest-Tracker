<?php

namespace Tests\Unit\Helpers\SpeedtestHelper;

use App\Helpers\SpeedtestHelper;
use PHPUnit\Framework\TestCase;

class ParseUnitsTest extends TestCase
{
    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testMbyteToMb()
    {
        $initial = '80 Mbyte';
        $expected = 640;

        $result = SpeedtestHelper::parseUnits($initial);

        $this->assertEquals($expected, $result['val']);
    }

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testKbitToMb()
    {
        $initial = '80 Kbit';
        $expected = 0.08;

        $result = SpeedtestHelper::parseUnits($initial);

        $this->assertEquals($expected, $result['val']);
    }

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testKbyteToMb()
    {
        $initial = '80 Kbyte';
        $expected = 0.64;

        $result = SpeedtestHelper::parseUnits($initial);

        $this->assertEquals($expected, $result['val']);
    }

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testMbToMb()
    {
        $initial = '80 Mbit';
        $expected = 80;

        $result = SpeedtestHelper::parseUnits($initial);

        $this->assertEquals($expected, $result['val']);
    }
}
