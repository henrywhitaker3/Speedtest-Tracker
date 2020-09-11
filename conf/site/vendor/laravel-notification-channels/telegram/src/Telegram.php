<?php

namespace NotificationChannels\Telegram;

use Exception;
use GuzzleHttp\Client as HttpClient;
use GuzzleHttp\Exception\ClientException;
use Illuminate\Support\Str;
use NotificationChannels\Telegram\Exceptions\CouldNotSendNotification;
use Psr\Http\Message\ResponseInterface;

/**
 * Class Telegram.
 */
class Telegram
{
    /** @var HttpClient HTTP Client */
    protected $http;

    /** @var string|null Telegram Bot API Token. */
    protected $token;

    /** @var string Telegram Bot API Base URI */
    protected $apiBaseUri;

    /**
     * @param string|null     $token
     * @param HttpClient|null $httpClient
     * @param string|null     $apiBaseUri
     */
    public function __construct($token = null, HttpClient $httpClient = null, $apiBaseUri = null)
    {
        $this->token = $token;
        $this->http = $httpClient ?? new HttpClient();
        $this->setApiBaseUri($apiBaseUri ?? 'https://api.telegram.org');
    }

    /**
     * Token getter.
     *
     * @return string
     */
    public function getToken(): string
    {
        return $this->token;
    }

    /**
     * Token setter.
     *
     * @param string $token
     *
     * @return $this
     */
    public function setToken(string $token): self
    {
        $this->token = $token;

        return $this;
    }

    /**
     * API Base URI getter.
     *
     * @return string
     */
    public function getApiBaseUri(): string
    {
        return $this->apiBaseUri;
    }

    /**
     * API Base URI setter.
     *
     * @param string $apiBaseUri
     *
     * @return $this
     */
    public function setApiBaseUri(string $apiBaseUri): self
    {
        $this->apiBaseUri = rtrim($apiBaseUri, '/');

        return $this;
    }

    /**
     * Get HttpClient.
     *
     * @return HttpClient
     */
    protected function httpClient(): HttpClient
    {
        return $this->http;
    }

    /**
     * Set HTTP Client.
     *
     * @param HttpClient $http
     *
     * @return $this
     */
    public function setHttpClient(HttpClient $http): self
    {
        $this->http = $http;

        return $this;
    }

    /**
     * Send text message.
     *
     * <code>
     * $params = [
     *   'chat_id'                  => '',
     *   'text'                     => '',
     *   'parse_mode'               => '',
     *   'disable_web_page_preview' => '',
     *   'disable_notification'     => '',
     *   'reply_to_message_id'      => '',
     *   'reply_markup'             => '',
     * ];
     * </code>
     *
     * @link https://core.telegram.org/bots/api#sendmessage
     *
     * @param array $params
     *
     * @throws CouldNotSendNotification
     *
     * @return ResponseInterface|null
     */
    public function sendMessage(array $params): ?ResponseInterface
    {
        return $this->sendRequest('sendMessage', $params);
    }

    /**
     * Send File as Image or Document.
     *
     * @param array  $params
     * @param string $type
     * @param bool   $multipart
     *
     * @throws CouldNotSendNotification
     *
     * @return ResponseInterface|null
     */
    public function sendFile(array $params, string $type, bool $multipart = false): ?ResponseInterface
    {
        return $this->sendRequest('send'.Str::studly($type), $params, $multipart);
    }

    /**
     * Send a Location.
     *
     * @param array $params
     *
     * @throws CouldNotSendNotification
     *
     * @return ResponseInterface|null
     */
    public function sendLocation(array $params): ?ResponseInterface
    {
        return $this->sendRequest('sendLocation', $params);
    }

    /**
     * Send an API request and return response.
     *
     * @param string $endpoint
     * @param array  $params
     * @param bool   $multipart
     *
     * @throws CouldNotSendNotification
     *
     * @return ResponseInterface|null
     */
    protected function sendRequest(string $endpoint, array $params, bool $multipart = false): ?ResponseInterface
    {
        if (blank($this->token)) {
            throw CouldNotSendNotification::telegramBotTokenNotProvided('You must provide your telegram bot token to make any API requests.');
        }

        $apiUri = sprintf('%s/bot%s/%s', $this->apiBaseUri, $this->token, $endpoint);

        try {
            return $this->httpClient()->post($apiUri, [
                $multipart ? 'multipart' : 'form_params' => $params,
            ]);
        } catch (ClientException $exception) {
            throw CouldNotSendNotification::telegramRespondedWithAnError($exception);
        } catch (Exception $exception) {
            throw CouldNotSendNotification::couldNotCommunicateWithTelegram($exception);
        }
    }
}
