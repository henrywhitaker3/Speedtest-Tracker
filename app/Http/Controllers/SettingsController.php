<?php

namespace App\Http\Controllers;

use App\Helpers\SettingsHelper;
use App\Setting;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class SettingsController extends Controller
{
    public function index()
    {
        return Setting::get();
    }

    public function get(Setting $setting)
    {
        return $setting;
    }

    public function store(Request $request)
    {
        $rule = [
            'name' => [ 'required', 'string', 'min:1' ],
        ];
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
}
