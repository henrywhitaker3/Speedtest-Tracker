<?php

namespace App\Utils\InfluxDB;

use App\Exceptions\InfluxDBConnectionErrorException;
use App\Exceptions\InfluxDBNotEnabledException;
use App\Helpers\SettingsHelper;
use App\Models\Speedtest;
use InfluxDB\Client as Version1;
use InfluxDB2\Client as Version2;
use App\Utils\InfluxDB\InfluxDBWrapperInterface as Client;

class InfluxDB
{
    private Client $client;

    private string $database;

    public function __construct(Client $client)
    {
        $this->client = $client;
    }

    /**
     * Connect to influx db
     *
     * @param string $host
     * @param integer $port
     * @param string $database
     * @return InfluxDB
     */
    public static function connect()
    {
        if (!(bool) SettingsHelper::get('influx_db_enabled')->value) {
            throw new InfluxDBNotEnabledException();
        }

        $host = SettingsHelper::get('influx_db_host')->value;
        $port = SettingsHelper::get('influx_db_port')->value;
        $token = '';
        $database = SettingsHelper::get('influx_db_database')->value;
        $version = (int) SettingsHelper::get('influx_db_version')->value;

        $wrapper = $version === 1
            ? new InfluxDBVersion1Wrapper(
                new Version1(str_replace(['http://', 'https://'], '', $host), $port)
            )
            : new InfluxDBVersion2Wrapper(
                new Version2([])
            );

        return (new self($wrapper))->setDatabase($database)
            ->testConnection();
    }

    /**
     * Set the database field
     *
     * @param string $database
     * @return InfluxDB
     */
    public function setDatabase(string $database): InfluxDB
    {
        $this->database = $database;

        return $this;
    }

    /**
     * Test the connection
     *
     * @throws InfluxDBConnectionErrorException
     * @return InfluxDB
     */
    public function testConnection(): InfluxDB
    {
        if (!$this->client->testConnection()) {
            throw new InfluxDBConnectionErrorException();
        }

        if (!$this->doesDatabaseExist()) {
            $this->createDatabase();
        }

        return $this;
    }

    public function doesDatabaseExist(): bool
    {
        return $this->client->doesDatabaseExist($this->database);
    }

    public function createDatabase(): bool
    {
        return $this->client->createDatabase($this->database);
    }

    public function store(Speedtest $speedtest): InfluxDB
    {
        $this->client->store($speedtest);

        return $this;
    }
}
