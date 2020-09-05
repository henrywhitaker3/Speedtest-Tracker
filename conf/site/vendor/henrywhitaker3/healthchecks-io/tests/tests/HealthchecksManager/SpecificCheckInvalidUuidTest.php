<?php

namespace Henrywhitaker3\Healthchecks\Tests;

use Henrywhitaker3\Healthchecks\HealthchecksManager;
use Ramsey\Uuid\Exception\InvalidUuidStringException;
use Ramsey\Uuid\Uuid;

class SpecificCheckInvalidUuidTest extends Test
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
        $this->uuid = 'abc';
    }

    /**
     * Tests listing specific check from API
     *
     * @return void
     */
    public function testGetCheckInvalidUuid()
    {
        try {
            $this->hm->getCheck($this->uuid);
        } catch(InvalidUuidStringException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }

    /**
     * Test resuming a check
     *
     * @return void
     */
    public function testResumeCheckInvalidUuid()
    {
        try {
            $this->hm->resumeCheck($this->uuid);
        } catch(InvalidUuidStringException $e) {
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
    public function testPauseCheckInvalidUuid()
    {
        try {
            $this->hm->pauseCheck($this->uuid);
        } catch(InvalidUuidStringException $e) {
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
    public function testDeleteCheckInvalidUuid()
    {
        try {
            $this->hm->deleteCheck($this->uuid);
        } catch(InvalidUuidStringException $e) {
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
    public function testGetCheckPingsInvalidUuid()
    {
        try {
            $this->hm->getCheckPings($this->uuid);
        } catch(InvalidUuidStringException $e) {
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
    public function testGetCheckStatusChangesInvalidUuid()
    {
        try {
            $this->hm->getCheckStatusChanges($this->uuid);
        } catch(InvalidUuidStringException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }
}
