<?php

namespace App\Helpers;

use App\Events\TestNotificationEvent;
use App\Setting;
use Cache;
use Carbon\Carbon;

class SettingsHelper
{

    /**
     * Get a Setting object by name
     *
     * @param   String                  $name   The name field in the setting table
     * @return  \App\Setting|bool|array    $name   The Setting object. Returns false if no mathcing obj.
     */
    public static function get(String $name)
    {
        $name = Setting::where('name', $name)->get();

        if (sizeof($name) == 0) {
            return false;
        } else if (sizeof($name) == 1) {
            return $name[0];
        } else {
            $name = $name->keyBy('name');
            return $name->all();
        }
    }

    /**
     * Create / update value for Setting object.
     *
     * @param   String  $name   Name of setting
     * @param   String|bool  $value  Value of setting
     * @return  \App\Setting
     */
    public static function set(String $name, $value)
    {
        $setting = SettingsHelper::get($name);

        if ($value === false) {
            $value = "0";
        }

        if ($setting !== false) {
            $setting->value = $value;
            $setting->save();
        } else {
            $setting = Setting::create([
                'name' => $name,
                'value' => $value,
            ]);
        }

        if ($name == 'show_failed_tests_on_graph') {
            Cache::flush();
        }

        return $setting;
    }

    /**
     * Get the app's base path
     *
     * @return string
     */
    public static function getBase()
    {
        $base = env('BASE_PATH', '/');
        if ($base == '') {
            $base = '/';
        } else {
            if ($base[0] != '/') {
                $base = '/' . $base;
            }
            if ($base[-1] != '/') {
                $base = $base . '/';
            }
        }
        return $base;
    }

    /**
     * Check whether a setting is defined in ENV vars or through DB
     *
     * @param string $key
     * @return boolean
     */
    public static function settingIsEditable(string $key)
    {
        $results = [];

        // Try exact key
        $val = exec('echo $' . $key);

        if ($val == "") {
            array_push($results, true);
        } else {
            array_push($results, false);
        }

        // Try key all caps
        $val = exec('echo $' . strtoupper($key));

        if ($val == "") {
            array_push($results, true);
        } else {
            array_push($results, false);
        }

        if (env($key, false) == false) {
            array_push($results, true);
        } else {
            array_push($results, false);
        }

        if (env(strtoupper($key), false) == false) {
            array_push($results, true);
        } else {
            array_push($results, false);
        }

        if (in_array(false, $results)) {
            return false;
        }

        return true;
    }

    /**
     * Get the application config
     *
     * @return array
     */
    public static function getConfig()
    {
        return [
            'base' => SettingsHelper::getBase(),
            'widgets' => [
                'show_average' => (bool)SettingsHelper::get('show_average')->value,
                'show_max' => (bool)SettingsHelper::get('show_max')->value,
                'show_min' => (bool)SettingsHelper::get('show_min')->value,
            ],
            'graphs' => [
                'download_upload_graph_enabled' => SettingsHelper::get('download_upload_graph_enabled'),
                'download_upload_graph_width' => SettingsHelper::get('download_upload_graph_width'),
                'ping_graph_enabled' => SettingsHelper::get('ping_graph_enabled'),
                'ping_graph_width' => SettingsHelper::get('ping_graph_width'),
                'failure_graph_enabled' => SettingsHelper::get('failure_graph_enabled'),
                'failure_graph_width' => SettingsHelper::get('failure_graph_width'),
            ],
            'editable' => [
                'slack_webhook' => SettingsHelper::settingIsEditable('slack_webhook'),
                'telegram_bot_token' => SettingsHelper::settingIsEditable('telegram_bot_token'),
                'telegram_chat_id' => SettingsHelper::settingIsEditable('telegram_chat_id'),
            ],
            'auth' => (bool)SettingsHelper::get('auth')->value
        ];
    }

    /**
     * Send test notification to agents
     *
     * @param boolean|string $agent
     * @return bool
     */
    public static function testNotification($agent = true)
    {
        $agents = ['slack', 'telegram'];

        if ($agent === true) {
            event(new TestNotificationEvent($agents));
            return true;
        }

        if (in_array($agent, $agents)) {
            event(new TestNotificationEvent([$agent]));
            return true;
        }

        return false;
    }

    public static function loadIntegrationConfig()
    {
        $settings = [
            'healthchecks_enabled' => (bool)SettingsHelper::get('healthchecks_enabled')->value,
            'healthchecks_uuid' => SettingsHelper::get('healthchecks_uuid')->value,
            'slack_webhook' => SettingsHelper::get('slack_webhook')->value,
            'telegram_bot_token' => SettingsHelper::get('telegram_bot_token')->value,
            'telegram_chat_id' => SettingsHelper::get('telegram_chat_id')->value,
        ];

        foreach ($settings as $key => $value) {
            $key = 'integrations.' . $key;

            if ($value === "") {
                $value = null;
            }

            config()->set([$key => $value]);
        }
    }
}
