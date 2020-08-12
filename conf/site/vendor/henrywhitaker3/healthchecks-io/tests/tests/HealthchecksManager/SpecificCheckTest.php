<?php

namespace Henrywhitaker3\Healthchecks\Tests;

use Henrywhitaker3\Healthchecks\Exceptions\HealthchecksUuidNotFoundException;
use Henrywhitaker3\Healthchecks\HealthchecksManager;
use Ramsey\Uuid\Exception\InvalidUuidStringException;
use Ramsey\Uuid\Uuid;

class SpecificCheckTest extends Test
{
    /**
     * Instance of manager
     *
     * @var HealthchecksManager
     */
    private $hm;

    /**
     * UUID of created check
     *
     * @var String
     */
    private $uuid;

    public function setUp(): void
    {
        parent::setUp();

        $this->hm = new HealthchecksManager($_ENV['APIKEY']);
        $this->uuid = substr($this->hm->createCheck()['ping_url'], 20);
    }

    public function tearDown(): void
    {
        $this->hm->deleteCheck($this->uuid);
    }

    /**
     * Tests listing specific check from API
     *
     * @return void
     */
    public function testValidUuid()
    {
        $check = $this->hm->getCheck($this->uuid);

        $this->assertIsArray($check);
        $this->assertArrayHasKey('name', $check);
    }

    /**
     * Test pausing a check
     *
     * @return void
     */
    public function testPauseCheck()
    {
        if($this->hm->pauseCheck($this->uuid)) {
            $this->assertTrue(true);
        }
    }

    /**
     * Test resuming a check
     *
     * @return void
     */
    public function testResumeCheck()
    {
        if($this->hm->resumeCheck($this->uuid)) {
            $this->assertTrue(true);
        }
    }

    /**
     * Test deleting a check
     *
     * @return void
     */
    public function testDeleteCheck()
    {
        $uuid = substr($this->hm->createCheck()['ping_url'], 20);
        if($this->hm->deleteCheck($uuid)) {
            $this->assertTrue(true);
        }
    }

    /**
     * Test getCheckPings
     *
     * @return void
     */
    public function testGetCheckPings()
    {
        if($this->hm->getCheckPings($this->uuid)) {
            $this->assertTrue(true);
        }
    }

    /**
     * Test getCheckStatusChanges
     *
     * @return void
     */
    public function testGetCheckStatusChanges()
    {
        if($this->hm->getCheckStatusChanges($this->uuid)) {
            $this->assertTrue(true);
        }
    }

    /**
     * Tests listing invalid check form API
     *
     * @return void
     */
    public function testInvalidUuid()
    {
        try{
            $check = $this->hm->getCheck('abc');
        } catch(InvalidUuidStringException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }
}
