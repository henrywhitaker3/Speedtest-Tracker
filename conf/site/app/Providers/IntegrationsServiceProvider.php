<?php

namespace App\Providers;

use App\Helpers\SettingsHelper;
use Exception;
use File;
use Henrywhitaker3\Healthchecks\Healthchecks;
use Illuminate\Support\Facades\App;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\ServiceProvider;
use Ramsey\Uuid\Exception\InvalidUuidStringException;
use Schema;

/**
 * This class updates the integrations.php config with the relevant values
 * from the databse.
 */
class IntegrationsServiceProvider extends ServiceProvider
{
    /**
     * Register services.
     *
     * @return void
     */
    public function register()
    {
        //
    }

    /**
     * Bootstrap services.
     *
     * @return void
     */
    public function boot()
    {
        if(File::exists(env('DB_DATABASE'))) {
            if(Schema::hasTable('settings')) {
                $setting = SettingsHelper::get('healthchecks_uuid');

                if($setting !== false) {
                    try {
                        App::bind('healthcheck', function() use ($setting) {
                            return new Healthchecks($setting->value);
                        });
                    } catch(InvalidUuidStringException $e) {
                        Log::error('Invalid healthchecks UUID');
                    } catch(Exception $e) {
                        Log::error($e->getMessage());
                    }
                }
            }
        }
    }
}
