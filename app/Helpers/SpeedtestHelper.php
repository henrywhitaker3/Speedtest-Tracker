<?php

namespace App\Helpers;

use App\Speedtest;

class SpeedtestHelper {
    public static function runSpeedtest($output = false)
    {
        if($output === false) {
            $output = shell_exec('speedtest-cli');
        }
        $output = preg_replace("/\r|\n/", "", $output);

        $pattern = '/([0-9\.]{1,}) ms.*Download: ([0-9\.]{1,} [A-Za-z]{1,}\/s).*Upload: ([0-9\.]{1,} [A-Za-z]{1,}\/s)/';
        $matches = [];
        preg_match_all($pattern, $output, $matches);

        $ping = $matches[1][0];
        $down = SpeedtestHelper::parseUnits($matches[2][0]);
        $up = SpeedtestHelper::parseUnits($matches[3][0]);

        $test = Speedtest::create([
            'ping' => $ping,
            'download' => $down['val'],
            'upload' => $up['val']
        ]);

        return $test;
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
