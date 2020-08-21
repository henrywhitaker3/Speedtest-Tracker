<?php

namespace App\Notifications;

use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Notifications\Messages\MailMessage;
use Illuminate\Notifications\Messages\SlackMessage;
use Illuminate\Notifications\Notification;

class SpeedtestAbsoluteThresholdNotificationSlack extends Notification
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
        return (new SlackMessage)
                ->warning()
                ->attachment(function ($attachment) {
                    $attachment->title('Speedtest absolute threshold error')
                               ->content($this->formatMessage());
                });
    }

    /**
     * Parse $this->errors and format message
     *
     * @return String
     */
    public function formatMessage()
    {
        $msg = 'For the latest speedtest, the ';

        for($i = 0; $i < sizeof($this->errors); $i++) {
            $key = $this->errors[$i];
            $msg = $msg . $key;
            if(sizeof($this->errors) > 1 && $i < (sizeof($this->errors) - 1)) {
                $msg = $msg . ', ';
            }
        }

        if($msg[-1] != '') {
            $msg = $msg . ' ';
        }

        if(sizeof($this->errors) > 1) {
            $msg = $msg . 'values ';
        } else {
            $msg = $msg . 'value ';
        }

        $msg = $msg . 'exceeded the absolute threshold';

        return $msg;
    }
}
