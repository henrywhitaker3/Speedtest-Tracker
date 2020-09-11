<?php

namespace Tests\Feature;

use App\Speedtest;
use Faker\Factory;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Illuminate\Support\Facades\DB;
use Tests\TestCase;

class APISpeedtestTest extends TestCase
{
    use RefreshDatabase;

    /**
     * Gets latest speedtest values
     *
     * @test
     * @return void
     */
    public function getLatestSpeedtest()
    {
        $faker = Factory::create();

        $ping = [];
        $dl = [];
        $ul = [];

        for ($i = 0; $i < 3; $i++) {
            $pingVal = $faker->randomFloat();
            array_push($ping, $pingVal);
            $dlVal = $faker->randomFloat();
            array_push($dl, $dlVal);
            $ulVal = $faker->randomFloat();
            array_push($ul, $ulVal);

            Speedtest::create([
                'ping' => $pingVal,
                'download' => $dlVal,
                'upload' => $ulVal,
            ]);
        }

        $avgVals = Speedtest::select(DB::raw('AVG(ping) as ping, AVG(download) as download, AVG(upload) as upload'))->get()[0];
        $maxVals = Speedtest::select(DB::raw('MAX(ping) as ping, MAX(download) as download, MAX(upload) as upload'))->get()[0];

        $pingAvg = $avgVals['ping'];
        $dlAvg = $avgVals['download'];
        $ulAvg = $avgVals['upload'];
        $pingMax = $maxVals['ping'];
        $dlMax = $maxVals['download'];
        $ulMax = $maxVals['upload'];

        $response = $this->get('/api/speedtest/latest');

        $response->assertStatus(200);
        $response->assertJsonStructure([
            'method',
            'data' => [
                'id',
                'ping',
                'download',
                'upload',
                'created_at',
                'updated_at',
            ],
            'average' => [
                'ping',
                'download',
                'upload',
            ],
            'maximum' => [
                'ping',
                'download',
                'upload',
            ],
            'minimum' => [
                'ping',
                'download',
                'upload',
            ],
        ]);
    }
}
