<?php

namespace App\Helpers;

use App\Speedtest;
use Carbon\Carbon;
use Exception;
use GuzzleHttp\Client;
use GuzzleHttp\Psr7\Request;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Log;
use JsonException;
use SimpleXMLElement;

class SpeedtestHelper {
    public static function runSpeedtest($output = false)
    {
        if($output === false) {
            $output = SpeedtestHelper::output();
        }

        try {
            $output = json_decode($output, true, 512, JSON_THROW_ON_ERROR);
            $test = Speedtest::create([
                'ping' => $output['ping']['latency'],
                'download' => SpeedtestHelper::convert($output['download']['bandwidth']),
                'upload' => SpeedtestHelper::convert($output['upload']['bandwidth']),
            ]);
        } catch(JsonException $e) {
            Log::error('Failed to parse speedtest JSON');
            Log::error($output);
        } catch(Exception $e) {
            Log::error($e->getMessage());
        }

        return (isset($test)) ? $test : false;
    }

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

            return shell_exec('HOME=/config && ' . $binPath . ' -f json -s ' . $server);
        }

        return shell_exec('HOME=/config && ' . $binPath . ' -f json');
    }

    public static function convert($bytes) {
        return ( $bytes * 8 ) / 1000000;
    }

    public static function latest()
    {
        $data = Speedtest::latest()->get();

        if($data->isEmpty()) {
            return false;
        }

        return $data->first();
    }

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
}
