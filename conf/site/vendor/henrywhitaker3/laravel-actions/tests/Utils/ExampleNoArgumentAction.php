<?php

namespace Henrywhitaker3\LaravelActions\Tests\Utils;

use Henrywhitaker3\LaravelActions\Interfaces\ActionInterface;

class ExampleNoArgumentAction implements ActionInterface
{
    public function run()
    {
        return 'Hello';
    }
}
