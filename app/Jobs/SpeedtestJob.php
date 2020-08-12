<?php

namespace App\Jobs;

use App\Events\SpeedtestCompleteEvent;
use App\Events\SpeedtestFailedEvent;
use App\Helpers\SettingsHelper;
use App\Helpers\SpeedtestHelper;
use Exception;
use Henrywhitaker3\Healthchecks\Healthchecks;
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
        $healthchecksEnabled = (bool)SettingsHelper::get('healthchecks_enabled')->value;
        $healthchecksUuid = SettingsHelper::get('healthchecks_uuid')->value;

        if($healthchecksEnabled === true) {
            try {
                $hc = new Healthchecks($healthchecksUuid);
                $hc->start();
            } catch(Exception $e) {
                Log::error($e->getMessage());
            }
        }
        $output = SpeedtestHelper::output();
        $speedtest = SpeedtestHelper::runSpeedtest($output, $this->scheduled);
        if($speedtest == false) {
            if(isset($hc)) {
                try {
                    $hc->fail();
                } catch(Exception $e) {
                    //
                }
            }
            event(new SpeedtestFailedEvent());
        } else {
            if(isset($hc)) {
                try {
                    $hc->success();
                } catch(Exception $e) {
                    //
                }
            }
            event(new SpeedtestCompleteEvent($speedtest));
        }
        return $speedtest;
    }
}
