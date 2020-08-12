<?php

namespace Henrywhitaker3\Healthchecks\Tests;

use Henrywhitaker3\Healthchecks\Exceptions\InvalidUrlException;
use Henrywhitaker3\Healthchecks\Healthchecks;
use Henrywhitaker3\Healthchecks\HealthchecksManager;

class InstantiationTest extends Test
{
    /**
     * Tests creating HM with invalid URL
     *
     * @return void
     */
    public function testInvalidURL()
    {
        try {
            $hm = new HealthchecksManager('abc', 'abc');
        } catch(InvalidUrlException $e) {
            $this->assertTrue(true);
            return true;
        }

        $this->assertTrue(false);
    }
}