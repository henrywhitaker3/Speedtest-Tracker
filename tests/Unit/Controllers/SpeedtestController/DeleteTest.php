<?php

namespace Tests\Unit\Controllers\SpeedtestController;

use App\Http\Controllers\SpeedtestController;
use App\Speedtest;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class DeleteTest extends TestCase
{
    use RefreshDatabase;

    /**
     * SpeedtestController
     *
     * @var SpeedtestController
     */
    private $controller;

    public function setUp() : void
    {
        parent::setUp();

        $this->controller = new SpeedtestController();
    }

    public function testDeleteAll()
    {
        for($i = 0; $i < 5; $i++) {
            Speedtest::create([
                'download' => 5,
                'upload' => 5,
                'ping' => 5
            ]);
        }

        $this->assertEquals(5, Speedtest::count());

        $resp = $this->controller->deleteAll()->original;

        $this->assertArrayHasKey('method', $resp);
        $this->assertArrayHasKey('success', $resp);

        $this->assertEquals(0, Speedtest::count());
    }

    public function testDeleteSpecific()
    {
        $test = Speedtest::create([
            'download' => 5,
            'upload' => 5,
            'ping' => 5
        ]);
        $id = $test->id;

        $this->assertNotNull(Speedtest::find($id));

        $resp = $this->controller->delete($test)->original;

        $this->assertArrayHasKey('method', $resp);
        $this->assertArrayHasKey('success', $resp);

        $this->assertNull(Speedtest::find($id));
    }
}
