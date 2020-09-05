<?php

namespace App\Console\Commands;

use App\Helpers\SettingsHelper;
use Illuminate\Console\Command;

class AuthenticationCommand extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'speedtest:auth {--enable} {--disable}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Toggle authentication for the app';

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
     * @return int
     */
    public function handle()
    {
        $opts = $this->options();

        if($opts['enable'] === true && $opts['disable'] === true) {
            $this->warn('Please specify only ONE of --enable and --disable');
        } else if($opts['enable'] === false && $opts['disable'] === false) {
            $this->warn('You need to specify either --enable OR --disable');
        } else {
            if($opts['enable'] === true) {
                $this->info('Enabling authentication');
                SettingsHelper::set('auth', true);
            }

            if($opts['disable'] === true) {
                $this->info('Disabling authentication');
                SettingsHelper::set('auth', false);
            }
        }
    }
}
