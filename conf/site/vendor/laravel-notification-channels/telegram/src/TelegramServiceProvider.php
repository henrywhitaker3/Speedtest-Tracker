<?php

namespace NotificationChannels\Telegram;

use GuzzleHttp\Client as HttpClient;
use Illuminate\Support\ServiceProvider;

/**
 * Class TelegramServiceProvider.
 */
class TelegramServiceProvider extends ServiceProvider
{
    /**
     * Bootstrap the application services.
     */
    public function boot(): void
    {
        $this->app->when(TelegramChannel::class)
            ->needs(Telegram::class)
            ->give(static function () {
                return new Telegram(
                    config('services.telegram-bot-api.token'),
                    app(HttpClient::class),
                    config('services.telegram-bot-api.base_uri')
                );
            });
    }
}
