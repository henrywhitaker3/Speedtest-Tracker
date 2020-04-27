<?php

namespace App\Notifications;

use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Notifications\Notification;
use Illuminate\Notifications\Messages\SlackMessage;

class SpeedtestComplete extends Notification
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
        return ['slack'];
    }

    public function toSlack($notifiable)
    {
        $speedtest = $this->speedtest;
        return (new SlackMessage)
                ->warning()
                ->attachment(function ($attachment) use ($speedtest) {
                    $attachment->title('New speedtest')
                               ->fields([
                                    'Ping' => number_format((float)$speedtest->ping, 1, '.', '') . ' ms',
                                    'Download' => number_format((float)$speedtest->download, 1, '.', '') . ' Mbit/s',
                                    'Upload' => number_format((float)$speedtest->upload, 1, '.', '') . ' Mbit/s',
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
