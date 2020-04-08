<?php

namespace App\Http\Controllers;

use App\Helpers\SpeedtestHelper;
use App\Jobs\SpeedtestJob;
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
            $data = SpeedtestJob::dispatch();
            return response()->json([
                'method' => 'run speedtest',
                'data' => 'a new speedtest has been added to the queue'
            ], 200);
        } catch(Exception $e) {
            return response()->json([
                'method' => 'run speedtest',
                'error' => $e
            ], 500);
        }
    }
}
