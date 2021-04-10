<?php

namespace App\Console;

use App\Events\SpeedtestOverviewEvent;
use App\Helpers\SettingsHelper;
use App\Helpers\SpeedtestHelper;
use App\Interfaces\SpeedtestProvider;
use App\Jobs\SpeedtestJob;
use Illuminate\Console\Scheduling\Schedule;
use Illuminate\Foundation\Console\Kernel as ConsoleKernel;

class Kernel extends ConsoleKernel
{
    /**
     * The Artisan commands provided by your application.
     *
     * @var array
     */
    protected $commands = [
        //
    ];

    /**
     * Define the application's command schedule.
     *
     * @param  \Illuminate\Console\Scheduling\Schedule  $schedule
     * @return void
     */
    protected function schedule(Schedule $schedule)
    {
        if ((bool)SettingsHelper::get('schedule_enabled')->value) {
            $schedule->job(
                new SpeedtestJob(
                    true,
                    config('integrations'),
                    app()->make(SpeedtestProvider::class)
                )
            )
                ->cron(SettingsHelper::get('schedule')['value'])
                ->timezone(env('TZ', 'UTC'));
        }
        $schedule->command('speedtest:overview')->cron('0 ' . SettingsHelper::get('speedtest_overview_time')->value . ' * * *')->timezone(env('TZ', 'UTC'));
        $schedule->command('speedtest:clear-sessions')->everyMinute();
    }

    /**
     * Register the commands for the application.
     *
     * @return void
     */
    protected function commands()
    {
        $this->load(__DIR__ . '/Commands');

        require base_path('routes/console.php');
    }
}
