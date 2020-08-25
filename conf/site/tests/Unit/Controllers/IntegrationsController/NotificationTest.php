<?php

namespace Tests\Unit\Controllers\IntegrationsController;

use App\Http\Controllers\IntegrationsController;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class NotificationTest extends TestCase
{
    use RefreshDatabase;

    /**
     * Controller
     *
     * @var IntegrationsController
     */
    private $controller;

    public function setUp() : void
    {
        parent::setUp();

        $this->controller = new IntegrationsController();
    }

    public function testNotificationsTest()
    {
        $resp = $this->controller->testNotification()->original;

        $this->assertEquals([
            'method' => 'test notification agents'
        ], $resp);
    }
}
