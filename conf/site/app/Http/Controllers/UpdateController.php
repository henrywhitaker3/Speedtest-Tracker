<?php

namespace App\Http\Controllers;

use Updater;
use Illuminate\Http\Request;

class UpdateController extends Controller
{
    public function checkForUpdate()
    {
        return response()->json([
            'method' => 'check for updates',
            'update' => Updater::check(),
        ], 200);
    }

    public function downloadUpdate()
    {
        $dl = Updater::downloadLatest();

        if($dl) {
            return response()->json([
                'method' => 'download latest version',
                'success' => true,
            ], 200);
        } else {
            return response()->json([
                'method' => 'download latest version',
                'success' => false,
            ], 500);
        }
    }

    public function extractUpdate()
    {
        $ex = Updater::extractFiles();

        if($ex) {
            return response()->json([
                'method' => 'extract latest version',
                'success' => true,
            ], 200);
        } else {
            return response()->json([
                'method' => 'extract latest version',
                'success' => false,
            ], 500);
        }
    }

    public function moveUpdate()
    {
        $cp = Updater::updateFiles();

        return response()->json([
            'method' => 'copy latest version',
            'success' => $cp,
        ], 200);
    }
}
