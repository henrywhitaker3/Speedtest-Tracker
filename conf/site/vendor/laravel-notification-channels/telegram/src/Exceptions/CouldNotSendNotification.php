<?php

namespace NotificationChannels\Telegram\Exceptions;

use Exception;
use GuzzleHttp\Exception\ClientException;

/**
 * Class CouldNotSendNotification.
 */
class CouldNotSendNotification extends Exception
{
    /**
     * Thrown when there's a bad request and an error is responded.
     *
     * @param ClientException $exception
     *
     * @return static
     */
    public static function telegramRespondedWithAnError(ClientException $exception): self
    {
        if (! $exception->hasResponse()) {
            return new static('Telegram responded with an error but no response body found');
        }

        $statusCode = $exception->getResponse()->getStatusCode();

        $result = json_decode($exception->getResponse()->getBody(), false);
        $description = $result->description ?? 'no description given';

        return new static("Telegram responded with an error `{$statusCode} - {$description}`", 0, $exception);
    }

    /**
     * Thrown when there's no bot token provided.
     *
     * @param string $message
     *
     * @return static
     */
    public static function telegramBotTokenNotProvided($message): self
    {
        return new static($message);
    }

    /**
     * Thrown when we're unable to communicate with Telegram.
     *
     * @param $message
     *
     * @return static
     */
    public static function couldNotCommunicateWithTelegram($message): self
    {
        return new static("The communication with Telegram failed. `{$message}`");
    }
}
