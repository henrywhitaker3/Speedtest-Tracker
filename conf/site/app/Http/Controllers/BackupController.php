<?php

namespace App\Http\Controllers;

use App\Helpers\BackupHelper;
use DateTime;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Storage;
use Illuminate\Support\Facades\Validator;

class BackupController extends Controller
{
    public function backup()
    {
        $data = BackupHelper::backup();
        $timestamp = new DateTime();
        $timestamp = $timestamp->format('Y-m-d_H:i:s');
        $name = 'speedtest_backup_' . $timestamp . '.json';
        Storage::disk('local')->put($name, $data);

        return Storage::disk('local')->download($name);
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
