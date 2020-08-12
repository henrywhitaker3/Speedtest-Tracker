<?php

namespace Henrywhitaker3\Healthchecks\Tests;

use Henrywhitaker3\Healthchecks\Exceptions\HealthchecksUnauthorisedException;
use Henrywhitaker3\Healthchecks\Exceptions\HealthchecksUuidNotFoundException;
use Henrywhitaker3\Healthchecks\HealthchecksManager;
use Ramsey\Uuid\Uuid;

class SpecificCheckUnauthTest extends Test
{
    /**
     * Instance of manager
     *
     * @var HealthchecksManager
     */
    private $hmAuth;

    /**
     * Instance of manager
     *
     * @var HealthchecksManager
     */
    private $hmUnauth;

    /**
     * UUID of created check
     *
     * @var String
     */
    private $uuid;

    public function setUp(): void
    {
        parent::setUp();

        $this->hmAuth = new HealthchecksManager($_ENV['APIKEY']);
        $this->uuid = substr($this->hmAuth->createCheck()['ping_url'], 20);
        $this->hmUnauth = new HealthchecksManager('abc');
    }

    public function tearDown(): void
    {
        $this->hmAuth->deleteCheck($this->uuid);
    }

    /**
     * Tests getting specific check from API
     *
     * @return void
     */
    public function testGetCheckUnauth()
    {
        try {
            $this->hmUnauth->getCheck($this->uuid);
        } catch(HealthchecksUnauthorisedException $e) {
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
    public function testPauseCheckUnauth()
    {
        try {
            $this->hmUnauth->pauseCheck($this->uuid);
        } catch(HealthchecksUnauthorisedException $e) {
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
    public function testDeleteCheckUnauth()
    {
        $uuid = substr($this->hmAuth->createCheck()['ping_url'], 20);

        try {
            $this->hmUnauth->deleteCheck($uuid);
        } catch(HealthchecksUnauthorisedException $e) {
            $this->assertTrue(true);
        }

        $this->hmAuth->deleteCheck($uuid);
    }

    /**
     * Test getCheckPings
     *
     * @return void
     */
    public function testGetCheckPingsUnauth()
    {
        try {
            $this->hmUnauth->getCheckPings($this->uuid);
        } catch(HealthchecksUnauthorisedException $e) {
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
    public function testGetCheckStatusChangesUnauth()
    {
        try {
            $this->hmUnauth->getCheckStatusChanges($this->uuid);
        } catch(HealthchecksUnauthorisedException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }
}
