<?php

namespace App\Helpers;

use App\Speedtest;
use Exception;
use Illuminate\Support\Facades\Log;
use JsonException;

class SpeedtestHelper {
    public static function runSpeedtest($output = false)
    {
        if($output === false) {
            $output = SpeedtestHelper::output();
        }

        try {
            $output = json_decode($output, true, 512, JSON_THROW_ON_ERROR);
            $test = Speedtest::create([
                'ping' => $output['ping'],
                'download' => $output['download'] / 1000000,
                'upload' => $output['upload'] / 1000000
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
        return shell_exec('speedtest-cli --json');
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
