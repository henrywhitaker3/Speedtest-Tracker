<?php

namespace Tests\Feature\Commands;

use App\Console\Commands\AuthenticationCommand;
use App\Helpers\SettingsHelper;
use Artisan;
use Illuminate\Console\Command;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Foundation\Testing\WithFaker;
use Symfony\Component\Console\Tester\CommandTester;
use Tests\TestCase;

class AuthenticationCommandTest extends TestCase
{
    use RefreshDatabase;

    /**
     * Enable app auth
     *
     * @return void
     */
    public function testEnableAuth()
    {
        SettingsHelper::set('auth', false);

        $this->assertEquals(Artisan::call('speedtest:auth', [ '--enable' => true ]), 0);

        $this->assertTrue((bool)SettingsHelper::get('auth')->value);
    }

    /**
     * Disable app auth
     *
     * @return void
     */
    public function testDisableAuth()
    {
        SettingsHelper::set('auth', true);

        $this->assertEquals(Artisan::call('speedtest:auth', [ '--disable' => true ]), 0);

        $this->assertFalse((bool)SettingsHelper::get('auth')->value);
    }

    /**
     * Test invalid params for command
     *
     * @return void
     */
    // public function testAuthBothOptions()
    // {
    //     $command = new AuthenticationCommand();
    //     $command->setLaravel($this->app);
    //     $tester = new CommandTester($command);
    //     $tester->setInputs([]);
    //     $tester->execute([ '--enable' => true, '--disable' => true ]);
    // }

    /**
     * Test invalid params for command
     *
     * @return void
     */
    // public function testAuthNoOptions()
    // {
    //     $command = new AuthenticationCommand();
    //     $command->setLaravel($this->app);
    //     $tester = new CommandTester($command);
    //     $tester->setInputs([]);
    //     $tester->execute([]);
    // }
}
