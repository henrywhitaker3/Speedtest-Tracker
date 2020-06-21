<?php

namespace App\Notifications;

use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Notifications\Notification;
use Illuminate\Notifications\Messages\SlackMessage;
use Illuminate\Support\Facades\Log;
use NotificationChannels\Telegram\TelegramChannel;
use NotificationChannels\Telegram\TelegramMessage;

class SpeedtestCompleteSlack extends Notification
{
    use Queueable;

    protected $speedtest;

    /**
     * Create a new notification instance.
     *
     * @return void
     */
    public function __construct($speedtest)
    {
        $speedtest->ping = number_format((float)$speedtest->ping, 1, '.', '');
        $speedtest->download = number_format((float)$speedtest->download, 1, '.', '');
        $speedtest->upload = number_format((float)$speedtest->upload, 1, '.', '');
        $this->speedtest = $speedtest;
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
            'slack',
        ];
    }

    /**
     * Format slack notification
     *
     * @param   mixed   $notifiable
     * @return  SlackMessage
     */
    public function toSlack($notifiable)
    {
        $speedtest = $this->speedtest;
        return (new SlackMessage)
                ->warning()
                ->attachment(function ($attachment) use ($speedtest) {
                    $attachment->title('New speedtest')
                               ->fields([
                                    'Ping' => $speedtest->ping . ' ms',
                                    'Download' => $speedtest->download . ' Mbit/s',
                                    'Upload' => $speedtest->upload . ' Mbit/s',
                                ]);
                });
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
