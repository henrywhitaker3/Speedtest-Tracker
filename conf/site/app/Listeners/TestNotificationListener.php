<?php

namespace App\Listeners;

use App\Helpers\SettingsHelper;
use App\Notifications\TestSlackNotification;
use App\Notifications\TestTelegramNotification;
use Exception;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Notification;
use NotificationChannels\Telegram\TelegramChannel;

class TestNotificationListener
{
    private $agents;

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
        foreach($event->agents as $agent) {
            if($agent == 'slack') {
                $this->slackNotification();
            }

            if($agent == 'telegram') {
                $this->telegramNotification();
            }
        }
    }

    /**
     * Send a slack notification
     *
     * @return void
     */
    private function slackNotification()
    {
        if(SettingsHelper::get('slack_webhook')->value == true) {
            try {
                Notification::route('slack', SettingsHelper::get('slack_webhook')->value)
                            ->notify(new TestSlackNotification());
            } catch(Exception $e) {
                Log::notice('Your sleck webhook is invalid');
                Log::notice($e);
            }
        }
    }

    /**
     * Send a telegram notification
     *
     * @return void
     */
    private function telegramNotification()
    {
        if(SettingsHelper::get('telegram_bot_token')->value == true && SettingsHelper::get('telegram_chat_id')->value == true) {
            try {
                config([ 'services.telegram-bot-api' => [ 'token' => SettingsHelper::get('telegram_bot_token')->value ] ]);
                Notification::route(TelegramChannel::class, SettingsHelper::get('telegram_bot_token')->value)
                            ->notify(new TestTelegramNotification());
            } catch(Exception $e) {
                Log::notice('Your telegram settings are invalid');
                Log::notice($e);
            }
        }
    }
}
