<?php

namespace App\Providers;

use App\Helpers\UpdateHelper;
use Illuminate\Support\ServiceProvider;

class UpdaterServiceProvider extends ServiceProvider
{
    /**
     * Register services.
     *
     * @return void
     */
    public function register()
    {
        $this->app->bind('updater', function() {
            return new UpdateHelper();
        });
    }

    /**
     * Bootstrap services.
     *
     * @return void
     */
    public function boot()
    {
        //
    }
}
