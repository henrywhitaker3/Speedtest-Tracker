<?php

namespace App\Actions;

use App\Helpers\SettingsHelper;
use App\Models\Speedtest;
use Cache;
use Carbon\Carbon;
use Henrywhitaker3\LaravelActions\Interfaces\ActionInterface;

class GetSpeedtestTimeData implements ActionInterface
{
    /**
     * Run the action.
     *
     * @return mixed
     */
    public function run($days = 7)
    {
        $ttl = Carbon::now()->addDays(1);

        return Cache::remember('speedtest-days-' . $days, $ttl, function () use ($days) {
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
    }
}
