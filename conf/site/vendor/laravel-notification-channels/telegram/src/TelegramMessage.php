<?php

namespace NotificationChannels\Telegram;

use Illuminate\Support\Facades\View;
use JsonSerializable;
use NotificationChannels\Telegram\Traits\HasSharedLogic;

/**
 * Class TelegramMessage.
 */
class TelegramMessage implements JsonSerializable
{
    use HasSharedLogic;

    /**
     * @param string $content
     *
     * @return self
     */
    public static function create(string $content = ''): self
    {
        return new self($content);
    }

    /**
     * Message constructor.
     *
     * @param string $content
     */
    public function __construct(string $content = '')
    {
        $this->content($content);
        $this->payload['parse_mode'] = 'Markdown';
    }

    /**
     * Notification message (Supports Markdown).
     *
     * @param string $content
     *
     * @return $this
     */
    public function content(string $content): self
    {
        $this->payload['text'] = $content;

        return $this;
    }

    /**
     * Attach a view file as the content for the notification.
     * Supports Laravel blade template.
     *
     * @param string $view
     * @param array  $data
     * @param array  $mergeData
     *
     * @return $this
     */
    public function view(string $view, array $data = [], array $mergeData = []): self
    {
        return $this->content(View::make($view, $data, $mergeData)->render());
    }
}
