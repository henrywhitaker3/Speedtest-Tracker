<?php

namespace Henrywhitaker3\LaravelActions\Tests\Unit;

use Henrywhitaker3\LaravelActions\Tests\Utils\ExampleInstantiatedAction;
use Henrywhitaker3\LaravelActions\Tests\Utils\ExampleMultiArgumentAction;
use Henrywhitaker3\LaravelActions\Tests\Utils\ExampleNoArgumentAction;
use Henrywhitaker3\LaravelActions\Tests\Utils\ExampleNonInstantiatedAction;

class ActionTest extends \Orchestra\Testbench\TestCase
{
    protected function getPackageProviders($app)
    {
        return ['Henrywhitaker3\LaravelActions\LaravelActionsServiceProvider'];
    }

    /**
     * Run an action when it's already been instantiated.
     *
     * @test
     * @return void
     */
    public function runInstantiatedAction()
    {
        $text = 'instantiated text';

        $action = new ExampleInstantiatedAction($text);
        $output = run($action);

        $this->assertEquals($text, $output);
    }

    /**
     * Run an action when it's not been instantiated.
     *
     * @test
     * @return void
     */
    public function runNonInstantiatedAction()
    {
        $text = 'non-instantiated text';

        $output = run(ExampleNonInstantiatedAction::class, $text);

        $this->assertEquals($text, $output);
    }

    /**
     * Run a multi-argaction when it's not been
     * instantiated.
     *
     * @test
     * @return void
     */
    public function runMultiArgumentAction()
    {
        $text = 'multi-argument ';
        $text2 = 'text';

        $output = run(ExampleMultiArgumentAction::class, $text, $text2);

        $this->assertEquals($text.$text2, $output);
    }

    /**
     * Run an action with no args when it's already
     * been instantiated.
     *
     * @test
     * @return void
     */
    public function runInstantiatedNoArgumentsAction()
    {
        $output = run(ExampleNoArgumentAction::class);

        $this->assertEquals('Hello', $output);
    }
}
