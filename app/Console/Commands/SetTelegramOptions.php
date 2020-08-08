<?php

namespace App\Console\Commands;

use App\Helpers\SettingsHelper;
use Illuminate\Console\Command;

class SetTelegramOptions extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'speedtest:telegram
                            {--bot= : The telegram bot token}
                            {--chat= : The telegram chat ID}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Set the telegram settings';

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
        $bot = $this->option('bot');
        $chat = $this->option('chat');

        SettingsHelper::set('telegram_bot_token', $bot);
        SettingsHelper::set('telegram_chat_id', $chat);

        $this->info('Telegram options updated');
    }
}
