<?php

namespace Henrywhitaker3\Healthchecks\Tests;

use Henrywhitaker3\Healthchecks\Exceptions\HealthchecksUuidNotFoundException;
use Henrywhitaker3\Healthchecks\HealthchecksManager;
use Ramsey\Uuid\Uuid;

class SpecificCheck404Test extends Test
{
    /**
     * Instance of manager
     *
     * @var HealthchecksManager
     */
    private $hm;

    /**
     * Random uuid
     *
     * @var String
     */
    private $uuid;

    public function setUp(): void
    {
        parent::setUp();

        $this->hm = new HealthchecksManager($_ENV['APIKEY']);
        $this->uuid = Uuid::uuid4();
    }

    /**
     * Tests listing specific check from API
     *
     * @return void
     */
    public function testGetCheck404()
    {
        try {
            $this->hm->getCheck($this->uuid);
        } catch(HealthchecksUuidNotFoundException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }

    /**
     * Test pausing a check
     *
     * @return void
     */
    public function testPauseCheck404()
    {
        try {
            $this->hm->pauseCheck($this->uuid);
        } catch(HealthchecksUuidNotFoundException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }

    /**
     * Test deleting a check
     *
     * @return void
     */
    public function testDeleteCheck()
    {
        try {
            $this->hm->deleteCheck($this->uuid);
        } catch(HealthchecksUuidNotFoundException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }

    /**
     * Test getCheckPings
     *
     * @return void
     */
    public function testGetCheckPings()
    {
        try {
            $this->hm->getCheckPings($this->uuid);
        } catch(HealthchecksUuidNotFoundException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }

    /**
     * Test getCheckStatusChanges
     *
     * @return void
     */
    public function testGetCheckStatusChanges()
    {
        try {
            $this->hm->getCheckStatusChanges($this->uuid);
        } catch(HealthchecksUuidNotFoundException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }
}