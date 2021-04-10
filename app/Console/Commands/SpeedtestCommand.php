<?php

namespace App\Console\Commands;

use App\Exceptions\SpeedtestFailureException;
use App\Helpers\SpeedtestHelper;
use App\Interfaces\SpeedtestProvider;
use Illuminate\Console\Command;

class SpeedtestCommand extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'speedtest:run';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Performs a new speedtest';

    private SpeedtestProvider $speedtestProvider;

    /**
     * Create a new command instance.
     *
     * @return void
     */
    public function __construct(SpeedtestProvider $speedtestProvider)
    {
        $this->speedtestProvider = $speedtestProvider;

        parent::__construct();
    }

    /**
     * Runs a speedtest synchroonously and displays the results..
     *
     * @return void
     */
    public function handle()
    {
        $this->info('Running speedtest, this might take a while...');

        try {
            $results = $this->speedtestProvider->run(false, false);
        } catch (SpeedtestFailureException $e) {
            $this->error('Something went wrong running the speedtest.');
            exit();
        }

        if (property_exists($results, 'ping') && property_exists($results, 'download') && property_exists($results, 'upload')) {
            $this->error('Something went wrong running the speedtest.');
            exit();
        }

        $this->info('Ping: ' . $results->ping . ' ms');
        $this->info('Download: ' . $results->download . ' Mbit/s');
        $this->info('Upload: ' . $results->upload . ' Mbit/s');
    }
}
