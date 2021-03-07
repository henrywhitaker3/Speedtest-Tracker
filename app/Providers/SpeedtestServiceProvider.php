<?php

namespace App\Providers;

use App\Helpers\SettingsHelper;
use App\Interfaces\SpeedtestProvider;
use App\Utils\OoklaTester;
use File;
use Illuminate\Support\ServiceProvider;
use Schema;

class SpeedtestServiceProvider extends ServiceProvider
{
    /**
     * Bootstrap services.
     *
     * @return void
     */
    public function boot()
    {
        if (File::exists(env('DB_DATABASE'))) {
            if (Schema::hasTable('settings')) {
                switch (SettingsHelper::get('speedtest_provider')) {
                    case 'ookla':
                    default:
                        $this->app->singleton(
                            SpeedtestProvider::class,
                            function () {
                                return new OoklaTester();
                            }
                        );
                        break;
                }
            }
        }
    }
}
