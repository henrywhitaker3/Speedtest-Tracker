<?php

namespace App\Listeners;

use App\Helpers\SettingsHelper;
use App\Notifications\SpeedtestCompleteSlack;
use App\Notifications\SpeedtestCompleteTelegram;
use Exception;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Notification;
use NotificationChannels\Telegram\TelegramChannel;

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
        if(SettingsHelper::get('speedtest_notifications')->value == true) {
            $data = $event->speedtest;
            if(env('SLACK_WEBHOOK')) {
                try {
                    Notification::route('slack', env('SLACK_WEBHOOK'))
                                ->notify(new SpeedtestCompleteSlack($data));
                } catch(Exception $e) {
                    Log::notice('Your sleck webhook is invalid');
                    Log::notice($e);
                }
            }

            if(env('TELEGRAM_BOT_TOKEN') && env('TELEGRAM_CHAT_ID')) {
                try {
                    Notification::route(TelegramChannel::class, env('TELEGRAM_CHAT_ID'))
                                ->notify(new SpeedtestCompleteTelegram($data));
                } catch(Exception $e) {
                    Log::notice('Your telegram settings are invalid');
                    Log::notice($e);
                }
            }
        }
    }
}
