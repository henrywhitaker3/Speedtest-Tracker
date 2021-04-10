<?php

namespace App\Actions;

use App\Helpers\SettingsHelper;
use App\Interfaces\SpeedtestProvider;
use App\Jobs\SpeedtestJob;
use Henrywhitaker3\LaravelActions\Interfaces\ActionInterface;

class QueueSpeedtest implements ActionInterface
{
    private SpeedtestProvider $speedtestProvider;

    /**
     * Create a new action instance.
     *
     * @return void
     */
    public function __construct(SpeedtestProvider $speedtestProvider)
    {
        $this->speedtestProvider = $speedtestProvider;
    }

    /**
     * Run the action.
     *
     * @return mixed
     */
    public function run()
    {
        SettingsHelper::loadIntegrationConfig();

        SpeedtestJob::dispatch(false, config('integrations'), $this->speedtestProvider);
    }
}
