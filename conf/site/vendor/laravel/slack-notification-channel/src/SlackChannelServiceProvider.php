<?php

namespace Illuminate\Notifications;

use GuzzleHttp\Client as HttpClient;
use Illuminate\Support\ServiceProvider;
use Illuminate\Support\Facades\Notification;

class SlackChannelServiceProvider extends ServiceProvider
{
    /**
     * Register the service provider.
     *
     * @return void
     */
    public function register()
    {
        Notification::resolved(function (ChannelManager $service) {
            $service->extend('slack', function ($app) {
                return new Channels\SlackWebhookChannel($app->make(HttpClient::class));
            });
        });
    }
}
