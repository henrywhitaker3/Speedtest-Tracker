<?php

namespace App\Console\Commands;

use App\Helpers\SettingsHelper;
use Illuminate\Console\Command;

class TestNotification extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'speedtest:notification';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Send test notifications to all notification agents';

    /**
     * Create a new command instance.
     *
     * @return void
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * Execute the console command.
     *
     * @return void
     */
    public function handle()
    {
        SettingsHelper::testNotification();
    }
}
