<?php

namespace App\Jobs;

use App\Events\SpeedtestCompleteEvent;
use App\Events\SpeedtestFailedEvent;
use App\Helpers\SettingsHelper;
use App\Helpers\SpeedtestHelper;
use App\Interfaces\SpeedtestProvider;
use App\Utils\OoklaTester;
use Exception;
use Healthcheck;
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

    /**
     * Scheduled bool
     *
     * @var bool
     */
    private $scheduled;

    /**
     * Integrations config array
     *
     * @var array
     */
    private $config;

    private SpeedtestProvider $speedtestProvider;

    /**
     * Create a new job instance.
     *
     * @return void
     */
    public function __construct($scheduled = true, $config = [], SpeedtestProvider $speedtestProvider)
    {
        $this->scheduled = $scheduled;
        $this->config = $config;
        $this->speedtestProvider = $speedtestProvider;
    }

    /**
     * Runs a speedtest
     *
     * @return \App\Speedtest
     */
    public function handle()
    {
        if ($this->config['healthchecks_enabled'] === true) {
            $this->healthcheck('start');
        }
        $output = $this->speedtestProvider->output();
        $speedtest = $this->speedtestProvider->run($output, $this->scheduled);
        if ($speedtest == false) {
            if ($this->config['healthchecks_enabled'] === true) {
                $this->healthcheck('fail');
            }

            event(new SpeedtestFailedEvent());
        } else {
            if ($this->config['healthchecks_enabled'] === true) {
                $this->healthcheck('success');
            }

            event(new SpeedtestCompleteEvent($speedtest));
        }
        return $speedtest;
    }

    /**
     * Wrapper to reduce duplication of try/catch for hc
     *
     * @param String $method
     * @return void
     */
    private function healthcheck(String $method)
    {
        try {
            $hc = new Healthchecks(SettingsHelper::get('healthchecks_uuid')->value, SettingsHelper::get('healthchecks_server_url')->value);
            if ($method === 'start') {
                $hc->start();
            }

            if ($method === 'success') {
                $hc->success();
            }

            if ($method === 'fail') {
                $hc->fail();
            }
        } catch (Exception $e) {
            Log::error($e->getMessage());
        }
    }
}
