<?php

namespace App\Helpers;

use Exception;

class UpdateHelper {
    public static function check()
    {
        $current = config('app.version', false);
        (!$current) ? false : '';
    }
}
