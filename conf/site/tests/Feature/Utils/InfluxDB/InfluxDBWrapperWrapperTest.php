<?php

namespace Tests\Feature\Utils\InfluxDB;

use App\Exceptions\InfluxDBNotEnabledException;
use App\Utils\InfluxDB\InfluxDB;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Tests\TestCase;

class InfluxDBWrapperWrapperTest extends TestCase
{
    use RefreshDatabase;

    /**
     * Test it throws the right exception.
     *
     * @return void
     */
    public function test_it_throws_excetpion_when_it_is_disabled()
    {
        $this->expectException(InfluxDBNotEnabledException::class);

        InfluxDB::connect();
    }
}
