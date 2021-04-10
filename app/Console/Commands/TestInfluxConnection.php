<?php

namespace App\Console\Commands;

use App\Exceptions\InfluxDBConnectionErrorException;
use App\Helpers\SettingsHelper;
use App\Utils\InfluxDB\InfluxDB;
use Illuminate\Console\Command;

class TestInfluxConnection extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'speedtest:influx';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Test connection settings for InfluxDB';

    private bool $enabled;

    /**
     * Create a new command instance.
     *
     * @return void
     */
    public function __construct()
    {
        parent::__construct();

        $this->enabled = (bool) SettingsHelper::get('influx_db_enabled')->value;
    }

    /**
     * Execute the console command.
     *
     * @return int
     */
    public function handle()
    {
        if (!$this->enabled) {
            $this->warn('InfluxDB is not enabled');
            exit;
        }

        try {
            InfluxDB::connect();
            $this->info('Connected successfully');
        } catch (InfluxDBConnectionErrorException $e) {
            $this->error('Couldn\'t connect');
            return 1;
        }

        return 0;
    }
}
