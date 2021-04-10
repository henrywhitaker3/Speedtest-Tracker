<?php

namespace Database\Seeders;

use App\Models\Speedtest;
use Illuminate\Database\Seeder;

class SpeedtestSeeder extends Seeder
{
    /**
     * Run the database seeds.
     *
     * @return void
     */
    public function run()
    {
        Speedtest::factory()
            ->count(250)
            ->create();
    }
}
