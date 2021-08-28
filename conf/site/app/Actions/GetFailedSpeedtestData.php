<?php

namespace App\Actions;

use App\Models\Speedtest;
use Cache;
use Carbon\Carbon;
use DB;
use Henrywhitaker3\LaravelActions\Interfaces\ActionInterface;

class GetFailedSpeedtestData implements ActionInterface
{
    /**
     * Run the action.
     *
     * @return mixed
     */
    public function run($days = 7)
    {
        $ttl = Carbon::now()->addDays(1);

        return Cache::remember('failure-rate-' . $days, $ttl, function () use ($days) {
            $range = [
                Carbon::today()
            ];
            for ($i = 0; $i < ($days - 1); $i++) {
                $prev = end($range);
                $new = $prev->copy()->subDays(1);
                array_push($range, $new);
            }

            $rate = [];

            foreach ($range as $day) {
                $success = Speedtest::select(DB::raw('COUNT(id) as rate'))->whereDate('created_at', $day)->where('failed', false)->get()[0]['rate'];
                $fail = Speedtest::select(DB::raw('COUNT(id) as rate'))->whereDate('created_at', $day)->where('failed', true)->get()[0]['rate'];

                array_push($rate, [
                    'date' => $day->toDateString(),
                    'success' => $success,
                    'failure' => $fail,
                ]);
            }

            return array_reverse($rate);
        });
    }
}
