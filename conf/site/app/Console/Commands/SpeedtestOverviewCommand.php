<?php

namespace App\Console\Commands;

use App\Events\SpeedtestOverviewEvent;
use App\Helpers\SpeedtestHelper;
use App\Notifications\SpeedtestOverviewSlack;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Notification;

class SpeedtestOverviewCommand extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'speedtest:overview';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Trigger a speedtest overview event';

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
        event(new SpeedtestOverviewEvent());
    }
}
