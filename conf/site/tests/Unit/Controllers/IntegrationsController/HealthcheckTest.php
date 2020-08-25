<?php

namespace Tests\Unit\Controllers\IntegrationsController;

use App;
use App\Helpers\SettingsHelper;
use App\Http\Controllers\IntegrationsController;
use Henrywhitaker3\Healthchecks\Healthchecks;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class HealthcheckTest extends TestCase
{
    use RefreshDatabase;

    /**
     * Controller
     *
     * @var IntegrationsController
     */
    private $controller;

    private $uuid;

    public function setUp() : void
    {
        parent::setUp();

        $this->controller = new IntegrationsController();

        $this->uuid = env('HEALTHCHECKS_UUID');
        $this->bindFacade($this->uuid);
    }

    public function testStartPing()
    {
        $resp = $this->controller->testHealthchecks('start')->original;

        $this->assertEquals([
            'method' => 'test healthchecks \'start\' endpoint',
            'success' => true
        ], $resp);
    }

    public function testSuccessPing()
    {
        $resp = $this->controller->testHealthchecks('success')->original;

        $this->assertEquals([
            'method' => 'test healthchecks \'success\' endpoint',
            'success' => true
        ], $resp);
    }

    public function testFailPing()
    {
        $resp = $this->controller->testHealthchecks('fail')->original;

        $this->assertEquals([
            'method' => 'test healthchecks \'fail\' endpoint',
            'success' => true
        ], $resp);
    }

    public function testInvalidUUID()
    {
        $this->bindFacade('test');

        $resp = $this->controller->testHealthchecks('start')->original;

        $this->assertEquals([
            'method' => 'test healthchecks \'start\' endpoint',
            'success' => false,
            'error' => 'Invalid UUID'
        ], $resp);

        $this->bindFacade($this->uuid);
    }

    /**
     * As clean install before setting up, there is no healthchecks
     * uuid in the db, so the facade doesn't get created during boot,
     * now just bind it in the container on test setup
     *
     * @param String $uuid
     * @return void
     */
    private function bindFacade(String $uuid)
    {
        App::bind('healthcheck', function() use ($uuid) {
            return new Healthchecks($uuid);
        });
    }
}
