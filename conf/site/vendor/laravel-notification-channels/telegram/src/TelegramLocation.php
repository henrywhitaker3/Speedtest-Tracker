<?php

namespace NotificationChannels\Telegram;

use JsonSerializable;
use NotificationChannels\Telegram\Traits\HasSharedLogic;

/**
 * Class TelegramLocation.
 */
class TelegramLocation implements JsonSerializable
{
    use HasSharedLogic;

    /**
     * @param float|string|null $latitude
     * @param float|string|null $longitude
     *
     * @return static
     */
    public static function create($latitude = null, $longitude = null): self
    {
        return new static($latitude, $longitude);
    }

    /**
     * Message constructor.
     *
     * @param float|string|null $latitude
     * @param float|string|null $longitude
     */
    public function __construct($latitude = null, $longitude = null)
    {
        $this->latitude($latitude);
        $this->longitude($longitude);
    }

    /**
     * Location's latitude.
     *
     * @param float|string $latitude
     *
     * @return $this
     */
    public function latitude($latitude): self
    {
        $this->payload['latitude'] = $latitude;

        return $this;
    }

    /**
     * Location's latitude.
     *
     * @param float|string $longitude
     *
     * @return $this
     */
    public function longitude($longitude): self
    {
        $this->payload['longitude'] = $longitude;

        return $this;
    }
}
