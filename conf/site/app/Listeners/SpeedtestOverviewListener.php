<?php

namespace App\Listeners;

use App\Helpers\SettingsHelper;
use App\Helpers\SpeedtestHelper;
use App\Notifications\SpeedtestOverviewSlack;
use App\Notifications\SpeedtestOverviewTelegram;
use Exception;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Notification;
use NotificationChannels\Telegram\TelegramChannel;

class SpeedtestOverviewListener
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
        if(SettingsHelper::get('speedtest_overview_notification')->value == true) {
            $data = SpeedtestHelper::last24Hours();
            if(env('SLACK_WEBHOOK')) {
                try {
                    Notification::route('slack', env('SLACK_WEBHOOK'))
                                ->notify(new SpeedtestOverviewSlack($data));
                } catch(Exception $e) {
                    Log::notice('Your sleck webhook is invalid');
                    Log::notice($e);
                }
            }

            if(env('TELEGRAM_BOT_TOKEN') && env('TELEGRAM_CHAT_ID')) {
                try {
                    Notification::route(TelegramChannel::class, env('TELEGRAM_CHAT_ID'))
                                ->notify(new SpeedtestOverviewTelegram($data));
                } catch(Exception $e) {
                    Log::notice('Your telegram settings are invalid');
                    Log::notice($e);
                }
            }
        }
    }
}
