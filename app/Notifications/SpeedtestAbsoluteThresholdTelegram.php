<?php

namespace App\Notifications;

use App\Helpers\NotificationsHelper;
use App\Helpers\SettingsHelper;
use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Notifications\Messages\MailMessage;
use Illuminate\Notifications\Notification;
use NotificationChannels\Telegram\TelegramChannel;
use NotificationChannels\Telegram\TelegramMessage;

class SpeedtestAbsoluteThresholdTelegram extends Notification
{
    use Queueable;

    private $errors;

    /**
     * Create a new notification instance.
     *
     * @return void
     */
    public function __construct($errors)
    {
        $this->errors = $errors;
    }

    /**
     * Get the notification's delivery channels.
     *
     * @param  mixed  $notifiable
     * @return array
     */
    public function via($notifiable)
    {
        return [
            TelegramChannel::class
        ];
    }

    /**
     * Format telegram notification
     *
     * @param   mixed   $notifiable
     * @return  TelegramMessage
     */
    public function toTelegram($notifiable)
    {
        $msg = NotificationsHelper::formatAbsoluteThresholdMessage($this->errors);

        return TelegramMessage::create()
                              ->to(SettingsHelper::get('telegram_chat_id')->value)
                              ->content($msg)
                              ->options(['parse_mode' => 'Markdown']);
    }
}
