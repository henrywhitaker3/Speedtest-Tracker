<?php

namespace App\Console\Commands;

use App\Helpers\SettingsHelper;
use Illuminate\Console\Command;

class SetSlackWebhook extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'speedtest:slack {webhook : The slack webhook to store}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Set the slack webhook setting';

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
        $webhook = $this->argument('webhook');

        SettingsHelper::set('slack_webhook', $webhook);

        $this->info('Slack webhook updated');
    }
}
