<?php

namespace NotificationChannels\Telegram\Traits;

/**
 * Trait HasSharedLogic.
 */
trait HasSharedLogic
{
    /** @var array Params payload. */
    protected $payload = [];

    /** @var array Inline Keyboard Buttons. */
    protected $buttons = [];

    /** @var string Bot Token. */
    public $token;

    /**
     * Recipient's Chat ID.
     *
     * @param string|int $chatId
     *
     * @return $this
     */
    public function to($chatId): self
    {
        $this->payload['chat_id'] = $chatId;

        return $this;
    }

    /**
     * Add an inline button.
     *
     * @param string $text
     * @param string $url
     * @param int    $columns
     *
     * @return $this
     */
    public function button(string $text, string $url, int $columns = 2): self
    {
        $this->buttons[] = compact('text', 'url');

        $this->payload['reply_markup'] = json_encode([
            'inline_keyboard' => array_chunk($this->buttons, $columns),
        ]);

        return $this;
    }

    /**
     * Send the message silently.
     * Users will receive a notification with no sound.
     *
     * @param bool $disableNotification
     *
     * @return $this
     */
    public function disableNotification(bool $disableNotification = true): self
    {
        $this->payload['disable_notification'] = $disableNotification;

        return $this;
    }

    /**
     * Bot Token.
     * Overrides default bot token with the given value for this notification.
     *
     * @param string $token
     *
     * @return $this
     */
    public function token(string $token): self
    {
        $this->token = $token;

        return $this;
    }

    /**
     * Determine if bot token is given for this notification.
     *
     * @return bool
     */
    public function hasToken(): bool
    {
        return null !== $this->token;
    }

    /**
     * Additional options to pass to sendMessage method.
     *
     * @param array $options
     *
     * @return $this
     */
    public function options(array $options): self
    {
        $this->payload = array_merge($this->payload, $options);

        return $this;
    }

    /**
     * Determine if chat id is not given.
     *
     * @return bool
     */
    public function toNotGiven(): bool
    {
        return ! isset($this->payload['chat_id']);
    }

    /**
     * Get payload value for given key.
     *
     * @param string $key
     *
     * @return mixed|null
     */
    public function getPayloadValue(string $key)
    {
        return $this->payload[$key] ?? null;
    }

    /**
     * Returns params payload.
     *
     * @return array
     */
    public function toArray(): array
    {
        return $this->payload;
    }

    /**
     * Convert the object into something JSON serializable.
     *
     * @return mixed
     */
    public function jsonSerialize()
    {
        return $this->toArray();
    }
}
