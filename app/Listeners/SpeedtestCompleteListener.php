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
            if(SettingsHelper::get('slack_webhook')) {
                try {
                    Notification::route('slack', SettingsHelper::get('slack_webhook')->value)
                                ->notify(new SpeedtestCompleteSlack($data));
                } catch(Exception $e) {
                    Log::notice('Your sleck webhook is invalid');
                    Log::notice($e);
                }
            }

            if(SettingsHelper::get('telegram_bot_token')->value == true && SettingsHelper::get('telegram_chat_id')->value == true) {
                try {
                    config([ 'services.telegram-bot-api' => [ 'token' => SettingsHelper::get('telegram_bot_token')->value ] ]);
                    Notification::route(TelegramChannel::class, SettingsHelper::get('telegram_chat_id')->value)
                                ->notify(new SpeedtestCompleteTelegram($data));
                } catch(Exception $e) {
                    Log::notice('Your telegram settings are invalid');
                    Log::notice($e);
                }
            }
        }
    }
}
