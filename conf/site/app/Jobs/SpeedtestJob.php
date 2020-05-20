<?php

namespace App\Jobs;

use App\Events\SpeedtestCompleteEvent;
use App\Helpers\SpeedtestHelper;
use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Bus\Dispatchable;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;

class SpeedtestJob implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

    /**
     * Create a new job instance.
     *
     * @return void
     */
    public function __construct()
    {
        //
    }

    /**
     * Execute the job.
     *
     * @return void
     */
    public function handle()
    {
        $output = SpeedtestHelper::output();
        $speedtest = SpeedtestHelper::runSpeedtest($output);
        event(new SpeedtestCompleteEvent($speedtest));
        return $speedtest;
    }
}
