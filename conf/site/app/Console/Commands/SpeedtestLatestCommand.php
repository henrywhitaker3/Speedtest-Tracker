<?php

namespace App\Console\Commands;

use App\Helpers\SpeedtestHelper;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Artisan;

class SpeedtestLatestCommand extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'speedtest:latest';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Returns the latest speedtest result';

    /**
     * Create a new command instance.
     *
     * @return void
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * Prints the latest speedtest values.
     *
     * @return void
     */
    public function handle()
    {
        $latest = SpeedtestHelper::latest();

        if($latest) {
            if($latest->scheduled) {
                $extra = '(scheduled)';
            } else {
                $extra = '(manual)';
            }

            $this->info('Last speedtest run at: ' . $latest->created_at . ' ' . $extra);

            if($latest->failed) {
                $this->error('Speedtest failed');
            } else {
                $this->info('Ping: ' . $latest->ping . ' ms');
                $this->info('Download: ' . $latest->download . ' Mbit/s');
                $this->info('Upload: ' . $latest->upload . ' Mbit/s');
            }
        } else {
            $this->info('No speedtests have been run yet.');

            $this->info('Running speedtest, this might take a while...');

            $results = SpeedtestHelper::runSpeedtest();

            $this->info('Ping: ' . $results->ping . ' ms');
            $this->info('Download: ' . $results->download . ' Mbit/s');
            $this->info('Upload: ' . $results->upload . ' Mbit/s');
        }
    }
}
