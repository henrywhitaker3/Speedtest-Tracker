<?php

namespace Tests\Unit\Controllers\SpeedtestController;

use App\Http\Controllers\SpeedtestController;
use App\Speedtest;
use DB;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class LatestTest extends TestCase
{
    use RefreshDatabase;

    /**
     * SpeedtestController
     *
     * @var SpeedtestController
     */
    private $controller;

    public function setUp(): void
    {
        parent::setUp();

        $this->controller = new SpeedtestController();
    }

    public function testLatestNoEntries()
    {
        DB::table('speedtests')->delete();

        $resp = $this->controller->latest();
        $resp = $resp->original;

        $this->assertEquals([
            'method' => 'get latest speedtest',
            'error' => 'no speedtests have been run'
        ], $resp);
    }

    public function testLatest()
    {
        $test = Speedtest::create([
            'download' => 5,
            'upload' => 5,
            'ping' => 5
        ]);
        $test = $test->attributesToArray();

        $resp = $this->controller->latest();
        $resp = $resp->original;

        $this->assertArrayHasKey('data', $resp);
        $this->assertArrayHasKey('average', $resp);
        $this->assertArrayHasKey('maximum', $resp);
        $this->assertArrayHasKey('minimum', $resp);
    }
}
