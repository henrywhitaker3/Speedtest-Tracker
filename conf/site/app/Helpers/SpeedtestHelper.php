<?php

namespace App\Helpers;

use App\Speedtest;
use Carbon\Carbon;
use Exception;
use Henrywhitaker3\Healthchecks\Healthchecks;
use Illuminate\Support\Facades\File;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use InvalidArgumentException;
use JsonException;

class SpeedtestHelper {

    /**
     * Runs/processes speedtest output to created a Speedtest object
     *
     * @param   boolean|string  $output If false, new speedtest runs. If anything else, will try to parse as JSON for speedtest results.
     * @return \App\Speedtest|bool
     */
    public static function runSpeedtest($output = false, $scheduled = true)
    {
        if($output === false) {
            $output = SpeedtestHelper::output();
        }

        try {
            $output = json_decode($output, true, 512, JSON_THROW_ON_ERROR);

            if(!SpeedtestHelper::checkOutputIsComplete($output)) {
                $test = false;
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
            $test = false;
        } catch(Exception $e) {
            Log::error($e->getMessage());
            $test = false;
        }

        if($test == false) {
            Speedtest::create([
                'ping' => 0,
                'upload' => 0,
                'download' => 0,
                'failed' => true,
                'scheduled' => $scheduled,
            ]);

            return false;
        }

        Cache::flush();

        return $test;
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
        $homePrefix = config('speedtest.home') . ' && ';

        if($server != '' && $server != false) {
            $server = explode(',', $server);
            $server = $server[array_rand($server)];
            if($server == false) {
                Log::error('Speedtest server undefined');
                return false;
            }

            return shell_exec($homePrefix . $binPath . ' -f json -s ' . $server);
        }

        return shell_exec($homePrefix . $binPath . ' -f json');
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
                      ->where('failed', false)
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
     * @param   string  $input
     * @return  array
     */
    public static function parseUnits($input)
    {
        $input = explode(' ', $input);

        $val = (float)$input[0];
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

    /**
     * Get a percentage rate of failure by days
     *
     * @param integer $days number of days to get rate for
     * @return integer percentage fail rate
     */
    public static function failureRate(int $days)
    {
        $ttl = Carbon::now()->addDays(1);
        $rate = Cache::remember('failure-rate-' . $days, $ttl, function () use ($days) {
            $range = [
                Carbon::today()
            ];
            for($i = 0; $i < ($days - 1); $i++) {
                $prev = end($range);
                $new = $prev->copy()->subDays(1);
                array_push($range, $new);
            }

            $rate = [];

            foreach($range as $day) {
                $success = Speedtest::select(DB::raw('COUNT(id) as rate'))->whereDate('created_at', $day)->where('failed', false)->get()[0]['rate'];
                $fail = Speedtest::select(DB::raw('COUNT(id) as rate'))->whereDate('created_at', $day)->where('failed', true)->get()[0]['rate'];

                array_push($rate, [
                    'date' => $day->toDateString(),
                    'success' => $success,
                    'failure' => $fail,
                ]);
            }

            return array_reverse($rate);
        });

        return $rate;
    }

    /**
     * Create a backup of the SQLite database
     *
     * @return null|boolean
     */
    public static function dbBackup()
    {
        if(env('DB_CONNECTION') === 'sqlite') {
            if(env('DB_DATABASE') !== null) {
                $current = env('DB_DATABASE');
                try {
                    if(File::copy($current, $current . '.bak')) {
                        return true;
                    }
                }catch(Exception $e) {
                    return false;
                }
            }

            return false;
        }

        return null;
    }

    /**
     * Delete all speedtests from the database
     *
     * @return array
     */
    public static function deleteAll()
    {
        Cache::flush();

        SpeedtestHelper::dbBackup();

        if(sizeof(Speedtest::whereNotNull('id')->get()) > 0) {
            if(Speedtest::whereNotNull('id')->delete()) {
                return [
                    'success' => true,
                ];
            }
        }

        return [
            'success' => true,
        ];
    }

    /**
     * Work out if a test is lower than the threshold for historic tests
     *
     * @param String $type
     * @param Speedtest $test
     * @return array
     */
    public static function testIsLowerThanThreshold(String $type, Speedtest $test)
    {
        if($type == 'percentage') {
            $avg = Speedtest::select(DB::raw('AVG(ping) as ping, AVG(download) as download, AVG(upload) as upload'))
                        ->where('failed', false)
                        ->get()
                        ->toArray()[0];

            $threshold = SettingsHelper::get('threshold_alert_percentage')->value;

            if($threshold == '') {
                return [];
            }

            $errors = [];

            foreach($avg as $key => $value) {
                if($key == 'ping') {
                    $threshold = (float)$value * (1 + ( $threshold / 100 ));

                    if($test->$key > $threshold) {
                        array_push($errors, $key);
                    }
                } else  {
                    $threshold = (float)$value * (1 - ( $threshold / 100 ));

                    if($test->$key < $threshold) {
                        array_push($errors, $key);
                    }
                }
            }

            return $errors;
        }

        if($type == 'absolute') {
            $thresholds = [
                'download' => SettingsHelper::get('threshold_alert_absolute_download')->value,
                'upload' => SettingsHelper::get('threshold_alert_absolute_upload')->value,
                'ping' => SettingsHelper::get('threshold_alert_absolute_ping')->value,
            ];

            $errors = [];

            foreach($thresholds as $key => $value) {
                if($value == '') {
                    continue;
                }

                if($key == 'ping') {
                    if($test->$key > $value) {
                        array_push($errors, $key);
                    }
                } else {
                    if($test->$key < $value) {
                        array_push($errors, $key);
                    }
                }
            }

            return $errors;
        }

        throw new InvalidArgumentException();
    }
}
