<?php

namespace App\Console\Commands;

use App\Helpers\SettingsHelper;
use Illuminate\Console\Command;

class GetConfig extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'speedtest:config';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Get the application configuration';

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
        $this->info(json_encode(SettingsHelper::getConfig(), JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES));
    }
}
