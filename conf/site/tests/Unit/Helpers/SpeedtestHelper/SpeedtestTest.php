<?php

namespace Tests\Unit\Helpers\SpeedtestHelper;

use App\Exceptions\SpeedtestFailureException;
use App\Helpers\SpeedtestHelper;
use App\Utils\OoklaTester;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\Mocks\OoklaTesterMocker;
use Tests\TestCase;

class SpeedtestTest extends TestCase
{
    use RefreshDatabase;

    private $output;

    private OoklaTester $speedtestProvider;

    public function setUp(): void
    {
        parent::setUp();

        $this->speedtestProvider = new OoklaTester();

        $this->output = (new OoklaTesterMocker())->output();
    }

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testRunSpeedtestWithExistingOutput()
    {
        $output = json_decode($this->output, true);

        $test = $this->speedtestProvider->run($this->output);

        $this->assertEquals($output['ping']['latency'], $test->ping);
        $this->assertEquals(SpeedtestHelper::convert($output['download']['bandwidth']), $test->download);
        $this->assertEquals(SpeedtestHelper::convert($output['upload']['bandwidth']), $test->upload);
    }

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testInvaidJson()
    {
        $this->expectException(SpeedtestFailureException::class);

        $json = '{hi: hi}';

        $o = $this->speedtestProvider->run($json);
    }

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testIncompleteJson()
    {
        $this->expectException(SpeedtestFailureException::class);

        $json = '{"hi": "hi"}';

        $o = $this->speedtestProvider->run($json);
    }
}
