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
            $name = $name->keyBy('name');
            return $name->all();
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

    public static function getBase()
    {
        $base = env('BASE_PATH', '/');
        if($base == '') {
            $base = '/';
        } else {
            if($base[0] != '/') {
                $base = '/' . $base;
            }
            if($base[-1] != '/') {
                $base = $base . '/';
            }
        }
        return $base;
    }
}
