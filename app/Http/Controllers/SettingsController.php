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
     * Bulk store/update a setting
     *
     * @param   Request $request
     * @return  Response
     */
    public function bulkStore(Request $request)
    {
        $rule = [
            'data' => [ 'array', 'required' ],
            'data.*.name' => [ 'string', 'required' ],
            'data.*.value' => [ 'required' ],
        ];

        $validator = Validator::make($request->all(), $rule);
        if($validator->fails()) {
            return response()->json([
                'method' => 'Bulk store a setting',
                'error' => $validator->errors()
            ], 422);
        }

        $settings = [];
        foreach($request->data as $d) {
            if($d['name'] == 'speedtest_overview_time') {
                $ok = [ '00', '01', '02', '03', '04', '05', '06', '07', '08', '09', '10', '11', '12', '13', '14', '15', '16', '17', '18', '19', '20', '21', '22', '23' ];
                if(!in_array($d['value'], $ok)) {
                    return response()->json([
                        'method' => 'Bulk store a setting',
                        'error' => 'Invalid speedtest_overview_time value'
                    ], 422);
                }
            }
            $setting = SettingsHelper::set($d['name'], $d['value']);
            array_push($settings, $setting);
        }

        return response()->json([
            'method' => 'Bulk store a setting',
            'data' => $settings,
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
            'base' => SettingsHelper::getBase(),
            'download_upload_graph_enabled' => SettingsHelper::get('download_upload_graph_enabled'),
            'download_upload_graph_width' => SettingsHelper::get('download_upload_graph_width'),
            'ping_graph_enabled' => SettingsHelper::get('ping_graph_enabled'),
            'ping_graph_width' => SettingsHelper::get('ping_graph_width'),
            'failure_graph_enabled' => SettingsHelper::get('failure_graph_enabled'),
            'failure_graph_width' => SettingsHelper::get('failure_graph_width'),
        ];

        return $config;
    }
}
