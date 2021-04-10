<?php

namespace App\Interfaces;

use App\Models\Speedtest;

interface SpeedtestProvider
{
    public function run(): Speedtest;
    public function output();
}
