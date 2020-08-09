<?php

namespace App\Helpers;

use App\Speedtest;
use DateTime;
use Exception;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Storage;

class BackupHelper {

    /**
     * Generates a backup of all speedtests.
     *
     * @param string    $format json|csv
     * @return string   $name   Returns the filename of the backup.
     */
    public static function backup(String $format = 'json')
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
                fputcsv($handle, array('id', 'ping', 'download', 'upload', 'created_at', 'updated_at', 'server_id', 'server_name', 'server_host', 'url', 'scheduled', 'failed'));

                foreach ($data as $d) {
                    fputcsv($handle, array($d->id, $d->ping, $d->download, $d->upload, $d->created_at, $d->updated_at, $d->server_id, $d->server_name, $d->server_host, $d->url, $d->scheduled, $d->failed));
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

    /**
     * Restore data from a backup in CSV or JSON format
     *
     * @param   array|string    $array  Backup data
     * @param   string          $format json|csv
     * @return  bool
     */
    public static function restore($array, $format)
    {
        if($format == 'json') {
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
            return true;
        } else if($format == 'csv') {
            $csv = explode(PHP_EOL, $array);
            $headers = 'id,ping,download,upload,created_at,updated_at';
            if($csv[0] != $headers) {
                Log::error('Incorrect CSV format');
                return false;
            }

            unset($csv[0]);
            $csv = array_values($csv);

            for($i = 0; $i < sizeof($csv); $i++) {
                $e = explode(',', $csv[$i]);
                try {
                    $st = Speedtest::create([
                        'ping' => $e[1],
                        'download' => $e[2],
                        'upload' => $e[3],
                        'created_at' => substr($e[4], 1, -1),
                    ]);
                } catch(Exception $e) {
                    Log::error($e);
                    continue;
                }
            }

            return true;
        }

        return false;
    }
}
