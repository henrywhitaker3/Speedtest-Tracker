<?php

namespace Tests\Unit\Helpers\SpeedtestHelper;

use App\Helpers\SpeedtestHelper;
use App\Utils\OoklaTester;
use Illuminate\Foundation\Testing\RefreshDatabase;
use JsonException;
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

        $this->output = $this->speedtestProvider->output();
    }

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testOutputFunction()
    {
        $this->assertJson($this->output);
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
        $json = '{hi: hi}';

        $o = $this->speedtestProvider->run($json);

        $this->assertFalse($o);
    }

    /**
     * A basic unit test example.
     *
     * @return void
     */
    public function testIncompleteJson()
    {
        $json = '{"hi": "hi"}';

        $o = $this->speedtestProvider->run($json);

        $this->assertFalse($o);
    }
}
