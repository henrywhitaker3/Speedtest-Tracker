<?php

namespace App\Http\Controllers;

use App\Helpers\BackupHelper;
use DateTime;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Storage;
use Illuminate\Support\Facades\Validator;

class BackupController extends Controller
{
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

    public function restore(Request $request)
    {
        $rule = [
            'data' => [ 'required', 'array' ],
        ];

        $validator = Validator::make($request->all(), $rule);

        if($validator->fails()) {
            return response()->json([
                'error' => $validator->errors()
            ], 422);
        }

        BackupHelper::restore($request->data);

        return response()->json([
            'method' => 'restore data from backup',
        ], 200);
    }
}
