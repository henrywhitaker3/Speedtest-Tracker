<?php

namespace App\Notifications;

use App\Helpers\SettingsHelper;
use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Notifications\Messages\MailMessage;
use Illuminate\Notifications\Notification;
use NotificationChannels\Telegram\TelegramChannel;
use NotificationChannels\Telegram\TelegramMessage;

class SpeedtestFailedTelegram extends Notification implements ShouldQueue
{
    use Queueable;

    /**
     * Create a new notification instance.
     *
     * @return void
     */
    public function __construct()
    {
        //
    }

    public function via($notifiable)
    {
        return [
            TelegramChannel::class
        ];
    }

    /**
     * Format tekegram notification
     *
     * @param   mixed   $notifiable
     * @return  TelegramMessage
     */
    public function toTelegram($notifiable)
    {
        $msg = "Error: something went wrong running your speedtest";
        return TelegramMessage::create()
                              ->to(SettingsHelper::get('telegram_chat_id')->value)
                              ->content($msg)
                              ->options(['parse_mode' => 'Markdown']);
    }
}
