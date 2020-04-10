<?php

namespace App\Http\Controllers;

use App\Helpers\UpdateHelper;
use Illuminate\Http\Request;

class UpdateController extends Controller
{
    public function checkForUpdate()
    {
        return response()->json([
            'update' => UpdateHelper::check(),
        ], 200);
    }
}
