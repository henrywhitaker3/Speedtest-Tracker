<?php

namespace App\Http\Controllers;

use Exception;
use Updater;
use Illuminate\Http\Request;

class UpdateController extends Controller
{

    /**
     * Check for new update
     *
     * @return  Response
     */
    public function checkForUpdate()
    {
        return response()->json([
            'method' => 'check for updates',
            'update' => Updater::check(),
        ], 200);
    }

    /**
     * Download new update
     *
     * @return  Response
     */
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

    /**
     * Trigger update extraction
     *
     * @return  Response
     */
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

    /**
     * Trigger update file move
     *
     * @return  Response
     */
    public function moveUpdate()
    {
        $cp = Updater::updateFiles();

        return response()->json([
            'method' => 'copy latest version',
            'success' => $cp,
        ], 200);
    }

    /**
     * Get local changelog
     *
     * @return  Response
     */
    public function changelog()
    {
        $url = base_path() . '/changelog.json';

        try {
            $changelog = json_decode(file_get_contents($url), true);
        } catch(Exception $e) {
            $changelog = [];
        }

        return response()->json([
            'method' => 'get changelog',
            'data' => $changelog
        ], 200);
    }
}
