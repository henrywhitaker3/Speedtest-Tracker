<?php

namespace Tests\Unit\Helpers\SpeedtestHelper;

use App\Helpers\SpeedtestHelper;
use Illuminate\Foundation\Testing\RefreshDatabase;
use JsonException;
use Tests\TestCase;

class SpeedtestTest extends TestCase
{
    use RefreshDatabase;

    private $output;

    public function setUp() : void
    {
        parent::setUp();

        $this->output = SpeedtestHelper::output();
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

        $test = SpeedtestHelper::runSpeedtest($this->output);

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

        $o = SpeedtestHelper::runSpeedtest($json);

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

        $o = SpeedtestHelper::runSpeedtest($json);

        $this->assertFalse($o);
    }
}
