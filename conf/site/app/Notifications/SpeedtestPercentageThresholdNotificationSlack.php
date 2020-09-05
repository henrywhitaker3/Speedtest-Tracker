<?php

namespace App\Notifications;

use App\Helpers\NotificationsHelper;
use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Notifications\Messages\MailMessage;
use Illuminate\Notifications\Messages\SlackMessage;
use Illuminate\Notifications\Notification;

class SpeedtestPercentageThresholdNotificationSlack extends Notification
{
    use Queueable;

    protected $errors;

    /**
     * Create a new notification instance.
     *
     * @param array $errors
     * @return void
     */
    public function __construct(array $errors)
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
        return ['slack'];
    }

    /**
     * Format slack notification
     *
     * @param   mixed   $notifiable
     * @return  SlackMessage
     */
    public function toSlack($notifiable)
    {
        $msg = NotificationsHelper::formatPercentageThresholdMessage($this->errors);

        return (new SlackMessage)
                ->warning()
                ->attachment(function ($attachment) use ($msg) {
                    $attachment->title('Speedtest percentage threshold error')
                               ->content($msg);
                });
    }
}
