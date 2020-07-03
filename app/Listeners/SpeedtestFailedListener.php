<?php

namespace App\Listeners;

use App\Helpers\SettingsHelper;
use App\Notifications\SpeedtestFailedSlack;
use App\Notifications\SpeedtestFailedTelegram;
use Exception;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Notification;

class SpeedtestFailedListener
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
     * Handle the event.
     *
     * @param  object  $event
     * @return void
     */
    public function handle($event)
    {
        if(env('SLACK_WEBHOOK')) {
            try {
                Notification::route('slack', env('SLACK_WEBHOOK'))
                            ->notify(new SpeedtestFailedSlack());
            } catch(Exception $e) {
                Log::notice('Your sleck webhook is invalid');
                Log::notice($e);
            }
        }

        if(env('TELEGRAM_BOT_TOKEN') && env('TELEGRAM_CHAT_ID')) {
            try {
                Notification::route(TelegramChannel::class, env('TELEGRAM_CHAT_ID'))
                            ->notify(new SpeedtestFailedTelegram());
            } catch(Exception $e) {
                Log::notice('Your telegram settings are invalid');
                Log::notice($e);
            }
        }
    }
}
