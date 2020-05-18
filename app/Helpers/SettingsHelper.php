<?php

namespace App\Helpers;

use App\Setting;

class SettingsHelper {
    public static function get(String $name)
    {
        $name = Setting::where('name', $name)->get();

        if(sizeof($name) == 0) {
            return false;
        } else if(sizeof($name) == 1) {
            return $name[0];
        } else {
            return $name;
        }
    }

    public static function set(String $name, String $value)
    {
        $setting = SettingsHelper::get($name);

        if($setting !== false) {
            $setting->value = $value;
            $setting->save();
        } else {
            $setting = Setting::create([
                'name' => $name,
                'value' => $value,
            ]);
        }

        return $setting;
    }
}
