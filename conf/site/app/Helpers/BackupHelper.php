<?php

namespace App\Helpers;

use App\Speedtest;
use Exception;

class BackupHelper {
    public static function backup()
    {
        $data = Speedtest::get();

        return $data;
    }

    public static function restore($array)
    {
        foreach($array as $test) {
            try {
                $st = Speedtest::create([
                    'ping' => $test['ping'],
                    'download' => $test['download'],
                    'upload' => $test['upload'],
                    'created_at' => $test['created_at'],
                ]);
            } catch(Exception $e) {
                continue;
            }
        }
    }
}
