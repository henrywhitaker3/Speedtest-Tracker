<?php

namespace Henrywhitaker3\LaravelActions\Tests\Utils;

use Henrywhitaker3\LaravelActions\Interfaces\ActionInterface;

class ExampleInstantiatedAction implements ActionInterface
{
    public string $text;

    public function __construct(string $text)
    {
        $this->text = $text;
    }

    public function run()
    {
        return $this->text;
    }
}
