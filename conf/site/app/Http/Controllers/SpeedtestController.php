<?php

namespace App\Http\Controllers;

use App\Helpers\SettingsHelper;
use App\Helpers\SpeedtestHelper;
use App\Jobs\SpeedtestJob;
use App\Speedtest;
use Carbon\Carbon;
use Exception;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Validator;
use Illuminate\Http\JsonResponse;

class SpeedtestController extends Controller
{
    public function __construct()
    {
        if ((bool)SettingsHelper::get('auth')->value === true) {
            $this->middleware('auth:api')
                ->only(['run', 'delete', 'deleteAll']);
        }
    }

    /**
     * Returns paginated list of speedtests
     *
     * @return  JsonResponse
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
     * @return  JsonResponse
     */
    public function time($days)
    {
        $rule = [
            'days' => ['required', 'integer'],
        ];

        $validator = Validator::make(['days' => $days], $rule);

        if ($validator->fails()) {
            return response()->json([
                'method' => 'get speedtests in last x days',
                'error' => $validator->errors(),
            ], 422);
        }

        $ttl = Carbon::now()->addDays(1);
        $data = Cache::remember('speedtest-days-' . $days, $ttl, function () use ($days) {
            $showFailed = (bool)SettingsHelper::get('show_failed_tests_on_graph')->value;

            if ($showFailed === true) {
                return Speedtest::where('created_at', '>=', Carbon::now()->subDays($days))
                    ->orderBy('created_at', 'asc')
                    ->get();
            }

            return Speedtest::where('created_at', '>=', Carbon::now()->subDays($days))
                ->where('failed', false)
                ->orderBy('created_at', 'asc')
                ->get();
        });

        return response()->json([
            'method' => 'get speedtests in last x days',
            'days' => $days,
            'data' => $data
        ], 200);
    }

    /**
     * Returns speedtest failure rate going back 'x' days
     *
     * @param   int     $days
     * @return  JsonResponse
     */
    public function fail($days)
    {
        $rule = [
            'days' => ['required', 'integer'],
        ];

        $validator = Validator::make(['days' => $days], $rule);

        if ($validator->fails()) {
            return response()->json([
                'method' => 'get speedtests in last x days',
                'error' => $validator->errors(),
            ], 422);
        }

        $data = SpeedtestHelper::failureRate($days);

        return response()->json([
            'method' => 'get speedtests in last x days',
            'days' => $days,
            'data' => $data
        ], 200);
    }

    /**
     * Return latest speedtest
     *
     * @return  JsonResponse
     */
    public function latest()
    {
        $data = SpeedtestHelper::latest();

        $response = [
            'method' => 'get latest speedtest',
            'data' => $data,
        ];

        if (SettingsHelper::get('show_average')) {
            $avg = Speedtest::select(DB::raw('AVG(ping) as ping, AVG(download) as download, AVG(upload) as upload'))
                ->where('failed', false)
                ->first()
                ->toArray();
            $response['average'] = $avg;
        }

        if (SettingsHelper::get('show_max')) {
            $max = Speedtest::select(DB::raw('MAX(ping) as ping, MAX(download) as download, MAX(upload) as upload'))
                ->where('failed', false)
                ->first()
                ->toArray();
            $response['maximum'] = $max;
        }

        if (SettingsHelper::get('show_min')) {
            $min = Speedtest::select(DB::raw('MIN(ping) as ping, MIN(download) as download, MIN(upload) as upload'))
                ->where('failed', false)
                ->first()
                ->toArray();
            $response['minimum'] = $min;
        }

        if ($data) {
            return response()->json($response, 200);
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
     * @return JsonResponse
     */
    public function run()
    {
        try {
            SettingsHelper::loadIntegrationConfig();
            $data = SpeedtestJob::dispatch(false, config('integrations'));
            return response()->json([
                'method' => 'run speedtest',
                'data' => 'a new speedtest has been added to the queue'
            ], 200);
        } catch (Exception $e) {
            return response()->json([
                'method' => 'run speedtest',
                'error' => $e
            ], 500);
        }
    }

    /**
     * Delete all speedtests from db
     *
     * @return JsonResponse
     */
    public function deleteAll()
    {
        $ret = SpeedtestHelper::deleteAll();

        if ($ret['success']) {
            return response()->json([
                'method' => 'delete all speedtests from the database',
                'success' => true
            ], 200);
        }

        return response()->json([
            'method' => 'delete all speedtests from the database',
            'success' => false,
            'error' => $ret['msg'],
        ], 500);
    }

    /**
     * Delete a specific speedtest from the database
     *
     * @param Speedtest $speedtest
     * @return JsonResponse
     */
    public function delete(Speedtest $speedtest)
    {
        $speedtest->delete();

        Cache::flush();

        return response()->json([
            'method' => 'delete a speedtest from the database',
            'success' => true,
        ], 200);
    }
}
