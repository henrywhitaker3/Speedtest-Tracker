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

    /**
     * Returns paginated list of speedtests
     *
     * @return  Response
     */
    public function index()
    {
        $data = Speedtest::orderBy('created_at', 'desc')
                         ->paginate();

        return response()->json([
            'method' => 'index of speedtests',
            'data' => $data,
        ], 200);
    }

    /**
     * Returns speedtest going back 'x' days
     *
     * @param   int     $days
     * @return  void
     */
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
                         ->orderBy('created_at', 'asc')
                         ->get();

        return response()->json([
            'method' => 'get speedtests in last x days',
            'days' => $days,
            'data' => $data
        ], 200);
    }

    /**
     * Return latest speedtest
     *
     * @return  Response
     */
    public function latest()
    {
        $data = SpeedtestHelper::latest();
        $avg = Speedtest::select(DB::raw('AVG(ping) as ping, AVG(download) as download, AVG(upload) as upload'))
                        ->where('failed', false)
                        ->get();
        $max = Speedtest::select(DB::raw('MAX(ping) as ping, MAX(download) as download, MAX(upload) as upload'))
                        ->where('failed', false)
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

    /**
     * Queue a new speedtest
     *
     * @return Response
     */
    public function run()
    {
        try {
            $data = SpeedtestJob::dispatch(false);
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
