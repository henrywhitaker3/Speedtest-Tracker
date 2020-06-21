<?php

namespace App\Providers;

use App\Events\SpeedtestCompleteEvent;
use App\Events\SpeedtestOverviewEvent;
use App\Listeners\SpeedtestCompleteListener;
use App\Listeners\SpeedtestOverviewListener;
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
