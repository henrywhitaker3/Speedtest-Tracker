<?php

namespace App\Providers;

use App\Helpers\SettingsHelper;
use Exception;
use Henrywhitaker3\Healthchecks\Healthchecks;
use Illuminate\Support\Facades\App;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\ServiceProvider;
use Ramsey\Uuid\Exception\InvalidUuidStringException;

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
        $settings = [
            'healthchecks_enabled' => (bool)SettingsHelper::get('healthchecks_enabled')->value,
            'healthchecks_uuid' => SettingsHelper::get('healthchecks_uuid')->value,
            'slack_webhook' => SettingsHelper::get('slack_webhook')->value,
            'telegram_bot_token' => SettingsHelper::get('telegram_bot_token')->value,
            'telegram_chat_id' => SettingsHelper::get('telegram_chat_id')->value,
        ];

        foreach($settings as $key => $value) {
            $key = 'integrations.' . $key;

            if($value === "") {
                $value = null;
            }

            config()->set([ $key => $value ]);
        }

        if($settings['healthchecks_enabled']) {
            try {
                App::bind('healthcheck', function() use ($settings) {
                    return new Healthchecks($settings['healthchecks_uuid']);
                });
            } catch(InvalidUuidStringException $e) {
                Log::error('Invalid healthchecks UUID');
            } catch(Exception $e) {
                Log::error($e->getMessage());
            }
        }
    }
}
