<?php

namespace App\Http\Controllers;

use App\Helpers\BackupHelper;
use App\Helpers\SettingsHelper;
use DateTime;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Storage;
use Illuminate\Support\Facades\Validator;
use Illuminate\Http\JsonResponse;

class BackupController extends Controller
{

    public function __construct()
    {
        if((bool)SettingsHelper::get('auth')->value === true) {
            $this->middleware('auth:api');
        }
    }

    /**
     * Get backup of speedtests
     *
     * @param   Request $request
     * @return  mixed|JsonResponse
     */
    public function backup(Request $request)
    {
        $validator = Validator::make($request->all(), [ 'format' => 'in:json,csv' ]);
        if($validator->fails()) {
            return response()->json([
                'method' => 'backup data',
                'error' => $validator->errors(),
            ], 422);
        }

        $filename = BackupHelper::backup($request->format);

        return Storage::disk('local')->download($filename);
    }

    /**
     * Retore from a backup
     *
     * @param   Request $request
     * @return  JsonResponse
     */
    public function restore(Request $request)
    {
        $rule = [
            'data' => [ 'required' ],
            'format' => [ 'required', 'in:json,csv' ]
        ];

        $validator = Validator::make($request->all(), $rule);

        if($validator->fails()) {
            return response()->json([
                'error' => $validator->errors()
            ], 422);
        }

        if(BackupHelper::restore($request->data, $request->format) != false) {
            return response()->json([
                'method' => 'restore data from backup',
            ], 200);
        } else {
            return response()->json([
                'method' => 'incorrect backup format',
            ], 422);
        }
    }
}
