<?php

namespace Henrywhitaker3\LaravelActions;

use Henrywhitaker3\LaravelActions\Commands\MakeActionCommand;
use Illuminate\Support\ServiceProvider;

class LaravelActionsServiceProvider extends ServiceProvider
{
    protected $commands = [
        MakeActionCommand::class,
    ];

    /**
     * Perform post-registration booting of services.
     *
     * @return void
     */
    public function boot(): void
    {
        if ($this->app->runningInConsole()) {
            $this->bootForConsole();

            $this->publishes([
                __DIR__.'/Stubs/ActionStub.php' => resource_path('stubs/ActionStub.php'),
            ], 'stubs');
        }

        require_once __DIR__.'/Helpers/ActionHelper.php';
    }

    /**
     * Get the services provided by the provider.
     *
     * @return array
     */
    public function provides()
    {
        return [];
    }

    /**
     * Console-specific booting.
     *
     * @return void
     */
    protected function bootForConsole(): void
    {
        // Registering package commands.
        $this->commands($this->commands);
    }
}
