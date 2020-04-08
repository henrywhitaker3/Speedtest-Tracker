<?php

namespace App\Http\Controllers;

use App\Helpers\SpeedtestHelper;
use Exception;
use Illuminate\Http\Request;

class SpeedtestController extends Controller
{
    public function latest()
    {
        $data = SpeedtestHelper::latest();

        if($data) {
            return response()->json([
                'method' => 'get latest speedtest',
                'data' => $data
            ], 200);
        } else {
            return response()->json([
                'method' => 'get latest speedtest',
                'error' => 'no speedtests have been run'
            ], 404);
        }
    }

    public function run()
    {
        try {
            $data = SpeedtestHelper::runSpeedtest();
            return response()->json([
                'method' => 'run speedtest',
                'data' => $data
            ], 200);
        } catch(Exception $e) {
            return response()->json([
                'method' => 'run speedtest',
                'error' => $e
            ], 500);
        }
    }
}
