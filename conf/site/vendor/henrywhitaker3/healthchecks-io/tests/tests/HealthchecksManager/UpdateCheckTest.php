<?php

namespace Henrywhitaker3\Healthchecks\Tests;

use Henrywhitaker3\Healthchecks\Exceptions\HealthchecksFailureException;
use Henrywhitaker3\Healthchecks\Exceptions\HealthchecksUnauthorisedException;
use Henrywhitaker3\Healthchecks\Exceptions\HealthchecksUuidNotFoundException;
use Henrywhitaker3\Healthchecks\Healthchecks;
use Henrywhitaker3\Healthchecks\HealthchecksManager;
use InvalidArgumentException;
use Ramsey\Uuid\Exception\InvalidUuidStringException;
use Ramsey\Uuid\Uuid;

class UpdateCheckTest extends Test
{
    /**
     * HM instance
     *
     * @var HealthchecksManager
     */
    private $hm;

    /**
     * Uuid of created check
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
     * Tests updating name
     *
     * @return void
     */
    public function testUpdateName()
    {
        $new = $this->hm->updateCheck($this->uuid, [
            'name' => 'testUpdateName'
        ]);

        $this->assertEquals('testUpdateName', $new['name']);
    }

    /**
     * Tests updating timeout
     *
     * @return void
     */
    public function testUpdateTimeout()
    {
        $new = $this->hm->updateCheck($this->uuid, [
            'timeout' => 3600
        ]);

        $this->assertEquals(3600, $new['timeout']);
    }

    /**
     * Tests updating cron
     *
     * @return void
     */
    public function testUpdateToCron()
    {
        $new = $this->hm->updateCheck($this->uuid, [
            'schedule' => '* * * * *'
        ]);

        $this->assertEquals('* * * * *', $new['schedule']);
    }

    /**
     * Tests updating cron->timeout
     *
     * @return void
     */
    public function testUpdateFromCronToTimeout()
    {
        $new = $this->hm->updateCheck($this->uuid, [
            'schedule' => '* * * * *'
        ]);
        $this->assertEquals('* * * * *', $new['schedule']);

        $new = $this->hm->updateCheck($this->uuid, [
            'timeout' => 3600
        ]);
        $this->assertEquals(3600, $new['timeout']);
    }

    /**
     * Tests invalid uuid
     *
     * @return void
     */
    public function testUpdateInvalidUuid()
    {
        try {
            $new = $this->hm->updateCheck('abc', []);
        } catch(InvalidUuidStringException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }

    /**
     * Tests invalid args param
     *
     * @return void
     */
    public function testUpdateInvalidArgs()
    {
        try {
            $new = $this->hm->updateCheck($this->uuid, [ 'name' => [ 'array not string' ] ]);
        } catch(InvalidArgumentException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }

    /**
     * Tests invalid 404
     *
     * @return void
     */
    public function testUpdate404()
    {
        try {
            $hm = new HealthchecksManager('abc', 'https://github.com/henrywhitaker3/404');
            $new = $hm->updateCheck(Uuid::uuid4(), []);
        } catch(HealthchecksUuidNotFoundException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }

    /**
     * Tests invalid unauth
     *
     * @return void
     */
    public function testUpdateUnauth()
    {
        try {
            $hm = new HealthchecksManager('abc');
            $new = $hm->updateCheck($this->uuid, []);
        } catch(HealthchecksUnauthorisedException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }
}
