<?php

namespace App\Facades;

use Illuminate\Support\Facades\Facade;

class UpdaterFacade extends Facade
{
    protected static function getFacadeAccessor()
    {
        return 'updater';
    }
}
