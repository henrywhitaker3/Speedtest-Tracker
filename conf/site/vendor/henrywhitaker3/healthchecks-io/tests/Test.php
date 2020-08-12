<?php

namespace Henrywhitaker3\Healthchecks\Tests;

use Dotenv\Dotenv;
use PHPUnit\Framework\TestCase;

class Test extends TestCase
{
    public function setUp() : void
    {
        $dotenv = Dotenv::createImmutable(__DIR__.'/../');
        $dotenv->load();
    }
}