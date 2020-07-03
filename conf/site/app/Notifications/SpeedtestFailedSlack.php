<?php

namespace App\Notifications;

use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Notifications\Messages\MailMessage;
use Illuminate\Notifications\Messages\SlackMessage;
use Illuminate\Notifications\Notification;

class SpeedtestFailedSlack extends Notification implements ShouldQueue
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
        return (new SlackMessage)
                ->error()
                ->attachment(function ($attachment) {
                    $attachment->title('Failed speedtest')
                               ->content('Something went wrong running your speedtest');
                });
    }
}
