<?php

namespace App\Listeners;

use App\Notifications\SpeedtestComplete;
use Exception;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Notification;

class SpeedtestCompleteListener
{
    /**
     * Create the event listener.
     *
     * @return void
     */
    public function __construct()
    {
        //
    }

    /**
     * Handle what to do after speedtest completes
     *
     * @param  object  $event
     * @return void
     */
    public function handle($event)
    {
        if(env('SLACK_WEBHOOK')) {
            $data = $event->speedtest;
            try {
                Notification::route('slack', env('SLACK_WEBHOOK'))
                            ->notify(new SpeedtestComplete($data));
            } catch(Exception $e) {
                Log::notice('Your sleck webhook is invalid');
                Log::notice($e);
            }
        }
    }
}
