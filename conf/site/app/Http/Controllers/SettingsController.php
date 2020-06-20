<?php

namespace App\Http\Controllers;

use App\Helpers\SettingsHelper;
use App\Rules\Cron;
use App\Setting;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class SettingsController extends Controller
{

    /**
     * Return all settings
     *
     * @return  array
     */
    public function index()
    {
        return Setting::get()->keyBy('name');
    }

    /**
     * Get setting by id
     *
     * @param   Setting $setting
     * @return  Setting
     */
    public function get(Setting $setting)
    {
        return $setting;
    }

    /**
     * Store/update a setting
     *
     * @param   Request $request
     * @return  Response
     */
    public function store(Request $request)
    {
        $rule = [
            'name' => [ 'required', 'string', 'min:1' ],
        ];
        if($request->name == 'schedule') {
            $rule['value'] = [ 'required', new Cron ];
        }

        $validator = Validator::make($request->all(), $rule);
        if($validator->fails()) {
            return response()->json([
                'method' => 'Store a setting',
                'error' => $validator->errors()
            ], 422);
        }

        if(!isset($request->value)) {
            $request->value = '';
        }

        $setting = SettingsHelper::set($request->name, $request->value);

        return response()->json([
            'method' => 'Store a setting',
            'data' => $setting
        ], 200);
    }

    /**
     * Returns instance config
     *
     * @return  array
     */
    public function config()
    {


        $config = [
            'base' => SettingsHelper::getBase()
        ];

        return $config;
    }
}
