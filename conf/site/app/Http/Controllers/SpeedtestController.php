<?php

namespace App\Http\Controllers;

use App\Helpers\SpeedtestHelper;
use App\Jobs\SpeedtestJob;
use App\Speedtest;
use Carbon\Carbon;
use Exception;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Validator;

class SpeedtestController extends Controller
{
    public function index()
    {
        $data = Speedtest::orderBy('id', 'desc')
                         ->paginate();

        return response()->json([
            'method' => 'index of speedtests',
            'data' => $data,
        ], 200);
    }

    public function time($days)
    {
        $rule = [
            'days' => [ 'required', 'integer' ],
        ];

        $validator = Validator::make([ 'days' => $days ], $rule);

        if($validator->fails()) {
            return response()->json([
                'method' => 'get speedtests in last x days',
                'error' => $validator->errors(),
            ], 422);
        }

        $data = Speedtest::where('created_at', '>=', Carbon::now()->subDays($days))
                         ->get();

        return response()->json([
            'method' => 'get speedtests in last x days',
            'days' => $days,
            'data' => $data
        ], 200);
    }

    public function latest()
    {
        $data = SpeedtestHelper::latest();
        $avg = Speedtest::select(DB::raw('AVG(ping) as ping, AVG(download) as download, AVG(upload) as upload'))
                        ->get();
        $max = Speedtest::select(DB::raw('MAX(ping) as ping, MAX(download) as download, MAX(upload) as upload'))
                        ->get();

        if($data) {
            return response()->json([
                'method' => 'get latest speedtest',
                'data' => $data,
                'average' => $avg[0],
                'max' => $max[0],
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
