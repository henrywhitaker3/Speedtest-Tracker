<?php

namespace App\Helpers;

use App\Speedtest;
use Cache;
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

                $csv = storage_path() . '/app/' . $name . '.csv';
                $name = $name . '.csv';
                $handle = fopen($csv, 'w+');
                fputcsv($handle, array('id', 'ping', 'download', 'upload', 'server_id', 'server_name', 'server_host', 'url', 'scheduled', 'failed', 'created_at', 'updated_at'));

                foreach($data as $d) {
                    fputcsv($handle, BackupHelper::createCSVBackupArray($d));
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
        Cache::flush();
        if($format == 'json') {
            foreach($array as $test) {
                try {
                    $data = BackupHelper::backupJSONToArray($test);

                    if($data === false) {
                        continue;
                    }

                    Speedtest::create($data);
                } catch(Exception $e) {
                    Log::error($e);
                    continue;
                }
            }
            return true;
        } else if($format == 'csv') {
            $csv = explode(PHP_EOL, $array);

            $headers = BackupHelper::validateCSV($csv[0]);
            if($headers === false) {
                return false;
            }

            unset($csv[0]);
            $csv = array_values($csv);

            for($i = 0; $i < sizeof($csv); $i++) {
                $data = BackupHelper::backupCSVToArray($csv[$i]);

                if($data === false) {
                    continue;
                }

                try {
                    Speedtest::create($data);
                } catch(Exception $e) {
                    Log::error($e);
                    continue;
                }
            }

            return true;
        }

        return false;
    }

    /**
     * Validate a CSV file passed for restore
     *
     * @param String $csv The line containing the CSV headers
     * @return bool|string
     */
    public static function validateCSV(String $csv)
    {
        $headers = [
            'old' => 'id,ping,download,upload,created_at,updated_at',
            'new' => 'id,ping,download,upload,server_id,server_name,server_host,url,scheduled,failed,created_at,updated_at',
        ];
        $backupHeaders = null;

        foreach($headers as $key => $h) {
            if($csv == $h) {
                $backupHeaders = $key;
            }
        }

        if($backupHeaders === null) {
            Log::info('Incorrect CSV format');
            return false;
        }

        return $backupHeaders;
    }

    /**
     * Return an array from the raw CSV data
     *
     * @param String $line The line of CSV data
     * @param String $header The type of backup header
     * @return array|bool
     */
    public static function backupCSVToArray(String $line, String $header = 'new')
    {
        $basic = explode(',', $line);

        if($header == 'old') {
            $array = [
                'ping' => $basic[1],
                'download' => $basic[2],
                'upload' => $basic[3],
                'created_at' => substr($basic[4], 1, -1),
            ];
        }

        if($header == 'new') {
            $array = [
                'ping' => $basic[1],
                'download' => $basic[2],
                'upload' => $basic[3],
                'server_id' => $basic[4],
                'server_name' => $basic[5],
                'server_host' => $basic[6],
                'url' => $basic[7],
                'scheduled' => $basic[8],
                'failed' => $basic[9],
                'created_at' => substr($basic[10], 1, -1),
            ];
        }

        if(!isset($array)) {
            return false;
        }

        return BackupHelper::cleanRestoreDataArray($array);
    }

    /**
     * Clean an array, setting values with '' to null
     *
     * @param array $array
     * @return array
     */
    public static function cleanRestoreDataArray(array $array)
    {
        foreach($array as $key => $val) {
            if($val === '') {
                $array[$key] = null;
            }
        }

        return $array;
    }

    /**
     * Return an array from the JSON data
     *
     * @param array $json json_decoded data
     * @return array|bool
     */
    public static function backupJSONToArray($json)
    {
        $required = [
            'ping',
            'upload',
            'download',
            'created_at',
        ];

        $extras = [
            'server_id',
            'server_name',
            'server_host',
            'url',
            'failed',
            'scheduled'
        ];

        $array = [];

        foreach($required as $req) {
            if(!array_key_exists($req, $json)) {
                return false;
            }

            $val = $json[$req];

            if($val === '') {
                $val = null;
            }

            $array[$req] = $val;
        }

        foreach($extras as $extra) {
            if(array_key_exists($extra, $json)) {
                $val = $json[$extra];

                if($val === '') {
                    $val = null;
                }
                $array[$extra] = $val;
            }
        }

        return $array;
    }

    /**
     * Return an array to store in CSV
     *
     * @param Speedtest $test
     * @return array
     */
    public static function createCSVBackupArray(Speedtest $test)
    {
        $data = [
            $test->id,
            $test->ping,
            $test->download,
            $test->upload,
            $test->server_id,
            $test->server_name,
            $test->server_host,
            $test->url,
            $test->scheduled,
            $test->failed,
            $test->created_at,
            $test->updated_at
        ];

        foreach($data as $key => $val) {
            if(strpos($val, ',') !== false) {
                $val = str_replace(',', ' -', $val);
            }

            $data[$key] = $val;
        }

        return $data;
    }
}
