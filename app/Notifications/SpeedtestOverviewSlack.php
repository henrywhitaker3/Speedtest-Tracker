<?php

namespace App\Notifications;

use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Notifications\Messages\MailMessage;
use Illuminate\Notifications\Messages\SlackMessage;
use Illuminate\Notifications\Notification;

class SpeedtestOverviewSlack extends Notification
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
            'slack'
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
        $data = $this->data;
        return (new SlackMessage)
                ->warning()
                ->attachment(function ($attachment) use ($data) {
                    $attachment->title('Speedtest Daily Overview')
                               ->fields([
                                    'Average ping' => $data['ping'] . ' ms',
                                    'Average download' => $data['download'] . ' Mbit/s',
                                    'Average upload' => $data['upload'] . ' Mbit/s',
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
