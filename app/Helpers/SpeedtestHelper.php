<?php

namespace App\Helpers;

use App\Speedtest;
use Carbon\Carbon;
use Exception;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use JsonException;

class SpeedtestHelper {

    /**
     * Runs/processes speedtest output to created a Speedtest object
     *
     * @param   boolean|string  $output If false, new speedtest runs. If anything else, will try to parse as JSON for speedtest results.
     * @return \App\Speedtest|boolean
     */
    public static function runSpeedtest($output = false, $scheduled = true)
    {
        if($output === false) {
            $output = SpeedtestHelper::output();
        }

        try {
            $output = json_decode($output, true, 512, JSON_THROW_ON_ERROR);

            if(!SpeedtestHelper::checkOutputIsComplete($output)) {
                return false;
            }

            $test = Speedtest::create([
                'ping' => $output['ping']['latency'],
                'download' => SpeedtestHelper::convert($output['download']['bandwidth']),
                'upload' => SpeedtestHelper::convert($output['upload']['bandwidth']),
                'server_id' => $output['server']['id'],
                'server_name' => $output['server']['name'],
                'server_host' => $output['server']['host'] . ':' . $output['server']['port'],
                'url' => $output['result']['url'],
                'scheduled' => $scheduled
            ]);
        } catch(JsonException $e) {
            Log::error('Failed to parse speedtest JSON');
            Log::error($output);
        } catch(Exception $e) {
            Log::error($e->getMessage());
        }

        return (isset($test)) ? $test : false;
    }

    /**
     * Gets the output of executing speedtest binary.
     *
     * @return boolean|string
     */
    public static function output()
    {
        $server = SettingsHelper::get('server')['value'];
        $binPath = app_path() . DIRECTORY_SEPARATOR . 'Bin' . DIRECTORY_SEPARATOR . 'speedtest';
        if($server != '' && $server != false) {
            $server = explode(',', $server);
            $server = $server[array_rand($server)];
            if($server == false) {
                Log::error('Speedtest server undefined');
                return false;
            }

            return shell_exec($binPath . ' -f json -s ' . $server);
        }

        return shell_exec($binPath . ' -f json');
    }

    /**
     * Get a 24 hour average of speedtest results
     *
     * @return array
     */
    public static function last24Hours()
    {
        $t = Carbon::now()->subDay();
        $s = Speedtest::select(DB::raw('AVG(ping) as ping, AVG(download) as download, AVG(upload) as upload'))
                      ->where('created_at', '>=', $t)
                      ->first()
                      ->toArray();

        return $s;
    }

    /**
     * Converts bytes/s to Mbps
     *
     * @param   int|float   $bytes
     * @return int|float
     */
    public static function convert($bytes) {
        return ( $bytes * 8 ) / 1000000;
    }

    /**
     * Returns the latest speedtest object.
     *
     * @return boolean|\App\Speedtest
     */
    public static function latest()
    {
        $data = Speedtest::latest()->get();

        if($data->isEmpty()) {
            return false;
        }

        return $data->first();
    }

    /**
     * Parses network speeds and return converted to Mbps
     *
     * @param   array  $input
     * @return  array
     */
    public static function parseUnits($input)
    {
        $input = explode(' ', $input);

        $val = $input[0];
        $unit = explode('/', $input[1])[0];

        switch($unit) {
            case 'Mbyte':
                $val = $val * 8;
                break;
            case 'Kbit':
                $val = $val / 1000;
                break;
            case 'Kbyte':
                $val = $val / 125;
                break;
            case 'Mbit':
            default:
                break;
        }

        return [
            'val' => $val,
            'unit' => 'Mbit/s'
        ];
    }

    /**
     * Checks that the speedtest JSON output is complete/valid
     *
     * @param array $output
     * @return boolean
     */
    public static function checkOutputIsComplete($output)
    {
        /**
         * Array of indexes that must exist in $output
         */
        $checks = [
            'type' => 'result',
            'download' => [ 'bandwidth' => '*' ],
            'upload' => [ 'bandwidth' => '*' ],
            'ping' => [ 'latency' => '*' ],
            'server' => [
                'id' => '*',
                'name' => '*',
                'host' => '*',
                'port' => '*'
            ],
            'result' => [
                'url' => '*'
            ],
        ];
        /**
         * Array of indexes that must not exist
         */
        $checkMissing = [
            'type' => 'error'
        ];

        foreach($checks as $key => $value) {
            if(!isset($output[$key])) {
                return false;
            }
        }

        return true;
    }
}
