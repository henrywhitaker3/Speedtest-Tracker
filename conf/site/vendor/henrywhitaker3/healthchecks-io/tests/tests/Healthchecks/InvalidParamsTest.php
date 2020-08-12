<?php

namespace Henrywhitaker3\Healthchecks\Tests;

use Henrywhitaker3\Healthchecks\Exceptions\HealthchecksUuidNotFoundException;
use Henrywhitaker3\Healthchecks\Exceptions\InvalidUrlException;
use Henrywhitaker3\Healthchecks\Healthchecks;
use Ramsey\Uuid\Exception\InvalidUuidStringException;
use Ramsey\Uuid\Uuid;

class InvalidUuidTest extends Test
{
    /**
     * Tests an invalid UUID for constructor
     *
     * @return void
     */
    public function testInvalidUuid()
    {
        try {
            $hc = new Healthchecks('');
        } catch(InvalidUuidStringException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }

    /**
     * Tets an invalid URL for constructor
     *
     * @return void
     */
    public function testInvalidUrl()
    {
        try {
            $uuid = Uuid::uuid4();
            $hc = new Healthchecks($uuid, 'test');
        } catch(InvalidUrlException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }

    /**
     * Tets an 404 URL for success
     *
     * @return void
     */
    public function test404UrlSuccess()
    {
        try {
            $uuid = Uuid::uuid4();
            $hc = new Healthchecks($uuid, 'https://github.com/henrywhitaker3/404');
            $hc->success();
        } catch(HealthchecksUuidNotFoundException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }

    /**
     * Tets an 404 URL for fail
     *
     * @return void
     */
    public function test404UrlFail()
    {
        try {
            $uuid = Uuid::uuid4();
            $hc = new Healthchecks($uuid, 'https://github.com/henrywhitaker3/404');
            $hc->fail();
        } catch(HealthchecksUuidNotFoundException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }

    /**
     * Tets an 404 URL for start
     *
     * @return void
     */
    public function test404UrlStart()
    {
        try {
            $uuid = Uuid::uuid4();
            $hc = new Healthchecks($uuid, 'https://github.com/henrywhitaker3/404');
            $hc->start();
        } catch(HealthchecksUuidNotFoundException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }
}