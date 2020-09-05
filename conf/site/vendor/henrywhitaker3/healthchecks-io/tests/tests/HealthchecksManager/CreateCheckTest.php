<?php

namespace Henrywhitaker3\Healthchecks\Tests;

use Exception;
use Henrywhitaker3\Healthchecks\Healthchecks;
use Henrywhitaker3\Healthchecks\HealthchecksManager;
use InvalidArgumentException;

class CreateCheckTest extends Test
{
    /**
     * Array of UUIDs to destroy on tearDown
     *
     * @var array
     */
    private $uuids;

    /**
     * HealthchecksManager instance
     *
     * @var HealthchecksManager
     */
    private $hm;

    public function setup(): void
    {
        $this->uuids = [];
        $this->hm = new HealthchecksManager($_ENV['APIKEY']);
    }

    public function tearDown(): void
    {
        foreach($this->uuids as $uuid) {
            try {
                $this->hm->deleteCheck($uuid);
            } catch(Exception $e) {
                continue;
            }
        }
    }

    /**
     * Tests creating a new check
     *
     * @return void
     */
    public function testCreateCheckNamed()
    {
        $args = [
            'name' => 'Test check named'
        ];

        $check = $this->hm->createCheck($args);
        array_push($this->uuids, substr($check['ping_url'], 20));

        $this->assertEquals('Test check named', $check['name']);
        $this->assertEquals(86400, $check['timeout']);
        $this->assertEquals(3600, $check['grace']);
        $this->assertEquals('', $check['desc']);
    }

    /**
     * Tests creating a new check with invalid args
     *
     * @return void
     */
    public function testCreateCheckNamedInvalidArgs()
    {
        $args = [
            'name' => 5/100
        ];

        try {
            $check = $this->hm->createCheck($args);
        } catch(InvalidArgumentException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }

    /**
     * Tests creating a new check with timeout specified
     *
     * @return void
     */
    public function testCreateCheckTimeout()
    {
        $args = [
            'timeout' => 500
        ];

        $check = $this->hm->createCheck($args);
        array_push($this->uuids, substr($check['ping_url'], 20));

        $this->assertEquals(500, $check['timeout']);
        $this->assertEquals(3600, $check['grace']);
        $this->assertEquals('', $check['desc']);
    }

    /**
     * Tests creating a new check with CRON specified
     *
     * @return void
     */
    public function testCreateCheckCron()
    {
        $args = [
            'schedule' => '* * * * *'
        ];

        $check = $this->hm->createCheck($args);
        array_push($this->uuids, substr($check['ping_url'], 20));

        $this->assertEquals('* * * * *', $check['schedule']);
        $this->assertEquals(3600, $check['grace']);
        $this->assertEquals('', $check['desc']);
    }
}