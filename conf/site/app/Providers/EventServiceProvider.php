<?php

namespace App\Providers;

use App\Events\SpeedtestCompleteEvent;
use App\Events\SpeedtestFailedEvent;
use App\Events\SpeedtestOverviewEvent;
use App\Events\TestNotificationEvent;
use App\Listeners\SpeedtestCompleteListener;
use App\Listeners\SpeedtestFailedListener;
use App\Listeners\SpeedtestOverviewListener;
use App\Listeners\TestNotificationListener;
use Illuminate\Auth\Events\Registered;
use Illuminate\Auth\Listeners\SendEmailVerificationNotification;
use Illuminate\Foundation\Support\Providers\EventServiceProvider as ServiceProvider;
use Illuminate\Support\Facades\Event;

class EventServiceProvider extends ServiceProvider
{
    /**
     * The event listener mappings for the application.
     *
     * @var array
     */
    protected $listen = [
        Registered::class => [
            SendEmailVerificationNotification::class,
        ],
        SpeedtestCompleteEvent::class => [
            SpeedtestCompleteListener::class,
        ],
        SpeedtestOverviewEvent::class => [
            SpeedtestOverviewListener::class
        ],
        SpeedtestFailedEvent::class => [
            SpeedtestFailedListener::class
        ],
        TestNotificationEvent::class => [
            TestNotificationListener::class
        ]
    ];

    /**
     * Register any events for your application.
     *
     * @return void
     */
    public function boot()
    {
        parent::boot();

        //
    }
}
