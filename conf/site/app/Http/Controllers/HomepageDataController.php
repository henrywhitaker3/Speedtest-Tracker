<?php

namespace App\Http\Controllers;

use App\Actions\GetFailedSpeedtestData;
use App\Actions\GetLatestSpeedtestData;
use App\Actions\GetSpeedtestTimeData;
use App\Helpers\SettingsHelper;
use Illuminate\Http\Request;
use Validator;

class HomepageDataController extends Controller
{
    public function __invoke($days)
    {
        $validator = Validator::make(
            ['days' => $days],
            ['days' => ['required', 'numeric']],
        );

        if ($validator->fails()) {
            return response()->json([
                'method' => 'get speedtests in last x days',
                'error' => $validator->errors(),
            ], 422);
        }

        return [
            'latest' => run(GetLatestSpeedtestData::class),
            'time' => run(GetSpeedtestTimeData::class),
            'fail' => run(GetFailedSpeedtestData::class),
            'config' => SettingsHelper::getConfig(),
        ];
    }
}
