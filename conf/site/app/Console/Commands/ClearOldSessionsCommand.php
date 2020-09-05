<?php

namespace App\Console\Commands;

use App\Auth\LoginSession;
use Carbon\Carbon;
use Illuminate\Console\Command;
use Log;

class ClearOldSessionsCommand extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'speedtest:clear-sessions';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Clear expired sessions from database';

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
        $now = Carbon::now()->timestamp;
        $sessions = LoginSession::where('expires', '<=', $now)
                    ->delete();
        $this->info('Invalidated expired sessions');
    }
}
