<?php

namespace Tests\Unit\Controllers\SpeedtestController;

use App\Http\Controllers\SpeedtestController;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class IndexTest extends TestCase
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

    public function testIndex()
    {
        $resp = $this->controller->index()->original;

        $this->assertArrayHasKey('method', $resp);
        $this->assertArrayHasKey('data', $resp);
    }
}
