<?php

namespace App\Utils;

use App\Exceptions\SpeedtestFailureException;
use App\Helpers\SettingsHelper;
use App\Helpers\SpeedtestHelper;
use App\Interfaces\SpeedtestProvider;
use App\Models\Speedtest;
use Cache;
use Exception;
use JsonException;
use Log;

class OoklaTester implements SpeedtestProvider
{
    public function run($output = false, $scheduled = true): Speedtest
    {
        if ($output === false) {
            $output = $this->output();
        }

        try {
            $output = json_decode($output, true, 512, JSON_THROW_ON_ERROR);

            if (!$this->isOutputComplete($output)) {
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
        } catch (JsonException $e) {
            Log::error('Failed to parse speedtest JSON');
            Log::error($output);
            $test = false;
        } catch (Exception $e) {
            Log::error($e->getMessage());
            $test = false;
        }

        if ($test == false) {
            Speedtest::create([
                'ping' => 0,
                'upload' => 0,
                'download' => 0,
                'failed' => true,
                'scheduled' => $scheduled,
            ]);

            throw new SpeedtestFailureException(json_encode($output));
        }

        Cache::flush();

        return $test;
    }

    public function output()
    {
        $server = SettingsHelper::get('server')['value'];

        $binPath = app_path() . DIRECTORY_SEPARATOR . 'Bin' . DIRECTORY_SEPARATOR . 'speedtest';
        $homePrefix = config('speedtest.home') . ' && ';

        if ($server != '' && $server != false) {
            $server = explode(',', $server);
            $server = $server[array_rand($server)];
            if ($server == false) {
                Log::error('Speedtest server undefined');
                return false;
            }

            return shell_exec($homePrefix . $binPath . ' -f json -s ' . $server);
        }

        return shell_exec($homePrefix . $binPath . ' -f json');
    }

    /**
     * Checks that the speedtest JSON output is complete/valid
     *
     * @param array $output
     * @return boolean
     */
    public static function isOutputComplete($output)
    {
        /**
         * Array of indexes that must exist in $output
         */
        $checks = [
            'type' => 'result',
            'download' => ['bandwidth' => '*'],
            'upload' => ['bandwidth' => '*'],
            'ping' => ['latency' => '*'],
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

        foreach ($checks as $key => $value) {
            if (!isset($output[$key])) {
                return false;
            }
        }

        return true;
    }
}
