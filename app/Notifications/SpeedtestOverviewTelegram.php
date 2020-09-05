<?php

namespace App\Notifications;

use App\Helpers\SettingsHelper;
use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Notifications\Messages\MailMessage;
use Illuminate\Notifications\Notification;
use NotificationChannels\Telegram\TelegramChannel;
use NotificationChannels\Telegram\TelegramMessage;

class SpeedtestOverviewTelegram extends Notification
{
    use Queueable;

    public $data;

    /**
     * Create a new notification instance.
     *
     * @return void
     */
    public function __construct($data)
    {
        $data['ping'] = number_format((float)$data['ping'], 1, '.', '');
        $data['download'] = number_format((float)$data['download'], 1, '.', '');
        $data['upload'] = number_format((float)$data['upload'], 1, '.', '');
        $this->data = $data;
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
     * Format tekegram notification
     *
     * @param   mixed   $notifiable
     * @return  TelegramMessage
     */
    public function toTelegram($notifiable)
    {
        $data = $this->data;
        $msg = "*Speedtest Daily Overview*
Average ping: *".$data["ping"]."*
Average download: *".$data["download"]."*
Average upload: *".$data["upload"]."*";
        return TelegramMessage::create()
                              ->to(SettingsHelper::get('telegram_chat_id')->value)
                              ->content($msg)
                              ->options(['parse_mode' => 'Markdown']);
    }

    /**
     * Get the array representation of the notification.
     *
     * @param  mixed  $notifiable
     * @return array
     */
    public function toArray($notifiable)
    {
        return [
            //
        ];
    }
}
