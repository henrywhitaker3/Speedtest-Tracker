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
        $server = SettingsHelper::get('server')['value'];
        $binPath = app_path() . DIRECTORY_SEPARATOR . 'Bin' . DIRECTORY_SEPARATOR . 'SpeedTest';
        if($server != '' && $server != false) {
            $server = explode(',', $server);
            $server = $server[array_rand($server)];
            $server = SpeedtestHelper::resolveServer($server);
            if($server == false) {
                Log::error('Speedtest server undefined');
                return false;
            }
            $server = $server['host'];

            return shell_exec($binPath . ' --output json --test-server ' . $server);
        }

        return shell_exec($binPath . ' --output json');
    }

    /*
    * Resolve the server host/port from speedtest server id
    */
    public static function resolveServer($id)
    {
        $ttl = Carbon::now()->addMinutes(120);
        $ids = Cache::remember('servers', $ttl, function () {
            $urls = [
                'http://www.speedtest.net/speedtest-servers-static.php',
                'http://c.speedtest.net/speedtest-servers-static.php',
                'http://www.speedtest.net/speedtest-servers.php',
                'http://c.speedtest.net/speedtest-servers.php',
            ];

            $client = new Client([]);

            $ids = [];
            foreach($urls as $url) {
                try {
                    $response = $client->get($url);
                    if($response->getStatusCode() == 200) {
                        $data = (string) $response->getBody();
                        $data = simplexml_load_string($data);
                        $data = SpeedtestHelper::xmlToArray($data);
                        $data = $data['settings']['servers']['server'];
                        foreach($data as $s) {
                            $ids[$s['attributes']['id']] = $s['attributes'];
                        }
                    } else {
                        continue;
                    }
                } catch(Exception $e) {
                    Log::error('Speedtest server resolver error');
                    Log::error($e->getMessage());
                    continue;
                }
            }
            return $ids;
        });

        if(array_key_exists($id, $ids)) {
            return $ids[$id];
        } else {
            return false;
        }
    }

    public static function xmlToArray(SimpleXMLElement $xml): array
    {
        $parser = function (SimpleXMLElement $xml, array $collection = []) use (&$parser) {
            $nodes = $xml->children();
            $attributes = $xml->attributes();

            if (0 !== count($attributes)) {
                foreach ($attributes as $attrName => $attrValue) {
                    $collection['attributes'][$attrName] = strval($attrValue);
                }
            }

            if (0 === $nodes->count()) {
                $collection['value'] = strval($xml);
                return $collection;
            }

            foreach ($nodes as $nodeName => $nodeValue) {
                if (count($nodeValue->xpath('../' . $nodeName)) < 2) {
                    $collection[$nodeName] = $parser($nodeValue);
                    continue;
                }

                $collection[$nodeName][] = $parser($nodeValue);
            }

            return $collection;
        };

        return [
            $xml->getName() => $parser($xml)
        ];
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
