<?php

namespace Henrywhitaker3\LaravelActions\Tests\Utils;

use Henrywhitaker3\LaravelActions\Interfaces\ActionInterface;

class ExampleNonInstantiatedAction implements ActionInterface
{
    public function run(string $text = null)
    {
        return $text;
    }
}
