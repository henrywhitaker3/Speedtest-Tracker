<?php

namespace App\Facades;

use Illuminate\Support\Facades\Facade;

class HealthchecksFacade extends Facade
{
    protected static function getFacadeAccessor()
    {
        return 'healthcheck';
    }
}
