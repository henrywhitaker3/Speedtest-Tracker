<?php

namespace App\Jobs;

use App\Events\SpeedtestCompleteEvent;
use App\Events\SpeedtestFailedEvent;
use App\Helpers\SpeedtestHelper;
use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Bus\Dispatchable;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;
use Illuminate\Support\Facades\Log;

class SpeedtestJob implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

    private $scheduled;

    /**
     * Create a new job instance.
     *
     * @return void
     */
    public function __construct($scheduled = true)
    {
        $this->scheduled = $scheduled;
    }

    /**
     * Runs a speedtest
     *
     * @return \App\Speedtest
     */
    public function handle()
    {
        $output = SpeedtestHelper::output();
        $speedtest = SpeedtestHelper::runSpeedtest($output, $this->scheduled);
        Log::info($speedtest);
        if($speedtest == false) {
            Log::info('speedtest == false');
            event(new SpeedtestFailedEvent());
        } else {
            event(new SpeedtestCompleteEvent($speedtest));
        }
        return $speedtest;
    }
}
