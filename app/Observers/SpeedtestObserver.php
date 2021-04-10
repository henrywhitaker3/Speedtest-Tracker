<?php

namespace App\Observers;

use App\Exceptions\InfluxDBNotEnabledException;
use App\Models\Speedtest;
use App\Utils\InfluxDB\InfluxDB;
use Exception;
use Log;

class SpeedtestObserver
{
    /**
     * Handle the Speedtest "created" event.
     *
     * @param  \App\Speedtest  $speedtest
     * @return void
     */
    public function created(Speedtest $speedtest)
    {
        try {
            InfluxDB::connect()
                ->store($speedtest);
        } catch (InfluxDBNotEnabledException $e) {
            // /
        } catch (Exception $e) {
            Log::error($e);
        }
    }
}
