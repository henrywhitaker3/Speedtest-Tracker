<?php

namespace App\Helpers;

use App\Speedtest;
use DateTime;
use Exception;
use Illuminate\Support\Facades\Storage;

class BackupHelper {
    public static function backup($format = 'json')
    {
        $timestamp = new DateTime();
        $timestamp = $timestamp->format('Y-m-d_H:i:s');
        $name = 'speedtest_backup_' . $timestamp;

        switch($format) {
            case 'csv':
                $data = Speedtest::get();

                $csv = Storage::disk('local')->getDriver()->getAdapter()->getPathPrefix() . $name . '.csv';
                $name = $name . '.csv';
                $handle = fopen($csv, 'w+');
                fputcsv($handle, array('id', 'ping', 'download', 'upload', 'created_at', 'updated_at'));

                foreach($data as $d) {
                    fputcsv($handle, array($d->id, $d->ping, $d->download, $d->upload, $d->created_at, $d->updated_at));
                }

                fclose($handle);

                break;
            case 'json':
            default:
                $data = Speedtest::get()->toJson();
                $name = $name . '.json';
                Storage::disk('local')->put($name, $data);
                break;

        }

        return $name;
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
