<?php

namespace Database\Seeders;

use Database\Seeders\SpeedtestSeeder;
use Illuminate\Database\Seeder;

class DatabaseSeeder extends Seeder
{
    /**
     * Seed the application's database.
     *
     * @return void
     */
    public function run()
    {
        $this->call(SpeedtestSeeder::class);
    }
}
