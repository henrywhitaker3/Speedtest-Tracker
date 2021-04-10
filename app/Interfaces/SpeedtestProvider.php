<?php

namespace App\Interfaces;

use App\Speedtest;

interface SpeedtestProvider
{
    public function run(): Speedtest;
    public function output();
}
