<?php

namespace App\Utils\InfluxDB;

use App\Interfaces\InfluxDBWrapperInterface;
use App\Models\Speedtest;
use Exception;
use InfluxDB\Client;
use InfluxDB\Database;
use InfluxDB\Database\RetentionPolicy;
use InfluxDB\Point;
use Log;

class InfluxDBVersion1Wrapper implements InfluxDBWrapperInterface
{
    private Client $client;
    private Database $database;

    public function __construct(Client $client)
    {
        $this->client = $client;
    }

    public function testConnection(): bool
    {
        try {
            $this->client->listDatabases();

            return true;
        } catch (Exception $e) {
            Log::error($e);

            return false;
        }
    }

    public function doesDatabaseExist(string $database): bool
    {
        $this->database = $this->client->selectDB($database);

        return (bool) $this->database->exists();
    }

    public function createDatabase(string $database): bool
    {
        try {
            $this->database->create(
                new RetentionPolicy(
                    'speedtest_retention_policy',
                    config('services.influxdb.retention'),
                    1,
                    true
                )
            );

            return true;
        } catch (Exception $e) {
            Log::error($e);
            return false;
        }
    }

    public function store(Speedtest $speedtest): bool
    {
        return $this->database->writePoints([
            new Point(
                'speedtest',
                null,
                ['host' => config('services.influxdb.host')],
                $speedtest->formatForInfluxDB(),
            )
        ]);
    }
}
