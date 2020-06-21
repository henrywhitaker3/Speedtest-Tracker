<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;

class SpeedtestVersionCommand extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'speedtest:version';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Displays the version number of this instance';

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
     * Execute the console command.
     *
     * @return mixed
     */
    public function handle()
    {
        $this->info('Speedtest Tracker');
        $this->info('Author: Henry Whitaker');
        $this->info('');
        $this->info('Installed version v' . config('speedtest.version'));
    }
}
