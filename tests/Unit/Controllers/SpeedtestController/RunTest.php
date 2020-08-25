<?php

namespace Tests\Unit\Controllers\SpeedtestController;

use App\Http\Controllers\SpeedtestController;
use App\Jobs\SpeedtestJob;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Queue;
use Tests\TestCase;

class RunTest extends TestCase
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

    public function testRun()
    {
        Queue::fake();

        $resp = $this->controller->run()->original;

        $this->assertArrayHasKey('method', $resp);
        $this->assertArrayHasKey('data', $resp);

        Queue::assertPushed(SpeedtestJob::class, function($job) {
            return true;
        });
    }
}
