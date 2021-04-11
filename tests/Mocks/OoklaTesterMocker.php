<?php

namespace Tests\Mocks;

use App\Exceptions\SpeedtestFailureException;
use App\Interfaces\SpeedtestProvider;
use App\Models\Speedtest;
use Exception;

class OoklaTesterMocker implements SpeedtestProvider
{
    private bool $passes;

    public function __construct(bool $passes = true)
    {
        $this->passes = $passes;
    }

    public function run($output = null): Speedtest
    {
        $output = $output ?? $this->output();

        try {
            $output = json_decode($output, true, 512, JSON_THROW_ON_ERROR);
        } catch (Exception $e) {
            throw new SpeedtestFailureException();
        }

        return $this->passes
            ? Speedtest::factory()->create()
            : Speedtest::factory()->create([
                'download' => 0,
                'upload' => 0,
                'ping' => 0,
                'failed' => true,
            ]);
    }

    public function output()
    {
        return !$this->passes
            ? null
            : json_encode([
                'type' => 'result',
                'download' => ['bandwidth' => '50'],
                'upload' => ['bandwidth' => '50'],
                'ping' => ['latency' => '50'],
                'server' => [
                    'id' => '1',
                    'name' => 'PHPUnit',
                    'host' => 'phpunit',
                    'port' => '443',
                ],
                'result' => [
                    'url' => 'some-url',
                ]
            ]);
    }
}
