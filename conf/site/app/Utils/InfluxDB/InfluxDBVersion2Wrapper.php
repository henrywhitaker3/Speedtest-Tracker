<?php

namespace App\Utils\InfluxDB;

use App\Interfaces\InfluxDBWrapperInterface;
use App\Models\Speedtest;
use InfluxDB2\Client;

class InfluxDBVersion2Wrapper implements InfluxDBWrapperInterface
{
    private Client $client;

    public function __construct(Client $client)
    {
        $this->client = $client;
    }

    public function testConnection(): bool
    {
        return $this->client->health()->getStatus() !== 'pass';
    }

    public function doesDatabaseExist(string $database): bool
    {
        return true;
    }

    public function createDatabase(string $database): bool
    {
        return true;
    }

    public function store(Speedtest $speedtest): bool
    {
        return true;
    }
}
